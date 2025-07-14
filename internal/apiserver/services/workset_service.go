package services

import (
	"log/slog"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
)

// WorksetService 接口定义了作品集服务的基本操作
type WorksetService interface {
	// GetBasicPageIdDesc 获取作品集列表，按 ID 倒序
	GetBasicPage(pageSerial, pageSize int) ([]*dtos.WorksetBasic, error)

	// GetProjectStats 获取特定作品集的项目统计信息
	GetProjectStats(worksetId uint) (*dtos.ProjectStats, error)
}

// worksetServiceImpl 是 WorksetService 的实现
type worksetServiceImpl struct {
	worksetRepo  repos.WorksetRepo
	materialView repos.MaterialView
	logger       *slog.Logger
}

// NewWorksetService 创建一个新的 WorksetService 实例
func NewWorksetService(
	worksetRepo repos.WorksetRepo,
	materialView repos.MaterialView,
	logger *slog.Logger,
) WorksetService {
	return &worksetServiceImpl{
		worksetRepo:  worksetRepo,
		materialView: materialView,
		logger:       logger,
	}
}

// GetBasicPage 实现 WorksetService 接口的 GetBasicPage 方法
func (s *worksetServiceImpl) GetBasicPage(pageSerial, pageSize int) ([]*dtos.WorksetBasic, error) {
	// 调用仓库方法获取作品集列表
	worksets, err := s.worksetRepo.SelectBasicPageIdDesc((pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetBasicPage 调用 SelectBasicPageIdDesc 中出现错误", slog.Any("error", err))
		return nil, err
	}

	// 将 dbmodels.Workset 转换为 dtos.WorksetBasic
	var worksetBasics []*dtos.WorksetBasic
	for _, workset := range worksets {
		worksetBasics = append(worksetBasics, &dtos.WorksetBasic{
			Id:     uint(workset.Id),
			Name:   workset.Name,
			TeamId: uint(workset.TeamId),
		})
	}

	return worksetBasics, nil
}

// GetProjectStats 实现 WorksetService 接口的 GetProjectStats 方法
func (s *worksetServiceImpl) GetProjectStats(worksetId uint) (*dtos.ProjectStats, error) {
	// 调用仓库方法获取项目统计信息
	stats, err := s.materialView.SelectProjectStats(dbmodels.PrimaryKey(worksetId))
	if err != nil {
		s.logger.Error("GetProjectStats 调用 GetProjectStats 中出现错误", slog.Any("error", err))
		return nil, err
	}

	// 将 dbmodels.ProjectStats 转换为 dtos.ProjectStats
	projectStats := &dtos.ProjectStats{
		WorksetId:      uint(stats.WorksetId),
		TotalCount:     stats.Total,
		PublishedCount: stats.Published,

		NotTranslatingCount: stats.NotTranslating,
		TranslatingCount:    stats.Translating,
		TranslatedCount:     stats.Translated,

		NotProovingCount: stats.NotProoving,
		ProovingCount:    stats.Prooving,
		ProovedCount:     stats.Prooved,

		NotLetteringCount: stats.NotLettering,
		LetteringCount:    stats.Lettering,
		LetteredCount:     stats.Letterred,

		NotReviewingCount: stats.NotReviewing,
		ReviewingCount:    stats.Reviewing,
		ReviewedCount:     stats.Reviewed,
	}

	return projectStats, nil
}
