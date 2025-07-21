package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"poplargrid/internal/apiserver/apiclient"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
	"poplargrid/internal/shared/txutils"

	"gorm.io/gorm"
)

// ProjectService 接口定义了项目服务的基本操作
type ProjectService interface {
	// GetBasicPageWithParams 获取项目列表，支持分页和排序以及复合条件查询
	GetBasicPageWithParams(worksetId uint, pageSerial, pageSize int, sort int, userId uint, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error)
	// GetBasicPageByUserId 获取指定用户参与的项目列表，支持分页
	GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.MyProjectBasic, error)

	// GetDetail 获取指定 ID 的项目详情
	GetDetail(projectId uint) (*dtos.ProjectDetail, error)

	// GetLaborDivision 获取指定项目 ID 的团队成员分工
	GetLaborDivision(projectId uint) ([]*dtos.LaborDivision, error)

	// CreateProject 创建一个新的项目
	CreateProject(request *dtos.CreateProjectInfo) (*dtos.ProjectCreatedInfo, error)

	// UpdateProject 更新指定 ID 的项目
	UpdateProject(projectId uint, request *dtos.UpdateProjectRequest) error

	// DeleteProject 删除指定 ID 的项目
	DeleteProject(projectId uint) error
}

// projectServiceImpl 是 ProjectService 的实现
type projectServiceImpl struct {
	projectRepo repos.ProjectRepo
	laborRepo   repos.LaborRepo
	userRepo    repos.UserRepo
	apiClient   apiclient.ApiClient
	logger      *slog.Logger
}

// NewProjectService 创建一个新的 ProjectService 实例
func NewProjectService(
	projectRepo repos.ProjectRepo,
	laborRepo repos.LaborRepo,
	userRepo repos.UserRepo,
	apiClient apiclient.ApiClient,
	logger *slog.Logger,
) ProjectService {
	return &projectServiceImpl{
		projectRepo: projectRepo,
		laborRepo:   laborRepo,
		userRepo:    userRepo,
		apiClient:   apiClient,
		logger:      logger,
	}
}

// GetBasicPageWithParams 实现 ProjectService 接口的 GetBasicPageWithParams 方法
func (s *projectServiceImpl) GetBasicPageWithParams(
	worksetId uint, pageSerial, pageSize int, sort int,
	userId uint, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error) {
	// 解析构建最终的查询条件
	queryParams := s.buildQueryParams(worksetId, userId, sort, status)

	var (
		page []*dbmodels.Project
		err  error
	)

	switch sort {
	case dtos.SORT_ID_DESC, dtos.SORT_UPDATED_AT_DESC:
		// 合法的 sort 参数，直接查询
		page, err = s.projectRepo.SelectBasicPageWithParam(
			(pageSerial-1)*pageSize, pageSize, queryParams)
		if err != nil {
			s.logger.Error("GetBasicPageWithParams 调用 SelectBasicPageWithParam 中出现错误", slog.Any("error", err))
			return nil, fmt.Errorf("获取项目列表失败")
		}

	default:
		s.logger.Error("GetBasicPage 不支持的排序方式", slog.Int("sort", sort))
		return nil, fmt.Errorf("不支持的排序方式：%d", sort)
	}

	// 将 dbmodels.Project 转换为 dtos.ProjectBasic
	var projectBasics []*dtos.ProjectBasic

	for _, project := range page {
		var status dtos.ProjectOverallStatus
		status.SetTranslatingStatus(uint(project.TranslateStatus))
		status.SetProofreadingStatus(uint(project.ProofStatus))
		status.SetLetteringStatus(uint(project.LetterStatus))
		status.SetReviewingStatus(uint(project.ReviewStatus))

		projectBasics = append(projectBasics, &dtos.ProjectBasic{
			Id:            uint(project.Id),
			Title:         project.Title,
			WorksetId:     uint(project.WorksetId),
			WorksetIndex:  project.WorksetIndex,
			LegacyId:      uint(project.LegacyId),
			MoetranId:     project.MoetranId,
			Status:        status,
			IsPublished:   project.IsPublished,
			AllowAutoJoin: project.AllowAutoJoin,
		})
	}

	return projectBasics, nil
}

