package services

import (
	"context"
	"errors"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"
	"poplargrid/internal/shared/txutils"

	"gorm.io/gorm"
)

// WorksetListParams 定义了获取作品集列表的查询参数
type WorksetListParams struct {
	Offset int  // 偏移量
	Limit  int  // 限制数量
	TeamId uint // 团队 ID
}

// WorksetInfo 定义了作品集的基本信息
type WorksetInfo struct {
	Id          uint   // 作品集 ID
	Name        string // 作品集名称
	Description string // 作品集描述
	TeamId      uint   // 所属团队 ID
	MoetranId   string // 龙译 ID
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

// CreateWorksetParams 定义了创建工作集的请求参数
type CreateWorksetParams struct {
	Name             string            // 工作集名称
	Description      string            // 工作集描述
	OperatorMemberId uint              // 操作人成员 ID
	CurrentMemberIds map[uint]struct{} // 当前成员 ID 集合
	TeamId           uint              // 所属团队 ID
}

// WorksetCreatedInfo 定义了创建工作集后的返回信息
type WorksetCreatedInfo struct {
	Message   string // 成功消息
	WorksetId uint   // 工作集 ID
	MoetranId string // 龙译 ID
}

// WorksetService 接口定义了作品集服务的基本操作
type WorksetService interface {
	// GetWorksets 根据参数获取作品集列表
	GetWorksets(params *WorksetListParams) ([]*WorksetInfo, error)
	// GetWorksetStats 获取特定作品集的项目统计信息
	GetWorksetStats(worksetId uint) (*WorksetStats, error)

	// CreateWorkset 创建一个新的工作集
	CreateWorkset(params *CreateWorksetParams) (*WorksetCreatedInfo, error)
}

// worksetServiceImpl 是 WorksetService 的实现
type worksetServiceImpl struct {
	handle    *gorm.DB
	apiClient apiclient.ApiClient // 用于与龙译 API 交互
	logger    *slog.Logger
}

// NewWorksetService 创建一个新的 WorksetService 实例
func NewWorksetService(
	hdl *gorm.DB,
	apiClient apiclient.ApiClient,
	lgr *slog.Logger,
) WorksetService {
	return &worksetServiceImpl{
		handle:    hdl,
		apiClient: apiClient,
		logger:    lgr,
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
			Id:          uint(w.Id),
			Name:        w.Name,
			Description: w.Description,
			TeamId:      uint(w.TeamId),
			MoetranId:   w.MoetranId,
		})
	}

	return worksetInfos, nil
}

// GetWorksetStats 实现 WorksetService 接口的 GetWorksetStats 方法
func (s *worksetServiceImpl) GetWorksetStats(worksetId uint) (*WorksetStats, error) {
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

// CreateWorkset 实现 WorksetService 接口的 CreateWorkset 方法
func (s *worksetServiceImpl) CreateWorkset(params *CreateWorksetParams) (*WorksetCreatedInfo, error) {
	// 检查是否是非法冒用
	// 检查申请者的成员 ID 是否在上下文中
	if _, exists := params.CurrentMemberIds[params.OperatorMemberId]; !exists {
		s.logger.Warn("CreateWorkset 所使用成员 ID 不在当前用户的成员列表中",
			slog.Uint64("operator_member_id", uint64(params.OperatorMemberId)))
		return nil, errors.New("非法引用申请者成员 ID")
	}

	// 检查当前用户是否有权利在本团队下创建工作集
	memberPKey := models.PKey(params.OperatorMemberId)
	memberSpec := &models.MemberSpec{
		Id: &memberPKey,
	}
	memberFields := &models.MemberFields{
		UserId:  true,
		IsAdmin: true, // 需要检查是否是管理员
		UserFields: &models.UserFields{
			MoetranJwt: true, // 需要获取用户的 JWT
		},
	}
	// 执行查询
	member, err := models.GetMember().SelectFirst(s.handle, memberSpec, memberFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("CreateWorkset 没有找到指定成员",
				slog.Any("member_id", params.OperatorMemberId),
				slog.Any("error", err))
			return nil, errors.New("没有找到指定成员")
		}
		s.logger.Error("CreateWorkset 查询指定成员失败",
			slog.Any("member_id", params.OperatorMemberId),
			slog.Any("error", err))
		return nil, errors.New("查询指定成员失败")
	}

	// 如果当前用户不是团队的管理员，则不允许创建工作集
	if !member.IsAdmin {
		s.logger.Warn("CreateWorkset 当前用户不是团队管理员",
			slog.Any("member_id", params.OperatorMemberId),
			slog.Any("team_id", params.TeamId))
		return nil, errors.New("当前用户没有权限在本团队下创建工作集")
	}

	// 获取对应汉化组的龙译 ID
	teamPKey := models.PKey(params.TeamId)
	teamSpec := &models.TeamSpec{
		Id: &teamPKey,
	}

	team, err := models.GetTeam().SelectFirst(s.handle, teamSpec)
	if err != nil {
		s.logger.Error("CreateWorkset 查询团队信息失败",
			slog.Any("team_id", params.TeamId),
			slog.Any("error", err))
		return nil, errors.New("查询团队信息失败")
	}

	// 创建事务协调器
	c := txutils.NewTransactionCoordinator(s.handle)

	var worksetCreatedInfo *WorksetCreatedInfo

	if err := c.RunInTransaction(context.Background(), func(tx *gorm.DB) (error, func() error) {
		// 先在本地创建工作集
		workset := &models.Workset{
			Name:        params.Name,
			Description: params.Description,
			TeamId:      models.PKey(params.TeamId),
		}

		if err := models.GetWorkset().Insert(tx, workset); err != nil {
			s.logger.Error("CreateWorkset 插入工作集失败",
				slog.Any("workset", workset),
				slog.Any("error", err))
			return errors.New("插入工作集失败"), nil
		}

		// 调用龙译创建作品集
		projSetInfo, err := s.apiClient.CreateProjectSet(&apiclient.CreateProjectSetParams{
			Name:          params.Name,
			MoetranAuth:   member.FkUser.MoetranJwt,
			MoetranTeamId: team.MoetranId,
		})
		if err != nil {
			s.logger.Error("CreateWorkset 调用龙译 API 创建作品集失败",
				slog.Any("team_id", params.TeamId),
				slog.Any("error", err))
			return errors.New("调用龙译 API 创建作品集失败"), nil
		}

		// 更新本地工作集的龙译 ID
		if err := models.GetWorkset().Update(tx, &models.Workset{
			BaseModel: models.BaseModel{
				Id: workset.Id,
			},
			MoetranId: projSetInfo.ProjectSet.Id,
		}); err != nil {
			s.logger.Error("CreateWorkset 更新工作集龙译 ID 失败",
				slog.Any("workset_id", workset.Id),
				slog.Any("moetran_id", projSetInfo.ProjectSet.Id),
				slog.Any("error", err))
			return errors.New("更新工作集龙译 ID 失败"), nil
		}

		// 组装返回信息
		worksetCreatedInfo = &WorksetCreatedInfo{
			Message:   "工作集创建成功",
			WorksetId: uint(workset.Id),
			MoetranId: projSetInfo.ProjectSet.Id,
		}

		return nil, nil

	}); err != nil {
		s.logger.Error("CreateWorkset 事务执行失败",
			slog.Any("error", err))
		return nil, errors.New("创建工作集失败")
	}

	return worksetCreatedInfo, nil
}
