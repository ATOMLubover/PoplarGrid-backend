package crawler

import (
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Crawler 定义了尨译爬虫
type Crawler interface {
	// RecurseUpdate 自动递归地更新当前用户所有汉化组的项目信息
	RecurseUpdate(userID uint, moetranAuth string) error
}

// crawlerImpl 实现了 Crawler 接口
type crawlerImpl struct {
	handle    *gorm.DB
	apiClient apiclient.ApiClient
	logger    *slog.Logger
}

// NewCrawler 创建一个新的 Crawler 实例
func NewCrawler(
	handle *gorm.DB,
	apiClient apiclient.ApiClient,
	logger *slog.Logger,
) Crawler {
	return &crawlerImpl{
		handle:    handle,
		apiClient: apiClient,
		logger:    logger,
	}
}

// RecurseUpdate 自动递归地更新当前用户所有汉化组的项目信息
func (c *crawlerImpl) RecurseUpdate(userID uint, moetranAuth string) error {
	return c.updateTeams(userID, moetranAuth)
}

// updateTeams 更新指定的汉化组信息，以及其中所有成员、项目集和项目信息
func (c *crawlerImpl) updateTeams(userID uint, moetranAuth string) error {
	// 调用尨译 API 获取汉化组信息
	moetranTeams, err := c.apiClient.GetUserTeams(moetranAuth)
	if err != nil {
		c.logger.Error("UpdateTeams 调用龙译 API 获取汉化组信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的汉化组信息转化为本地模型
	poplarTeams := teamMoetranToPoplar(moetranTeams.Teams)

	// 将转化后的汉化组信息保存到数据库
	// 此处根据 moetran_id 来进行 ON CONFLICT 更新
	if err := c.handle.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "moetran_id"}},
			// 使用 DoUpdates 保证返回 ID
			DoUpdates: clause.AssignmentColumns([]string{"name"}),
		}).
		Create(&poplarTeams).Error; err != nil {
		c.logger.Error("UpdateTeams 保存汉化组信息到数据库失败",
			slog.Any("error", err))
		return err
	}

	// 同时为该成员更新 member 信息（所有因为他而被 upsert 的汉化组）
	members := make([]*models.Member, len(poplarTeams))

	for i, team := range poplarTeams {
		members[i] = &models.Member{
			UserId: models.PKey(userID),
			TeamId: team.Id,
		}
	}

	if err := c.handle.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "team_id"},
			},
			// 使用 DoUpdates 保证返回 ID
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		}).
		Create(&members).Error; err != nil {
		c.logger.Error("UpdateTeams 保存成员信息到数据库失败",
			slog.Any("error", err))
		return err
	}

	// 遍历所有汉化组，更新每个汉化组的项目集信息
	for _, team := range poplarTeams {
		var memberID uint
		for _, m := range members {
			if m.TeamId == team.Id {
				memberID = uint(m.Id)
				break
			}
		}

		if err := c.updateProjectSets(memberID, &team, moetranAuth); err != nil {
			c.logger.Error("UpdateTeams 更新项目集信息失败",
				slog.Any("error", err))
			return err
		}
	}

	return nil
}

// updateProjectSets 更新指定作品集的信息，以及其中所有的项目信息
func (c *crawlerImpl) updateProjectSets(memberID uint, team *models.Team, moetranAuth string) error {
	// 调用尨译 API 获取项目集信息
	moetranProjectSets, err := c.apiClient.GetTeamProjectSets(team.MoetranId, moetranAuth)
	if err != nil {
		c.logger.Error("UpdateProjectSets 调用龙译 API 获取项目集信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的项目集信息转化为本地模型
	poplarProjectSets := setMoetranToPoplar(team, moetranProjectSets.Sets)

	// 将转化后的项目集信息保存到数据库
	if err := c.handle.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "moetran_id"}},
			// 使用 DoUpdates 保证返回 ID
			DoUpdates: clause.AssignmentColumns([]string{"name"}),
		}).
		Create(&poplarProjectSets).Error; err != nil {
		c.logger.Error("UpdateProjectSets 保存项目集信息到数据库失败", slog.Any("error", err))
		return err
	}

	// 遍历每个作品集，更新其中的项目信息
	for _, workset := range poplarProjectSets {
		if err := c.updateProjects(memberID, team, &workset, moetranAuth); err != nil {
			c.logger.Error("UpdateProjectSets 更新项目集下的项目信息失败", slog.Any("error", err))
			return err
		}
	}

	return nil
}

// updateProjects 更新指定作品集的所有项目
func (c *crawlerImpl) updateProjects(memberID uint, team *models.Team, workset *models.Workset, moetranAuth string) error {
	const LIMIT = 50 // 每次请求获取的项目数量

	for page := 1; ; page++ {
		// 调用尨译 API 获取项目集下的项目信息
		moetranProjects, err := c.apiClient.GetProjects(
			team.MoetranId, workset.MoetranId,
			page, LIMIT, moetranAuth)
		if err != nil {
			c.logger.Error("UpdateProjects 调用龙译 API 获取项目信息失败", slog.Any("error", err))
			return err
		}

		if len(moetranProjects.Projects) == 0 {
			// 如果没有更多项目了，结束循环
			c.logger.Info("UpdateProjects 结束更新作品集",
				slog.String("team_moetran_id", team.MoetranId),
				slog.String("project_set_moetran_id", workset.MoetranId))
			break
		}

		// 将获取到的项目信息转化为本地模型
		poplarProjects := projectMoetranToPoplar(memberID, workset, moetranProjects.Projects)

		// 将转化后的项目信息保存到数据库
		if err := c.handle.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "moetran_id"}},
				DoNothing: true,
			}).
			Create(&poplarProjects).Error; err != nil {
			c.logger.Error("UpdateProjects 保存项目信息到数据库失败", slog.Any("error", err))
			return err
		}

		if len(moetranProjects.Projects) < LIMIT {
			// 如果本次获取的项目数量少于 LIMIT，说明已经没有更多项目了
			c.logger.Info("UpdateProjects 结束更新作品集",
				slog.String("team_moetran_id", team.MoetranId),
				slog.String("project_set_moetran_id", workset.MoetranId))
			break
		}
	}

	return nil
}
