package crawler

import (
	"log/slog"
	"poplargrid/internal/api_server/apiclient"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Crawler 定义了尨译爬虫
type Crawler interface {
	// UpdateTeams 更新汉化组信息
	UpdateTeams(moetranAuth string) error
	// UpdateProjectSets 更新项目集信息
	UpdateProjectSets(teamMoetranID, moetranAuth string) error
	// UpdateProjects 更新指定作品集下的项目信息
	UpdateProjects(teamMoetranID, projectSetMoetranID, moetranAuth string) error
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

// UpdateTeams 实现了 Crawler 接口的 UpdateTeams 方法
func (c *crawlerImpl) UpdateTeams(moetranAuth string) error {
	// 调用尨译 API 获取汉化组信息
	teams, err := c.apiClient.GetUserTeams(moetranAuth)
	if err != nil {
		c.logger.Error("UpdateTeams 调用龙译 API 获取汉化组信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的汉化组信息转化为本地模型
	poplarTeams := teamMoetranToPoplar(teams.Teams)

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

	return nil
}

// UpdateProjectSets 实现了 Crawler 接口的 UpdateProjectSets 方法
func (c *crawlerImpl) UpdateProjectSets(teamMoetranID, moetranAuth string) error {
	// 调用尨译 API 获取项目集信息
	projectSets, err := c.apiClient.GetTeamProjectSets(teamMoetranID, moetranAuth)
	if err != nil {
		c.logger.Error("UpdateProjectSets 调用龙译 API 获取项目集信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的项目集信息转化为本地模型
	poplarProjectSets := setMoetranToPoplar(projectSets.Sets)

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

	return nil
}

// UpdateProjects 实现了 Crawler 接口的 UpdateProjects 方法
func (c *crawlerImpl) UpdateProjects(teamMoetranID, projectSetMoetranID, moetranAuth string) error {
	// 调用尨译 API 获取项目集下的项目信息
	projects, err := c.apiClient.GetProjectInfo(teamMoetranID, projectSetMoetranID, moetranAuth)
	if err != nil {
		c.logger.Error("UpdateProjects 调用龙译 API 获取项目信息失败", slog.Any("error", err))
		return err
	}

	// 将获取到的项目信息转化为本地模型
	poplarProjects := projectMoetranToPoplar(projects.Projects)

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

	return nil
}
