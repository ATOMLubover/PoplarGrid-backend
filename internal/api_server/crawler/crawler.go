package crawler

import (
	"log/slog"
	"poplargrid/internal/api_server/apiclient"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Crawler 定义了尨译爬虫
type Crawler interface {
	// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
	AutoUpdateAll(moetranAuth string) error
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

// AutoUpdateAll 自动递归地更新当前用户所有汉化组的项目信息
func (c *crawlerImpl) AutoUpdateAll(moetranAuth string) error {
	return c.updateTeams(moetranAuth)
}

// updateTeams 更新指定的汉化组信息，以及其中所有的项目集和项目信息
func (c *crawlerImpl) updateTeams(moetranAuth string) error {
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
			Columns:   []clause.Column{{Name: "moetran_id"}},
			DoNothing: true,
		}).
		Create(&poplarTeams).Error; err != nil {
		c.logger.Error("UpdateTeams 保存汉化组信息到数据库失败", slog.Any("error", err))
		return err
	}

	// 遍历所有汉化组，更新每个汉化组的项目集信息
	for _, team := range poplarTeams {
		if err := c.updateProjectSets(team.MoetranId, moetranAuth); err != nil {
			c.logger.Error("UpdateTeams 更新项目集信息失败", slog.Any("error", err))
			// 选择直接跳过该汉化组
			continue
		}
	}

	return nil
}

// updateProjectSets 更新指定作品集的信息，以及其中所有的项目信息
func (c *crawlerImpl) updateProjectSets(teamMoetranID, moetranAuth string) error {
	// 调用尨译 API 获取项目集信息
	moetranProjectSets, err := c.apiClient.GetTeamProjectSets(teamMoetranID, moetranAuth)
	if err != nil {
		c.logger.Error("UpdateProjectSets 调用龙译 API 获取项目集信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的项目集信息转化为本地模型
	poplarProjectSets := setMoetranToPoplar(moetranProjectSets.Sets)

	// 将转化后的项目集信息保存到数据库
	if err := c.handle.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "moetran_id"}},
			DoNothing: true,
		}).
		Create(&poplarProjectSets).Error; err != nil {
		c.logger.Error("UpdateProjectSets 保存项目集信息到数据库失败", slog.Any("error", err))
		return err
	}

	// 遍历每个作品集，更新其中的项目信息
	for _, projectSet := range poplarProjectSets {
		if err := c.updateProjects(teamMoetranID, projectSet.MoetranId, moetranAuth); err != nil {
			c.logger.Error("UpdateProjectSets 更新项目集下的项目信息失败", slog.Any("error", err))
			// 选择直接跳过该作品集
			continue
		}
	}

	return nil
}

// updateProjects 更新指定作品集的所有项目
func (c *crawlerImpl) updateProjects(teamMoetranID, projectSetMoetranID, moetranAuth string) error {
	const LIMIT = 50 // 每次请求获取的项目数量

	for page := 1; ; page++ {
		// 调用尨译 API 获取项目集下的项目信息
		moetranProjects, err := c.apiClient.GetProjects(
			teamMoetranID, projectSetMoetranID,
			page, LIMIT, moetranAuth)
		if err != nil {
			c.logger.Error("UpdateProjects 调用龙译 API 获取项目信息失败", slog.Any("error", err))
			return err
		}

		if len(moetranProjects.Projects) == 0 {
			// 如果没有更多项目了，结束循环
			c.logger.Info("UpdateProjects 结束更新作品集",
				slog.String("team_moetran_id", teamMoetranID),
				slog.String("project_set_moetran_id", projectSetMoetranID))
			break
		}

		// 将获取到的项目信息转化为本地模型
		poplarProjects := projectMoetranToPoplar(moetranProjects.Projects)

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
				slog.String("team_moetran_id", teamMoetranID),
				slog.String("project_set_moetran_id", projectSetMoetranID))
			break
		}
	}

	return nil
}
