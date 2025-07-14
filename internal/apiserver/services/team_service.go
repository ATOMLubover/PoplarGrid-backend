package services

import (
	"log/slog"
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/repos"
)

// TeamService 接口定义了团队服务的基本操作
type TeamService interface {
	// GetBasicPage 获取团队列表，支持分页
	GetBasicPage(pageSerial, pageSize int) ([]*dtos.TeamBasic, error)
}

// teamServiceImpl 是 TeamService 的实现
type teamServiceImpl struct {
	teamRepo repos.TeamRepo
	logger   *slog.Logger
}

// NewTeamService 创建一个新的 TeamService 实例
func NewTeamService(
	teamRepo repos.TeamRepo,
) TeamService {
	return &teamServiceImpl{
		teamRepo: teamRepo,
	}
}

// GetBasicPage 实现 TeamService 接口的 GetBasicPage 方法
func (s *teamServiceImpl) GetBasicPage(pageSerial, pageSize int) ([]*dtos.TeamBasic, error) {
	// 调用仓库方法获取团队列表
	teams, err := s.teamRepo.SelectBasicPage((pageSerial-1)*pageSize, pageSize)
	if err != nil {
		s.logger.Error("GetBasicPage 调用 SelectBasicPage 中出现错误", slog.Any("error", err))
		return nil, err
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
