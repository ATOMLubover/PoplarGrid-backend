package services

import (
	"errors"
	"fmt"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// ProjectStatus 定义了项目的整体状态
// 其本质上也是一个位掩码，表示项目在各个阶段的状态
type ProjectStatus uint32

// ProjectStatus 常量定义了项目状态的位掩码
const (
	PROJECT_STATUS_NOT_USED    uint32 = iota // 未使用
	PROJECT_STATUS_UNSET                     // 未设置
	PROJECT_STATUS_IN_PROGRESS               // 进行中
	PROJECT_STATUS_COMPLETED                 // 已完成
)

const (
	kProjectTranslateShift uint32 = iota * 2 // 翻译阶段状态位
	kProjectProofShift                       // 校对阶段状态位
	kProjectLetterShift                      // 嵌字阶段状态位
	kProjectReviewShift                      // 审核阶段状态位
	kProjectPublishShift                     // 发布阶段状态位
)

// 各阶段状态对应的掩码常量 (用于写入和读取)
const (
	PROJECT_STATUS_TRANSLATE_MASK uint32 = 3 << (2 * iota) // 二进制 3
	PROJECT_STATUS_PROOF_MASK                              // 二进制 12
	PROJECT_STATUS_LETTER_MASK                             // 二进制 48
	PROJECT_STATUS_REVIEW_MASK                             // 二进制 192
	PROJECT_STATUS_PUBLISH_MASK                            // 二进制 768
)

// GetTranslatingStatus 获取翻译阶段的状态
func (s ProjectStatus) GetTranslatingStatus() uint32 {
	return (uint32(s) & PROJECT_STATUS_TRANSLATE_MASK) >> kProjectTranslateShift
}

// GetProofreadingStatus 获取校对阶段的状态
func (s ProjectStatus) GetProofreadingStatus() uint32 {
	return (uint32(s) & PROJECT_STATUS_PROOF_MASK) >> kProjectProofShift
}

// GetLetteringStatus 获取嵌字阶段的状态
func (s ProjectStatus) GetLetteringStatus() uint32 {
	return (uint32(s) & PROJECT_STATUS_LETTER_MASK) >> kProjectLetterShift
}

// GetReviewingStatus 获取审核阶段的状态
func (s ProjectStatus) GetReviewingStatus() uint32 {
	return (uint32(s) & PROJECT_STATUS_REVIEW_MASK) >> kProjectReviewShift
}

// GetPublishedStatus 获取发布阶段的状态
func (s ProjectStatus) GetPublishedStatus() uint32 {
	return (uint32(s) & PROJECT_STATUS_PUBLISH_MASK) >> kProjectPublishShift
}

// SetTranslatingStatus 设置翻译阶段的状态
func (s *ProjectStatus) SetTranslatingStatus(status uint32) {
	*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_TRANSLATE_MASK | (status << kProjectTranslateShift))
}

// SetProofreadingStatus 设置校对阶段的状态
func (s *ProjectStatus) SetProofreadingStatus(status uint32) {
	*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_PROOF_MASK | (status << kProjectProofShift))
}

// SetLetteringStatus 设置嵌字阶段的状态
func (s *ProjectStatus) SetLetteringStatus(status uint32) {
	*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_LETTER_MASK | (status << kProjectLetterShift))
}

// SetReviewingStatus 设置审核阶段的状态
func (s *ProjectStatus) SetReviewingStatus(status uint32) {
	*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_REVIEW_MASK | (status << kProjectReviewShift))
}

// SetPublishedStatus 设置发布阶段的状态
func (s *ProjectStatus) SetPublishedStatus(status bool) {
	if status {
		*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_PUBLISH_MASK | (PROJECT_STATUS_COMPLETED << kProjectPublishShift))
		return
	}
	*s = ProjectStatus(uint32(*s)&^PROJECT_STATUS_PUBLISH_MASK | (PROJECT_STATUS_UNSET << kProjectPublishShift))
}

// ProjectListParams 定义了获取项目列表的查询参数
type ProjectListParams struct {
	Offset    int             // 偏移量
	Limit     int             // 限制数量
	WorksetId uint            // 作品集 ID，为 0 时忽略
	Index     uint            // 作品集索引，为 0 时忽略
	MemberId  uint            // 成员 ID，为 0 时忽略
	Status    []ProjectStatus // 项目状态，默认值时忽略
}

