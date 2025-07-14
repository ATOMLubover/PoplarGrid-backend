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
	"poplargrid/internal/shared/transaction"

	"gorm.io/gorm"
)

// ProjectService 接口定义了项目服务的基本操作
type ProjectService interface {
	// GetBasicPageWithParams 获取项目列表，支持分页和排序以及复合条件查询
	GetBasicPageWithParams(worksetId uint, pageSerial, pageSize int, sort int, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error)
	// GetBasicPageByUserId 获取指定用户参与的项目列表，支持分页
	GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.ProjectBasic, error)

	// GetDetailById 获取指定 ID 的项目详情
	GetDetailById(id uint) (*dtos.ProjectDetail, error)

	// GetLaborDivisionByProjectId 获取指定项目 ID 的团队成员分工
	GetLaborDivisionByProjectId(projectId uint) ([]*dtos.LaborDivision, error)

	// CreateProject 创建一个新的项目，如果成功返回新主键和 nil，否则返回错误
	CreateProject(request *dtos.CreateProjectRequest) (uint, error)
}

// projectServiceImpl 是 ProjectService 的实现
type projectServiceImpl struct {
	projectRepo repos.ProjectRepo
	laborRepo   repos.LaborRepo
	apiClient   apiclient.ApiClient
	logger      *slog.Logger
}

// NewProjectService 创建一个新的 ProjectService 实例
func NewProjectService(
	projectRepo repos.ProjectRepo,
	laborRepo repos.LaborRepo,
	apiClient apiclient.ApiClient,
	logger *slog.Logger,
) ProjectService {
	return &projectServiceImpl{
		projectRepo: projectRepo,
		laborRepo:   laborRepo,
		apiClient:   apiClient,
		logger:      logger,
	}
}

// GetBasicPageWithParams 实现 ProjectService 接口的 GetBasicPageWithParams 方法
func (s *projectServiceImpl) GetBasicPageWithParams(worksetId uint, pageSerial, pageSize int, sort int, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error) {
	// 解析 status 位掩码，获取各个阶段的状态
	queryParams := s.buildQueryParams(status)

	// 根据 sort 选择不同的 service 函数处理
	var page []*dbmodels.Project
	var err error

	switch sort {
	case dtos.SORT_ID_DESC:
		page, err = s.projectRepo.SelectBasicPageIdDescWithParam(dbmodels.PrimaryKey(worksetId), (pageSerial-1)*pageSize, pageSize, queryParams)
		if err != nil {
			s.logger.Error("GetBasicPage 调用 SelectBasicPageIdDesc 中出现错误", slog.Any("error", err))
			return nil, errors.New("按 ID 降序获取项目列表失败")
		}

	case dtos.SORT_UPDATED_AT_DESC:
		page, err = s.projectRepo.SelectBasicPageUpdatedAtDescWithParam(dbmodels.PrimaryKey(worksetId), pageSerial, pageSize, queryParams)
		if err != nil {
			s.logger.Error("GetBasicPage 调用 SelectBasicPageUpdatedAtDesc 中出现错误", slog.Any("error", err))
			return nil, errors.New("按更新时间降序获取项目列表失败")
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
			Status:        status,
			IsPublished:   project.IsPublished,
			AllowAutoJoin: project.AllowAutoJoin,
		})
	}

	return projectBasics, nil
}

// GetBasicPageByUserId 实现 ProjectService 接口的 GetBasicPageByUserId 方法
func (s *projectServiceImpl) GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.ProjectBasic, error) {
	// 获取用户参与的项目列表
	projects, err := s.laborRepo.SelectProjectBasicPageIdDescByUserId(dbmodels.PrimaryKey(userId), (pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetBasicPageByUserId 调用 SelectBasicPageIdDescByUserId 中出现错误", slog.Any("error", err))
		return nil, fmt.Errorf("获取用户参与的项目列表失败")
	}

	// 将 dbmodels.Project 转换为 dtos.ProjectBasic
	var projectBasics []*dtos.ProjectBasic

	for _, project := range projects {
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
			Status:        status,
			IsPublished:   project.IsPublished,
			AllowAutoJoin: project.AllowAutoJoin,
		})
	}

	return projectBasics, nil
}

