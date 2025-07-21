package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
	"poplargrid/internal/shared/dbmodels"
)

// TeamService 接口定义了团队服务的基本操作
type TeamService interface {
	// GetBasicPageByUserId 获取指定用户 ID 的团队列表
	GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.TeamBasic, error)
	// GetMemberBasicPage 获取指定团队 ID 的成员列表，支持分页
	GetMemberBasicPage(teamId uint, pageSerial, pageSize int) ([]*dtos.MemberBasic, error)
}

// teamServiceImpl 是 TeamService 的实现
type teamServiceImpl struct {
	teamMemberRepo repos.TeamMemberRepo
	logger         *slog.Logger
}

// NewTeamService 创建一个新的 TeamService 实例
func NewTeamService(
	teamMemberRepo repos.TeamMemberRepo,
) TeamService {
	return &teamServiceImpl{
		teamMemberRepo: teamMemberRepo,
	}
}

// GetBasicPageByUserId 实现 TeamService 接口的 GetBasicPageByUserId 方法
func (s *teamServiceImpl) GetBasicPageByUserId(userId uint, pageSerial, pageSize int) ([]*dtos.TeamBasic, error) {
	// 调用仓库方法获取团队列表
	teams, err := s.teamMemberRepo.SelectTeamBasicByUserId(dbmodels.PrimaryKey(userId))
	if err != nil {
		s.logger.Error("GetBasicPage 调用 SelectBasicPage 中出现错误", slog.Any("error", err))
		return nil, errors.New("无法获取指定用户的团队列表")
	}

	// 将 dbmodels.Team 转换为 dtos.TeamBasic
	var teamBasics []*dtos.TeamBasic
	for _, team := range teams {
		teamBasics = append(teamBasics, &dtos.TeamBasic{
			Id:   uint(team.Id),
			Name: team.Name,
		})
	}

	return teamBasics, nil
}

// GetMemberBasicPage 实现 TeamService 接口的 GetMemberBasicPage 方法
func (s *teamServiceImpl) GetMemberBasicPage(teamId uint, pageSerial, pageSize int) ([]*dtos.MemberBasic, error) {
	// 调用仓库方法获取成员列表
	members, err := s.teamMemberRepo.SelectUserBasicPage(dbmodels.PrimaryKey(teamId), (pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetMemberListPage 调用 SelectMemberBasicPage 中出现错误", slog.Any("error", err))
		return nil, err
	}

	// 将 dbmodels.TeamMember 转换为 dtos.MemberBasic
	var memberBasics []*dtos.MemberBasic

	for _, member := range members {
		memberBasics = append(memberBasics, &dtos.MemberBasic{
			Id:       uint(member.Id),
			UserId:   uint(member.UserId),
			Role:     uint(member.Role),
			Nickname: member.FkUser.Nickname,
		})
	}

	return memberBasics, nil
}
