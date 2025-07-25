package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// WorksetService 接口定义了作品集服务的基本操作
type WorksetService interface {
	// GetWorksets 根据参数获取作品集列表
	GetWorksets(params *WorksetListParams) ([]*WorksetInfo, error)

	// GetWorksetStatsByWorksetId 获取特定作品集的项目统计信息
	GetWorksetStatsByWorksetId(worksetId uint) (*WorksetStats, error)
}

// WorksetListParams 定义了获取作品集列表的查询参数
type WorksetListParams struct {
	Offset int  // 偏移量
	Limit  int  // 限制数量
	TeamId uint // 团队 ID
}

// WorksetInfo 定义了作品集的基本信息
type WorksetInfo struct {
	Id        uint   // 作品集 ID
	Name      string // 作品集名称
	TeamId    uint   // 所属团队 ID
	MoetranId string // 龙译 ID
}

// WorksetStats 定义了项目统计信息
type WorksetStats struct {
	WorksetId      uint // 作品集 ID
	TotalCount     int  // 总项目数
	PublishedCount int  // 已发布项目数

	NotTranslatingCount int // 未翻译项目数
	TranslatingCount    int // 正在翻译项目数
	TranslatedCount     int // 已翻译项目数

	NotProovingCount int // 未校对项目数
	ProovingCount    int // 正在校对项目数
	ProovedCount     int // 已校对项目数

	NotLetteringCount int // 未嵌字项目数
	LetteringCount    int // 正在嵌字项目数
	LetteredCount     int // 已嵌字项目数

	NotReviewingCount int // 未审核项目数
	ReviewingCount    int // 正在审核项目数
	ReviewedCount     int // 已审核项目数
}

// worksetServiceImpl 是 WorksetService 的实现
type worksetServiceImpl struct {
	handle *gorm.DB
	logger *slog.Logger
}

// NewWorksetService 创建一个新的 WorksetService 实例
func NewWorksetService(
	hdl *gorm.DB,
	lgr *slog.Logger,
) WorksetService {
	return &worksetServiceImpl{
		handle: hdl,
		logger: lgr,
	}
}

// GetWorksets 实现 WorksetService 接口的 GetWorksets 方法
func (s *worksetServiceImpl) GetWorksets(params *WorksetListParams) ([]*WorksetInfo, error) {
	// 查询条件为汉化组 ID
	teamPKey := models.PKey(params.TeamId)
	worksetSpec := &models.WorksetSpec{
		TeamId: &teamPKey,
	}

	// 执行查询
	worksets, err := models.GetWorkset().SelectMany(
		s.handle, worksetSpec,
		&params.Offset, &params.Limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("没有找到指定团队的作品集",
				slog.Any("team_id", params.TeamId),
				slog.Any("error", err))
			return nil, errors.New("没有找到指定团队的作品集")
		}
		s.logger.Error("获取指定团队的作品集失败",
			slog.Any("team_id", params.TeamId),
			slog.Any("error", err))
		return nil, errors.New("查询团队的作品集列表失败")
	}

	// 将查询结果转换为 WorksetInfo
	worksetInfos := make([]*WorksetInfo, 0, len(worksets))

	for _, w := range worksets {
		worksetInfos = append(worksetInfos, &WorksetInfo{
			Id:        uint(w.Id),
			Name:      w.Name,
			TeamId:    uint(w.TeamId),
			MoetranId: w.MoetranId,
		})
	}

	return worksetInfos, nil
}

// GetWorksetStatsByWorksetId 实现 WorksetService 接口的 GetWorksetStatsByWorksetId 方法
func (s *worksetServiceImpl) GetWorksetStatsByWorksetId(worksetId uint) (*WorksetStats, error) {
	// 查询条件为作品集 ID
	worksetPKey := models.PKey(worksetId)

	// 执行对物化视图的查询
	stat, err := models.GetWorkStatsMv().Select(s.handle, worksetPKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("没有找到指定的作品集",
				slog.Any("workset_id", worksetPKey),
				slog.Any("error", err))
			return nil, errors.New("没有找到指定的作品集")
		}
		s.logger.Error("获取指定的作品集统计数据失败",
			slog.Any("workset_id", worksetPKey),
			slog.Any("error", err))
		return nil, errors.New("查询指定的作品集统计数据失败")
	}

	// 转化为 WorksetStats
	statsInfo := &WorksetStats{
		WorksetId:      uint(stat.WorksetId),
		TotalCount:     stat.Total,
		PublishedCount: stat.Published,

		NotTranslatingCount: stat.NotTranslating,
		TranslatingCount:    stat.Translating,
		TranslatedCount:     stat.Translated,

		NotProovingCount: stat.NotProoving,
		ProovingCount:    stat.Prooving,
		ProovedCount:     stat.Prooved,

		NotLetteringCount: stat.NotLettering,
		LetteringCount:    stat.Lettering,
		LetteredCount:     stat.Letterred,

		NotReviewingCount: stat.NotReviewing,
		ReviewingCount:    stat.Reviewing,
		ReviewedCount:     stat.Reviewed,
	}

	return statsInfo, nil
}
