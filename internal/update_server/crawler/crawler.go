package crawler

import (
	"log/slog"
	"poplargrid/internal/update_server/persistence"
	"poplargrid/internal/update_server/transformer"
)

// Crawler 是更新服务器的爬虫结构体
type Crawler struct {
	// 作品 repo
	worksRepo *persistence.WorksRepo
	// 作品集 repo
	worksetsRepo *persistence.WorksetsRepo
	// 汉化组 repo
	teamsRepo *persistence.TeamsRepo

	// 用于与尨译 API 交互的 HTTP 客户端
	apiClient *ApiClient
	// transformer 用于将从尨译获取的数据转换成数据库模型
	transformer *transformer.Transformer

	// 关闭 channel
	stopChan chan struct{}

	// 日志器
	logger *slog.Logger
}

// NewCrawler 创建一个新的 Crawler 实例
func NewCrawler(
	worksRepo *persistence.WorksRepo,
	worksetsRepo *persistence.WorksetsRepo,
	teamsRepo *persistence.TeamsRepo,
	apiClient *ApiClient,
	transformer *transformer.Transformer,
	logger *slog.Logger,
) *Crawler {
	return &Crawler{
		worksRepo:    worksRepo,
		worksetsRepo: worksetsRepo,
		teamsRepo:    teamsRepo,

		apiClient:   apiClient,
		transformer: transformer,

		stopChan: make(chan struct{}),

		logger: logger,
	}
}

// Start 启动爬虫，这会启动一个后台 goroutine 用来 tick
func (c *Crawler) Start() {
	// // 获取一个每 1 个小时刷新一次的 ticker
	// go func() {
	// 	ticker := time.NewTicker(time.Hour)
	// 	defer ticker.Stop()

	// 	for {
	// 		select {
	// 		case <-c.stopChan:
	// 			c.logger.Info("爬虫停止")
	// 			return

	// 		case <-ticker.C:
	// 			c.logger.Info("开始从尨译获取汉化组信息")
	// 			if err := c.sFetchAndUpsert(); err != nil {
	// 				c.logger.Error("从尨译获取并更新信息失败", "error", err)
	// 				continue
	// 			}
	// 			// 暂时不处理错误和重试，直接打印日志
	// 			c.logger.Info("从尨译获取并更新信息成功")
	// 		}
	// 	}
	// }()

	if err := c.sFetchAndUpsert(); err != nil {
		c.logger.Error("从尨译获取并更新信息失败", "error", err)
	}
}

// Stop 停止爬虫
func (c *Crawler) Stop() {
	// 发送停止信号到 stopChan
	close(c.stopChan)
	c.logger.Info("爬虫停止信号已发送")
}

// sFetchAndUpdate 从尨译获取并更新整个数据库的信息
func (c *Crawler) sFetchAndUpsert() error {
	// 从尨译获取作品集信息，由于作品集不大，因此选择直接全部读入
	projectSets, err := c.sFetchWorksets("64aac7d91b30e3645d7f96c6")
	if err != nil {
		c.logger.Error("fetch尨译作品集信息失败", "error", err)
		return err
	}
	c.logger.Info("fetch尨译作品集信息成功", "length", len(projectSets))

	// 将获取到的作品集信息转换成数据库模型
	worksets, err := c.transformer.ProjectSetsToWorksets(projectSets)
	if err != nil {
		c.logger.Error("转换作品集信息失败", "error", err)
		return err
	}

	c.logger.Info("转换作品集信息成功", "length", len(worksets))

	// 将转换后的作品集信息存入数据库
	if err := c.worksetsRepo.BulkUpsert(worksets); err != nil {
		c.logger.Error("批量插入作品集信息到数据库失败", "error", err)
		return err
	}

	return nil
}

// // sFetchTeams 从尨译获取汉化组信息
// func (c *Crawler) sFetchTeams() ([]MoetranTeam, error) {
//
// }

// sFetchWorksets 从尨译获取指定汉化组的作品集信息
func (c *Crawler) sFetchWorksets(teamMoetranId string) ([]transformer.MoetranProjSet, error) {
	// 调用客户端获取 teamMoetranId 的所有作品集
	allProjSets, err := c.apiClient.GetProjectSetUri(teamMoetranId)
	if err != nil {
		c.logger.Error("从尨译获取作品集信息失败", "error", err)
		return nil, err
	}

	c.logger.Info("从尨译获取作品集信息成功", "length", len(allProjSets))

	return allProjSets, nil
}

// // sFetchWorks 从尨译获取作品信息
// func (c *Crawler) sFetchWorks() ([]MoetranProj, error) {
// 	// 先获取所有汉化组的信息
// 	teams, err := c.teamsRepo.SelectNameAndId()
// 	if err != nil {
// 		c.logger.Error("从数据库获取所有汉化组列表信息失败", "error", err)
// 		return nil, err
// 	}

// 	// 根据汉化组列表，逐个获取其下所有作品信息

// 	return nil, nil
// }