// ProjectInfo 定义了项目的基本信息
type ProjectInfo struct {
	Id            uint          // 项目 ID
	Title         string        // 项目标题
	Description   string        // 项目简介
	WorksetId     uint          // 所属作品集 ID
	WorksetIndex  uint          // 作品集索引
	LegacyId      uint          // 旧版 ID
	MoetranId     string        // 龙译 ID
	Status        ProjectStatus // 项目状态
	IsPublished   bool          // 是否已发布
	AllowAutoJoin bool          // 是否允许自动加入
	Labors        []*MemberInfo // 成员参与的分工信息
}

// CreateProjectParams 定义了创建项目所需的信息
type CreateProjectParams struct {
	Title       string // 项目标题
	Description string // 项目简介
	WorksetId   uint   // 作品集 ID

	CreatorUserId uint // 创建者的用户 ID
	AllowAutoJoin bool // 允许加入的权限
	IsHidden      bool // 是否隐藏项目
}

// ProjectService 接口定义了项目服务的基本操作
type ProjectService interface {
	// GetProjects 获取指定条件下的项目列表
	GetProjects(params *ProjectListParams) ([]*ProjectInfo, error)
	// GetProjectDetail 获取指定项目的详细信息
	GetProjectDetail(projectId uint) (*ProjectInfo, error)

	// // GetBasicPageWithParams 获取项目列表，支持分页和排序以及复合条件查询
	// GetBasicPageWithParams(
	// 	worksetId, index uint, userId uint,
	// 	pageSerial, pageSize int, sort int,
	// 	status dtos.ProjectOverallStatus,
	// ) ([]*dtos.ProjectBasic, error)

	// // GetDetail 获取指定 ID 的项目详情
	// GetDetail(projectId uint) (*dtos.ProjectDetail, error)

	// // GetLaborDivision 获取指定项目 ID 的团队成员分工
	// GetLaborDivision(projectId uint) ([]*dtos.LaborDivision, error)

	// // CreateProject 创建一个新的项目
	// CreateProject(request *dtos.CreateProjectParams) (*dtos.ProjectCreatedInfo, error)

	// // UpdateProject 更新指定 ID 的项目
	// UpdateProject(projectId uint, request *dtos.UpdateProjectRequest) error

	// // DeleteProject 删除指定 ID 的项目
	// DeleteProject(projectId uint) error
}

// projectServiceImpl 是 ProjectService 的实现
type projectServiceImpl struct {
	handle    *gorm.DB
	apiClient apiclient.ApiClient
	logger    *slog.Logger
}

// NewProjectService 创建一个新的 ProjectService 实例
func NewProjectService(
	hdl *gorm.DB,
	apiClient apiclient.ApiClient,
	logger *slog.Logger,
) ProjectService {
	return &projectServiceImpl{
		handle:    hdl,
		apiClient: apiClient,
		logger:    logger,
	}
}