// GetBasicPageByUserId 实现 ProjectService 接口的 GetBasicPageByUserId 方法
func (s *projectServiceImpl) GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.MyProjectBasic, error) {
	// 获取用户参与的项目列表
	labors, err := s.laborRepo.SelectProjectPageByUserId(dbmodels.PrimaryKey(userId), (pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetBasicPageByUserId 调用 SelectBasicPageIdDescByUserId 中出现错误", slog.Any("error", err))
		return nil, fmt.Errorf("获取用户参与的项目列表失败")
	}

	// 将 dbmodels.ProjectLaborDivision 转换为 dtos.MyProjectBasic
	var myProjects []*dtos.MyProjectBasic

	for _, labor := range labors {
		var status dtos.ProjectOverallStatus
		status.SetTranslatingStatus(uint(labor.FkProject.TranslateStatus))
		status.SetProofreadingStatus(uint(labor.FkProject.ProofStatus))
		status.SetLetteringStatus(uint(labor.FkProject.LetterStatus))
		status.SetReviewingStatus(uint(labor.FkProject.ReviewStatus))

		myProjects = append(myProjects, &dtos.MyProjectBasic{
			ProjectBasic: dtos.ProjectBasic{
				Id:            uint(labor.FkProject.Id),
				Title:         labor.FkProject.Title,
				WorksetId:     uint(labor.FkProject.WorksetId),
				WorksetIndex:  labor.FkProject.WorksetIndex,
				LegacyId:      uint(labor.FkProject.LegacyId),
				MoetranId:     labor.FkProject.MoetranId,
				Status:        status,
				IsPublished:   labor.FkProject.IsPublished,
				AllowAutoJoin: labor.FkProject.AllowAutoJoin,
			},
			Role: uint(labor.LaborRole),
		})
	}

	return myProjects, nil
}

// GetDetail 实现 ProjectService 接口的 GetDetail 方法
func (s *projectServiceImpl) GetDetail(projectId uint) (*dtos.ProjectDetail, error) {
	// 调用仓库方法获取项目详情
	project, err := s.projectRepo.SelectById(dbmodels.PrimaryKey(projectId))
	if err != nil {
		s.logger.Error("GetDetail 调用 SelectById 中出现错误", slog.Any("error", err))
		return nil, fmt.Errorf("获取项目详情失败")
	}

	// 将 dbmodels.Project 转换为 dtos.ProjectDetail
	var status dtos.ProjectOverallStatus
	status.SetTranslatingStatus(uint(project.TranslateStatus))
	status.SetProofreadingStatus(uint(project.ProofStatus))
	status.SetLetteringStatus(uint(project.LetterStatus))
	status.SetReviewingStatus(uint(project.ReviewStatus))

	detail := &dtos.ProjectDetail{
		ProjectBasic: dtos.ProjectBasic{
			Id:            uint(project.Id),
			Title:         project.Title,
			WorksetId:     uint(project.WorksetId),
			WorksetIndex:  project.WorksetIndex,
			LegacyId:      uint(project.LegacyId),
			Status:        status,
			IsPublished:   project.IsPublished,
			AllowAutoJoin: project.AllowAutoJoin,
		},
		Description: project.Description,
		CreatedAt:   project.CreatedAt.Format(dtos.DTO_TIME_FORMAT),
		UpdatedAt:   project.UpdatedAt.Format(dtos.DTO_TIME_FORMAT),
	}

	return detail, nil
}

// GetLaborDivision 实现 ProjectService 接口的 GetLaborDivision 方法
func (s *projectServiceImpl) GetLaborDivision(projectId uint) ([]*dtos.LaborDivision, error) {
	// 调用仓库方法获取团队成员分工
	members, err := s.laborRepo.SelectByProjectId(dbmodels.PrimaryKey(projectId))
	if err != nil {
		s.logger.Error("GetLaborDivision 调用 SelectLaborDivisionByProjectId 中出现错误", slog.Any("error", err))
		return nil, fmt.Errorf("获取项目成员分工失败")
	}

	// 将 dbmodels.TeamMember 转换为 dtos.LaborDivision
	var divisions []*dtos.LaborDivision
	for _, member := range members {
		divisions = append(divisions, &dtos.LaborDivision{
			MemberId: uint(member.Id),
			Nickname: member.FkUser.Nickname,
			Role:     uint(member.LaborRole),
		})
	}

	return divisions, nil
}

// CreateProject 实现 ProjectService 接口的 CreateProject 方法
func (s *projectServiceImpl) CreateProject(createInfo *dtos.CreateProjectInfo) (*dtos.ProjectCreatedInfo, error) {
	// 构建一个事务协调器
	coordinater := txutils.NewTransactionCoordinator(s.projectRepo.GetHandle())

	var createdInfo *dtos.ProjectCreatedInfo

	// 在一个事务中协调数据库和调用龙译 API 的操作
	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 基于 transaction 上下文获取仓库实例
		txProjectRepo := repos.NewProjectRepo(tx)
		txLaborRepo := repos.NewLaborRepo(tx)

		// 创建一个新的项目
		project := &dbmodels.Project{
			Title:         createInfo.Title,
			Description:   createInfo.Description,
			WorksetId:     dbmodels.PrimaryKey(createInfo.WorksetId),
			AllowAutoJoin: createInfo.AllowAutoJoin,
			IsHidden:      createInfo.IsHidden,
		}

		// 在数据库中创建项目，如果成功 project 的 Id、WorksetIndex 和 FkWorkset 应当被填充
		if err := txProjectRepo.CreateProject(project); err != nil {
			s.logger.Error("CreateProject 调用 CreateProject 中出现错误", slog.Any("error", err))
			return errors.New("创建项目失败"), nil
		}

		// 创建 creator 的分工记录
		creatorLabor := dbmodels.LaborMask(0)
		creatorLabor.AddRole(dbmodels.LABOR_CREATOR_MASK)

		laborDivision := &dbmodels.ProjectLaborDivision{
			ProjectId: dbmodels.PrimaryKey(project.Id),
			UserId:    dbmodels.PrimaryKey(createInfo.CreatorUserId),
			LaborRole: creatorLabor,
			PrincipalId: dbmodels.PrimaryKey(createInfo.CreatorUserId), // 负责人为创建者
		}

		// 然后在团队成员分工表中插入 creator 记录
		if err := txLaborRepo.CreateLaborDivision(laborDivision); err != nil {
			s.logger.Error("CreateProject 调用 CreateLaborDivision 中出现错误", slog.Any("error", err))
			return fmt.Errorf("创建团队成员分工失败"), nil
		}

		// 获取对应 user 的龙译 JWT
		user, err := s.userRepo.SelectByUserId(dbmodels.PrimaryKey(createInfo.CreatorUserId))
		if err != nil {
			s.logger.Error("CreateProject 获取用户信息失败", slog.Any("error", err))
			return fmt.Errorf("获取用户信息失败"), nil
		}

		// 随后调用龙译 API 创建项目
		projInfo := s.buildMoetranProjInfo(createInfo, project, user)

		moetranRes, err := s.apiClient.CreateProject(projInfo)
		if err != nil {
			s.logger.Error("CreateProject 调用 CreateProject API 中出现错误", slog.Any("error", err))
			// TODO：由于不确定尨译的 API 是否是幂等的，这里需要考虑补偿操作，比如删除对应项目
			// 但在不确定尨译实现的情况下，先不处理补偿
			return errors.New("调用龙译 API 创建项目失败"), nil
		}

		// 将龙译返回的项目 ID 更新到本地项目中
		if err := txProjectRepo.UpdateMoetranId(project.Id, moetranRes.Project.Id); err != nil {
			s.logger.Error("CreateProject 更新本地项目龙译 ID 失败", slog.Any("error", err))
			return fmt.Errorf("更新本地龙译 ID 失败"), nil
		}

		// 写入到返回结果
		createdInfo = &dtos.ProjectCreatedInfo{
			Message:   moetranRes.Message,
			ProjectId: uint(project.Id),
			MoetranId: moetranRes.Project.Id,
		}

		// 一切正常，则返回 nil 提交事务
		return nil, nil

	}); err != nil {
		// 这里的 err 是上述事务中抛出的错误
		s.logger.Error("CreateProject 事务执行失败", slog.Any("error", err))
		return nil, errors.New("创建项目失败")
	}

	return createdInfo, nil
}

