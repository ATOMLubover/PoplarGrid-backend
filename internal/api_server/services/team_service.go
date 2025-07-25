package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
)

// TeamService 接口定义了团队服务的基本操作
type TeamService interface {
	// GetTeams 获取指定用户的团队列表
	GetTeams(params *TeamListParams) ([]*TeamInfo, error)
}

// TeamListParams 定义了团队列表的查询参数
type TeamListParams struct {
	Offset int  // 偏移量
	Limit  int  // 限制数量
	UserId uint // 用户 ID
}

// TeamInfo 定义了团队的基本信息
type TeamInfo struct {
	Id          uint   // 团队 ID
	Name        string // 团队名称
	Description string // 团队描述
	MoetranId   string // 龙译 ID
}

// // CreateProjectParams 定义了创建项目所需的信息
// type CreateProjectParams struct {
// 	Title       string // 项目标题
// 	Description string // 项目简介
// 	WorksetId   uint   // 作品集 ID

// 	CreatorUserId uint // 创建者的用户 ID
// 	AllowAutoJoin bool // 允许加入的权限
// 	IsHidden      bool // 是否隐藏项目
// }

// teamServiceImpl 是 TeamService 的实现
type teamServiceImpl struct {
	handle *gorm.DB
	logger *slog.Logger
}

// NewTeamService 创建一个新的 TeamService 实例
func NewTeamService(
	hdl *gorm.DB,
	lgr *slog.Logger,
) TeamService {
	return &teamServiceImpl{
		handle: hdl,
		logger: lgr,
	}
}

// GetTeams 实现 TeamService 接口的 GetTeams 方法
func (s *teamServiceImpl) GetTeams(params *TeamListParams) ([]*TeamInfo, error) {
	// 查询条件为用户 ID
	userPKey := models.PKey(params.UserId)
	memberSpec := &models.MemberSpec{
		UserId: &userPKey,
	}
	// 所需字段
	memberFields := &models.MemberFields{
		TeamId:     true,
		TeamFields: true, // 需要 Preload 团队信息
	}

	// 查询其对应的团队
	members, err := models.GetMember().SelectMany(
		s.handle, memberSpec, memberFields,
		&params.Offset, &params.Limit)
	if err != nil {
		s.logger.Error("GetTeams 查询成员信息失败", slog.Any("error", err))
		return nil, errors.New("查询用户的汉化组列表失败")
	}

	// 将查询结果转换为团队信息
	teamInfos := make([]*TeamInfo, 0, len(members))

	for _, member := range members {
		if member.FkTeam == nil {
			s.logger.Error("GetTeams 查询到无团队成员",
				slog.Any("member_id", member.BaseModel.Id),
				slog.Any("user_id", member.UserId),
				slog.Any("team_id", member.TeamId))
			return nil, errors.New("出现异常的无团队成员")
		}

		teamInfos = append(teamInfos, &TeamInfo{
			Id:          uint(member.FkTeam.Id),
			Name:        member.FkTeam.Name,
			Description: member.FkTeam.Description,
			MoetranId:   member.FkTeam.MoetranId,
		})
	}

	return teamInfos, nil
}
