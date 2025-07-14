package services

import (
	"errors"
	"fmt"
	"log/slog"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
)

// ProjectService 接口定义了项目服务的基本操作
type ProjectService interface {
	// GetBasicPage 获取所有项目列表，支持分页和排序，包括发布和未发布的项目
	GetBasicPage(worksetId uint, pageSerial, pageSize int, sort int) ([]*dtos.ProjectBasic, error)
	// GetBasicPageWithParams 获取项目列表，支持分页和排序以及复合条件查询
	GetBasicPageWithParams(worksetId uint, pageSerial, pageSize int, sort int, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error)
}

// projectServiceImpl 是 ProjectService 的实现
type projectServiceImpl struct {
	projectRepo repos.ProjectRepo
	logger      *slog.Logger
}

// NewProjectService 创建一个新的 ProjectService 实例
func NewProjectService(
	projectRepo repos.ProjectRepo,
	logger *slog.Logger,
) ProjectService {
	return &projectServiceImpl{
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// GetBasicPage 实现 ProjectService 接口的 GetBasicPage 方法
func (s *projectServiceImpl) GetBasicPage(worksetId uint, pageSerial, pageSize int, sort int) ([]*dtos.ProjectBasic, error) {
	// 调用仓库方法获取项目列表
	var page []*dbmodels.Project
	var err error

	switch sort {
	case dtos.SORT_ID_DESC:
		page, err = s.projectRepo.SelectBasicPageIdDesc(dbmodels.PrimaryKey(worksetId), (pageSerial-1)*pageSize, pageSize)
		if err != nil {
			s.logger.Error("GetBasicPage 调用 SelectBasicPageIdDesc 中出现错误", slog.Any("error", err))
			return nil, errors.New("按 ID 降序获取项目列表失败")
		}

	case dtos.SORT_UPDATED_AT_DESC:
		page, err = s.projectRepo.SelectBasicPageUpdatedAtDesc(dbmodels.PrimaryKey(worksetId), (pageSerial-1)*pageSize, pageSize)
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

// GetBasicPageWithParams 实现 ProjectService 接口的 GetBasicPageWithParams 方法
func (s *projectServiceImpl) GetBasicPageWithParams(worksetId uint, pageSerial, pageSize int, sort int, status dtos.ProjectOverallStatus) ([]*dtos.ProjectBasic, error) {
	// 解析 status 位掩码，获取各个阶段的状态
	queryParams := &dtos.ProjectStatusQueryParams{}
	{
		if s := status.GetTranslatingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.TranslateStatus = &s
		}
		if s := status.GetProofreadingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.ProofStatus = &s
		}
		if s := status.GetLetteringStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.LetterStatus = &s
		}
		if s := status.GetReviewingStatus(); s != dtos.PROJECT_STATUS_NOT_USED {
			queryParams.ReviewStatus = &s
		}
	}

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