// UpdateProject 实现 ProjectService 接口的 UpdateProject 方法
func (s *projectServiceImpl) UpdateProject(projectId uint, request *dtos.UpdateProjectRequest) error {
	// 处理 status 字段
	var (
		translateStatus uint8
		proofStatus     uint8
		letterStatus    uint8
		reviewStatus    uint8
		isPublished     bool
	)

	if request.Status != 0 {
		status := dtos.ProjectOverallStatus(request.Status)

		translateStatus = uint8(status.GetTranslatingStatus())
		proofStatus = uint8(status.GetProofreadingStatus())
		letterStatus = uint8(status.GetLetteringStatus())
		reviewStatus = uint8(status.GetReviewingStatus())

		publishStatus := uint8(status.GetPublishedStatus())
		if publishStatus == uint8(dtos.PROJECT_STATUS_COMPLETED) {
			isPublished = true
		}
	}

	// 构造要 save 的项目信息
	project := &dbmodels.Project{
		BaseModel: dbmodels.BaseModel{
			Id: dbmodels.PrimaryKey(projectId),
		},
		Title:       request.Title,
		Description: request.Description,

		TranslateStatus: translateStatus,
		ProofStatus:     proofStatus,
		LetterStatus:    letterStatus,
		ReviewStatus:    reviewStatus,
		IsPublished:     isPublished,
	}

	// 更新项目信息数据到数据库
	if err := s.projectRepo.SaveInfo(project); err != nil {
		s.logger.Error("UpdateProject 调用 Update 中出现错误", slog.Any("error", err))
		return fmt.Errorf("更新项目失败")
	}

	return nil
}