// GetBasicPageWithParams 实现 ProjectService 接口的 GetBasicPageWithParams 方法
func (s *projectServiceImpl) GetProjects(params *ProjectListParams) ([]*ProjectInfo, error) {
	// 当指定的 userId 不为 0 时，使用 labor 进行查询
	// 其他的时候使用 project 进行查询
	switch params.MemberId {
	case 0:
		// 初步构建查询参数
		projectSpec := s.buildQueryParams(params)
		// 查询的字段
		projectFields := &models.ProjectFields{
			Id:            true,
			Title:         true,
			WorksetId:     true,
			WorksetIndex:  true,
			LegacyId:      true,
			MoetranId:     true,
			Status:        true,
			AllowAutoJoin: true,
		}

		// 直接基于 project 查询
		projects, err := models.GetProject().SelectMany(
			s.handle, projectSpec, projectFields,
			&params.Offset, &params.Limit)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("GetProjects 查询没有结果", slog.Any("params", params))
				return nil, errors.New("没有找到符合条件的项目")
			}
			s.logger.Error("GetProjects 查询项目列表时出现错误", slog.Any("error", err))
			return nil, fmt.Errorf("查询项目列表时出现错误: %w", err)
		}

		// 将查询结果转换为 ProjectInfo
		projectInfos := make([]*ProjectInfo, 0, len(projects))
		for _, p := range projects {
			projectInfos = append(projectInfos, &ProjectInfo{
				Id:            uint(p.Id),
				Title:         p.Title,
				WorksetId:     uint(p.WorksetId),
				WorksetIndex:  p.WorksetIndex,
				LegacyId:      p.LegacyId,
				MoetranId:     p.MoetranId,
				Status:        s.buildProjectStatus(p),
				IsPublished:   p.IsPublished,
				AllowAutoJoin: p.AllowAutoJoin,
			})
		}

		return projectInfos, nil

	default:
		// 将 Labor、Project、Member、User 关联查询
		type result struct {
			UserId          uint
			Nickname        string
			MemberId        uint
			ProjectId       uint
			ProjectTitle    string
			WorksetId       uint
			WorksetIndex    uint
			LegacyId        uint
			MoetranId       string
			IsPublished     bool
			TranslateStatus uint8
			ProofreadStatus uint8
			LetterStatus    uint8
			ReviewStatus    uint8
			AllowAutoJoin   bool
		}

		var results []result

		query := s.handle.Model(&models.Labor{}).
			// LEFT JOIN project 表
			Joins("LEFT JOIN projects ON labors.project_id = projects.id").
			// LEFT JOIN member 表
			Joins("LEFT JOIN members ON labors.member_id = members.id").
			// LEFT JOIN user 表 (通过 members 表连接)
			Joins("LEFT JOIN users ON members.user_id = users.id").
			// 筛选条件：根据 member_id
			Where("members.id = ?", params.MemberId)

		// 为 query 加上要筛选的 project 字段
		s.buildQueryParams(params).Apply(query)

		if err := query.
			// 选择所有需要的列
			Select(`
            	users.id AS user_id,
            	users.nickname AS nickname,
            	members.id AS member_id,
            	projects.id AS project_id,
				projects.title AS project_title,
				projects.workset_id AS workset_id,
				projects.workset_index AS workset_index,
				projects.legacy_id AS legacy_id,
				projects.moetran_id AS moetran_id,
				projects.is_published AS is_published,
				projects.translate_status AS translate_status,
				projects.proofread_status AS proofread_status,
				projects.letter_status AS letter_status,
				projects.review_status AS review_status,
				projects.allow_auto_join AS allow_auto_join
			`).
			Scan(&results).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("GetProjects 查询没有结果", slog.Any("params", params))
				return nil, errors.New("没有找到符合条件的项目")
			}

			s.logger.Error("GetProjects 查询项目列表时出现错误", slog.Any("error", err))
			return nil, errors.New("查询项目列表失败")
		}

		// 将查询结果转换为 ProjectInfo
		projectInfos := make([]*ProjectInfo, 0, len(results))

		for _, r := range results {
			var status ProjectStatus
			status.SetTranslatingStatus(uint32(r.TranslateStatus))
			status.SetProofreadingStatus(uint32(r.ProofreadStatus))
			status.SetLetteringStatus(uint32(r.LetterStatus))
			status.SetReviewingStatus(uint32(r.ReviewStatus))
			status.SetPublishedStatus(r.IsPublished)

			projectInfos = append(projectInfos, &ProjectInfo{
				Id:            r.ProjectId,
				Title:         r.ProjectTitle,
				WorksetId:     r.WorksetId,
				WorksetIndex:  r.WorksetIndex,
				LegacyId:      r.LegacyId,
				MoetranId:     r.MoetranId,
				Status:        status,
				IsPublished:   r.IsPublished,
				AllowAutoJoin: r.AllowAutoJoin,
				Labors: []*MemberInfo{
					{
						Id: r.MemberId,
						User: UserInfo{
							Id:       r.UserId,
							Nickname: r.Nickname,
						},
					},
				},
			})
		}

		return projectInfos, nil
	}
}

// GetProjectDetail 实现 ProjectService 接口的 GetProjectDetail 方法
func (s *projectServiceImpl) GetProjectDetail(projectId uint) (*ProjectInfo, error) {
	// 查询项目的基本信息
	projectPKey := models.PKey(projectId)
	projectSpec := &models.ProjectSpec{
		Id: &projectPKey,
	}

	project, err := models.GetProject().SelectFirst(s.handle, projectSpec, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("GetProjectDetail 查询没有结果", slog.Uint64("projectId", uint64(projectId)))
			return nil, errors.New("没有找到对应的项目")
		}
		s.logger.Error("GetProjectDetail 查询项目详情失败", slog.Any("error", err))
		return nil, fmt.Errorf("查询项目详情失败")
	}

	projectInfo := &ProjectInfo{
		Id:            uint(project.Id),
		Title:         project.Title,
		Description:   project.Description,
		WorksetId:     uint(project.WorksetId),
		WorksetIndex:  project.WorksetIndex,
		LegacyId:      project.LegacyId,
		MoetranId:     project.MoetranId,
		Status:        s.buildProjectStatus(project),
		IsPublished:   project.IsPublished,
		AllowAutoJoin: project.AllowAutoJoin,
	}

	return projectInfo, nil
}