// GetDetailById 实现 ProjectService 接口的 GetDetailById 方法
func (s *projectServiceImpl) GetDetailById(id uint) (*dtos.ProjectDetail, error) {
	// 调用仓库方法获取项目详情
	project, err := s.projectRepo.SelectById(dbmodels.PrimaryKey(id))
	if err != nil {
		s.logger.Error("GetDetailById 调用 SelectById 中出现错误", slog.Any("error", err))
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

// ================ 辅助函数 ================

// buildQueryParams 通过 service 参数为 queryParams 构建查询条件
func (s *projectServiceImpl) buildQueryParams(status dtos.ProjectOverallStatus) *dtos.ProjectStatusQueryParams {
	// 解析 status 位掩码，获取各个阶段的状态
	withParams := false

	queryParams := &dtos.ProjectStatusQueryParams{}
	{
		// 翻译状态
		if s := status.GetTranslatingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.TranslateStatus = &s
			withParams = true
		}

		// 校对状态
		if s := status.GetProofreadingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.ProofStatus = &s
			withParams = true
		}

		// 嵌字状态
		if s := status.GetLetteringStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.LetterStatus = &s
			withParams = true
		}

		// 审核状态
		if s := status.GetReviewingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.ReviewStatus = &s
			withParams = true
		}

		// 发布状态要特别处理
		if s := status.GetPublishedStatus(); s != dtos.PROJECT_STATUS_NOT_USED && s != dtos.PROJECT_STATUS_IN_PROGRESS {
			switch s {
			case dtos.PROJECT_STATUS_UNSET:
				isPublished := false
				queryParams.PublishStatus = &isPublished
			case dtos.PROJECT_STATUS_COMPLETED:
				isPublished := true
				queryParams.PublishStatus = &isPublished
			}
			withParams = true
		}
	}

	// 如果没有有效查询条件，退化为无条件查询
	if !withParams {
		queryParams = nil
	}

	return queryParams
}

// GetLaborDivisionByProjectId 实现 ProjectService 接口的 GetLaborDivisionByProjectId 方法
func (s *projectServiceImpl) GetLaborDivisionByProjectId(projectId uint) ([]*dtos.LaborDivision, error) {
	// 调用仓库方法获取团队成员分工
	members, err := s.laborRepo.SelectByProjectId(dbmodels.PrimaryKey(projectId))
	if err != nil {
		s.logger.Error("GetLaborDivisionByProjectId 调用 SelectLaborDivisionByProjectId 中出现错误", slog.Any("error", err))
		return nil, fmt.Errorf("获取项目成员分工失败")
	}

	// 将 dbmodels.TeamMember 转换为 dtos.LaborDivision
	var divisions []*dtos.LaborDivision
	for _, member := range members {
		divisions = append(divisions, &dtos.LaborDivision{
			MemberId: uint(member.Id),
			Nickname: member.FkMember.FkUser.Nickname,
			Role:     uint(member.LaborRole),
		})
	}

	return divisions, nil
}

// CreateProject 实现 ProjectService 接口的 CreateProject 方法
func (s *projectServiceImpl) CreateProject(request *dtos.CreateProjectRequest) (uint, error) {
	// 创建一个将要插入的 dbmodels.Project 实例
	project := dbmodels.Project{
		Title:       request.Title,
		Description: request.Description,
		WorksetId:   dbmodels.PrimaryKey(request.WorksetId),
	}

	// 构建一个事务协调器
	coordinater := transaction.NewTransactionCoordinator(s.projectRepo.GetHandle())

	// 在一个事务中协调数据库和调用龙译 API 的操作
	if err := coordinater.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 先在数据库中创建项目，如果成功 project 的 Id 和 WorksetIndex 应当被填充
		if err := s.projectRepo.CreateProject(&project); err != nil {
			s.logger.Error("CreateProject 调用 CreateProject 中出现错误", slog.Any("error", err))
			return fmt.Errorf("创建项目失败：%w", err), nil
		}

		// 随后调用龙译 API 创建项目
		if err := s.apiClient.CreateProject(request, project.WorksetIndex); err != nil {
			s.logger.Error("CreateProject 调用 CreateProject API 中出现错误", slog.Any("error", err))
			// TODO：由于不确定尨译的 API 是否是幂等的，这里需要考虑补偿操作，比如删除对应项目
			// 但在不确定尨译实现的情况下，先不处理补偿
			return fmt.Errorf("调用龙译 API 创建项目失败：%w", err), nil
		}

		// 一切正常，则返回 nil 提交事务
		return nil, nil

	}); err != nil {
		// 这里的 err 是上述事务中抛出的错误
		s.logger.Error("CreateProject 事务执行失败", slog.Any("error", err))
		return 0, errors.New("创建项目失败")
	}

	return uint(project.Id), nil
}