// DeleteProject 实现 ProjectService 接口的 DeleteProject 方法
func (s *projectServiceImpl) DeleteProject(projectId uint) error {
	// 检查项目 ID 是否有效
	if projectId == 0 {
		return errors.New("项目 ID 不能为 0")
	}

	// 调用仓库方法删除项目
	if err := s.projectRepo.DeleteById(dbmodels.PrimaryKey(projectId)); err != nil {
		s.logger.Error("DeleteProject 调用 DeleteById 中出现错误", slog.Any("error", err))
		return fmt.Errorf("删除项目失败")
	}

	return nil
}

// ================ 辅助函数 ================

// buildQueryParams 通过 service 参数为 queryParams 构建查询条件
func (s *projectServiceImpl) buildQueryParams(worksetId, userId uint, sort int, status dtos.ProjectOverallStatus) *dtos.ProjectSearchParams {
	queryParams := &dtos.ProjectSearchParams{}

	// 解析 workset id，其为 0 时不筛选
	if worksetId != 0 {
		queryParams.WorksetId = &worksetId
	}

	// 解析 user id，其为 0 时不筛选
	if userId != 0 {
		queryParams.UserId = &userId
	}

	// 加入 sort 参数
	queryParams.Sort = sort

	// 解析 status 位掩码，获取各个阶段的状态
	statusQueryParams := dtos.ProjectStatusQueryParams{}
	{
		withParams := false

		// 翻译状态
		if s := status.GetTranslatingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			statusQueryParams.TranslateStatus = &s
			withParams = true
		}

		// 校对状态
		if s := status.GetProofreadingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			statusQueryParams.ProofStatus = &s
			withParams = true
		}

		// 嵌字状态
		if s := status.GetLetteringStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			statusQueryParams.LetterStatus = &s
			withParams = true
		}

		// 审核状态
		if s := status.GetReviewingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			statusQueryParams.ReviewStatus = &s
			withParams = true
		}

		// 发布状态要特别处理
		if s := status.GetPublishedStatus(); s != dtos.PROJECT_STATUS_NOT_USED && s != dtos.PROJECT_STATUS_IN_PROGRESS {
			switch s {
			case dtos.PROJECT_STATUS_UNSET:
				isPublished := false
				statusQueryParams.PublishStatus = &isPublished
			case dtos.PROJECT_STATUS_COMPLETED:
				isPublished := true
				statusQueryParams.PublishStatus = &isPublished
			}
			withParams = true
		}

		// 如果有 status 查询条件，则设置到 queryParams 中
		if withParams {
			queryParams.Status = &statusQueryParams
		}
	}

	return queryParams
}

// buildMoetranProjInfo 构建 Moetran 项目信息
func (s *projectServiceImpl) buildMoetranProjInfo(
	info *dtos.CreateProjectInfo,
	project *dbmodels.Project,
	user *dbmodels.User,
) *apiclient.CreateProjectInfo {
	// 组装尨译的 title
	title := fmt.Sprintf("[%d-%d] %s",
		project.WorksetId, project.WorksetIndex, project.Title)

	// 构造 AllowApplyType 和 ApplicationCheckType
	allowApplyType := -1
	applicationCheckType := -1

	switch info.AllowAutoJoin {
	case true:
		allowApplyType = apiclient.ALLOW_ANY_APPLI
		applicationCheckType = apiclient.APPLI_NON_CHECK
	case false:
		allowApplyType = apiclient.ALLOW_MEMBER_ONLY
		applicationCheckType = apiclient.APPLI_ADMIN_CHECK
	}

	return &apiclient.CreateProjectInfo{
		MoetranAuth: user.MoetranAuth,

		Title:            title,
		Description:      project.Description,
		MoetranProjSetId: project.FkWorkset.MoetranId,
		MoetranTeamId:    project.FkWorkset.FkTeam.MoetranId,
		WorksetIndex:     project.WorksetIndex,

		// TODO：先写死为 ja 到 zh-TW
		SourceLanguage:  apiclient.LangJapanese,
		TargetLanguages: []string{apiclient.LangTraditionalChinese},

		AllowApplyType:       allowApplyType,
		ApplicationCheckType: applicationCheckType,

		// TODO：先写死默认角色为实习翻译
		DefaultRole: apiclient.SystemRoleIDsMap[apiclient.ROLE_INTERN],
	}
}