// // CreateProject 实现 ProjectService 接口的 CreateProject 方法
// func (s *projectServiceImpl) CreateProject(createInfo *dtos.CreateProjectParams) (*dtos.ProjectCreatedInfo, error) {
// 	// 构建一个事务协调器
// 	coordinater := txutils.NewTransactionCoordinator(s.projectRepo.GetHandle())

// 	var createdInfo *dtos.ProjectCreatedInfo

// 	// 在一个事务中协调数据库和调用龙译 API 的操作
// 	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
// 		// 基于 transaction 上下文获取仓库实例
// 		txProjectRepo := repos.NewProjectRepo(tx)
// 		txLaborRepo := repos.NewLaborRepo(tx)

// 		// 创建一个新的项目
// 		project := &dbmodels.Project{
// 			Title:         createInfo.Title,
// 			Description:   createInfo.Description,
// 			WorksetId:     dbmodels.PrimaryKey(createInfo.WorksetId),
// 			AllowAutoJoin: createInfo.AllowAutoJoin,
// 			IsHidden:      createInfo.IsHidden,
// 			PrincipalId:   dbmodels.PrimaryKey(createInfo.CreatorUserId), // 设置负责人为创建者
// 		}

// 		// 在数据库中创建项目，如果成功 project 的 Id、WorksetIndex 和 FkWorkset 应当被填充
// 		if err := txProjectRepo.CreateProject(project); err != nil {
// 			s.logger.Error("CreateProject 调用 CreateProject 中出现错误", slog.Any("error", err))
// 			return errors.New("创建项目失败"), nil
// 		}

// 		// 创建 creator 的分工记录
// 		creatorLabor := dbmodels.LaborMask(0)
// 		creatorLabor.AddRole(dbmodels.LABOR_CREATOR_MASK)
// 		creatorLabor.AddRole(dbmodels.LABOR_PRINCIPAL_MASK)

// 		laborDivision := &dbmodels.ProjectLaborDivision{
// 			ProjectId: dbmodels.PrimaryKey(project.Id),
// 			UserId:    dbmodels.PrimaryKey(createInfo.CreatorUserId),
// 			LaborRole: creatorLabor,
// 		}

// 		// 然后在团队成员分工表中插入 creator 记录
// 		if err := txLaborRepo.CreateLaborDivision(laborDivision); err != nil {
// 			s.logger.Error("CreateProject 调用 CreateLaborDivision 中出现错误", slog.Any("error", err))
// 			return fmt.Errorf("创建团队成员分工失败"), nil
// 		}

// 		// 获取对应 user 的龙译 JWT
// 		user, err := s.userRepo.SelectByUserId(dbmodels.PrimaryKey(createInfo.CreatorUserId))
// 		if err != nil {
// 			s.logger.Error("CreateProject 获取用户信息失败", slog.Any("error", err))
// 			return fmt.Errorf("获取用户信息失败"), nil
// 		}

// 		// 随后调用龙译 API 创建项目
// 		projInfo := s.buildMoetranProjInfo(createInfo, project, user)

// 		moetranRes, err := s.apiClient.CreateProject(projInfo)
// 		if err != nil {
// 			s.logger.Error("CreateProject 调用 CreateProject API 中出现错误", slog.Any("error", err))
// 			// TODO：由于不确定尨译的 API 是否是幂等的，这里需要考虑补偿操作，比如删除对应项目
// 			// 但在不确定尨译实现的情况下，先不处理补偿
// 			return errors.New("调用龙译 API 创建项目失败"), nil
// 		}

// 		// 将龙译返回的项目 ID 更新到本地项目中
// 		if err := txProjectRepo.UpdateMoetranId(project.Id, moetranRes.Project.Id); err != nil {
// 			s.logger.Error("CreateProject 更新本地项目龙译 ID 失败", slog.Any("error", err))
// 			return fmt.Errorf("更新本地龙译 ID 失败"), nil
// 		}

// 		// 写入到返回结果
// 		createdInfo = &dtos.ProjectCreatedInfo{
// 			Message:   moetranRes.Message,
// 			ProjectId: uint(project.Id),
// 			MoetranId: moetranRes.Project.Id,
// 		}

// 		// 一切正常，则返回 nil 提交事务
// 		return nil, nil

// 	}); err != nil {
// 		// 这里的 err 是上述事务中抛出的错误
// 		s.logger.Error("CreateProject 事务执行失败", slog.Any("error", err))
// 		return nil, errors.New("创建项目失败")
// 	}

// 	return createdInfo, nil
// }

// ================ 辅助函数 ================

// buildProjectSpec 通过 service 参数为 projectSpec 构建查询条件
func (srv *projectServiceImpl) buildQueryParams(params *ProjectListParams) *models.ProjectSpec {
	var projectSpec models.ProjectSpec

	if params.WorksetId > 0 {
		// 如果指定了作品集 ID，则添加查询条件
		worksetPKey := models.PKey(params.WorksetId)
		projectSpec.WorksetId = &worksetPKey
	}

	if params.Index > 0 {
		// 如果指定了作品集索引，则添加查询条件
		projectSpec.WorksetIndex = &params.Index
	}

	if len(params.Status) != 0 {
		// 如果传入了 Status 参数，则添加状态条件
		for _, status := range params.Status {
			if status == 0 {
				// 如果 s 是默认值，说明有问题，直接跳过
				continue
			}

			// 当 Status 不是全默认时，解析位掩码
			ps := ProjectStatus(status)

			if s := ps.GetTranslatingStatus(); s != PROJECT_STATUS_NOT_USED {
				if projectSpec.TranslateStatus == nil {
					slice := make([]uint8, 0)
					projectSpec.TranslateStatus = &slice
				}
				*projectSpec.TranslateStatus = append(*projectSpec.TranslateStatus, uint8(s))
			}
			if s := ps.GetProofreadingStatus(); s != PROJECT_STATUS_NOT_USED {
				if projectSpec.ProofreadStatus == nil {
					slice := make([]uint8, 0)
					projectSpec.ProofreadStatus = &slice
				}
				*projectSpec.ProofreadStatus = append(*projectSpec.ProofreadStatus, uint8(s))
			}
			if s := ps.GetLetteringStatus(); s != PROJECT_STATUS_NOT_USED {
				if projectSpec.LetterStatus == nil {
					slice := make([]uint8, 0)
					projectSpec.LetterStatus = &slice
				}
				*projectSpec.LetterStatus = append(*projectSpec.LetterStatus, uint8(s))
			}
			if s := ps.GetReviewingStatus(); s != PROJECT_STATUS_NOT_USED {
				if projectSpec.ReviewStatus == nil {
					slice := make([]uint8, 0)
					projectSpec.ReviewStatus = &slice
				}
				*projectSpec.ReviewStatus = append(*projectSpec.ReviewStatus, uint8(s))
			}
			if s := ps.GetPublishedStatus(); s != PROJECT_STATUS_NOT_USED {
				switch s {
				case PROJECT_STATUS_COMPLETED:
					b := true
					projectSpec.IsPublished = &b
				case PROJECT_STATUS_UNSET:
					b := false
					projectSpec.IsPublished = &b
				default:
					// 跳过未知状态
					srv.logger.Warn("buildQueryParams 未知的项目状态", slog.Any("status", s))
					continue
				}
			}
		}
	}

	return &projectSpec
}

// buildProjectStatus 根据传入的状态参数构建 ProjectStatus
func (s *projectServiceImpl) buildProjectStatus(model *models.Project) ProjectStatus {
	var status ProjectStatus

	status.SetTranslatingStatus(uint32(model.TranslateStatus))
	status.SetProofreadingStatus(uint32(model.ProofreadStatus))
	status.SetLetteringStatus(uint32(model.LetterStatus))
	status.SetReviewingStatus(uint32(model.ReviewStatus))
	status.SetPublishedStatus(model.IsPublished)

	return status
}

// buildMoetranProjInfo 构建 Moetran 项目信息
func (s *projectServiceImpl) buildMoetranProjInfo(
	info *CreateProjectParams,
	project *models.Project,
	user *models.User,
) *apiclient.CreateProjectInfo {
	// 组装尨译的 title
	title := fmt.Sprintf("[%d-%d]%s",
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
		MoetranAuth: user.MoetranJwt,

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
