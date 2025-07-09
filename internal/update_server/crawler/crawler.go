package crawler

import (
	"log/slog"
	"poplargrid/internal/shared/dbmodel"
	"poplargrid/internal/update_server/repository"
	"poplargrid/internal/update_server/transformer"
	"time"
)

// Crawler 是更新服务器的爬虫结构体
type Crawler struct {
	// 作品 repo
	projsRepo *repository.ProjectsRepo
	// 作品集 repo
	worksetsRepo *repository.WorksetsRepo
	// 汉化组 repo
	teamsRepo *repository.TeamsRepo

	// 用于与尨译 API 交互的 HTTP 客户端
	apiClient *ApiClient

	// 关闭 channel
	stopChan chan struct{}
}

// NewCrawler 创建一个新的 Crawler 实例
func NewCrawler(
	projsRepo *repository.ProjectsRepo,
	worksetsRepo *repository.WorksetsRepo,
	teamsRepo *repository.TeamsRepo,
	apiClient *ApiClient,
) *Crawler {
	return &Crawler{
		projsRepo:    projsRepo,
		worksetsRepo: worksetsRepo,
		teamsRepo:    teamsRepo,

		apiClient: apiClient,

		stopChan: make(chan struct{}),
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
	// 			slog.Info("爬虫停止")
	// 			return

	// 		case <-ticker.C:
	// 			slog.Info("开始从尨译获取汉化组信息")
	// 			if err := c.sFetchAndUpsert(); err != nil {
	// 				slog.Error("从尨译获取并更新信息失败", "error", err)
	// 				continue
	// 			}
	// 			// 暂时不处理错误和重试，直接打印日志
	// 			slog.Info("从尨译获取并更新信息成功")
	// 		}
	// 	}
	// }()

	if err := c.SyncMoetran(); err != nil {
		slog.Error("从尨译获取并更新信息失败", "error", err)
	}
}

// Stop 停止爬虫
func (c *Crawler) Stop() {
	// 发送停止信号到 stopChan
	close(c.stopChan)
	slog.Info("爬虫停止信号已发送")
}

// sSyncMoetran 从尨译获取并更新整个数据库的信息
// 目前仅获取并跟新作品集和作品信息
func (c *Crawler) SyncMoetran() error {
	// 从数据库获取所有汉化组信息（这个不依赖尨译直接拉取）
	teams, err := c.teamsRepo.Select()
	if err != nil {
		slog.Error("读取数据库汉化组信息失败", "error", err)
		return err
	}

	// 开始遍历更新各个汉化组的作品集和作品信息
	for _, team := range teams {
		if err := c.SyncWorksetsOfTeam(team); err != nil {
			slog.Error("同步汉化组作品集和作品失败",
				"team_name", team.Name, "team_id", team.Id, "error", err)
			continue // 继续处理下一个汉化组，避免直接退出
		}
	}

	return nil
}

// SyncWorksetsOfTeam 从尨译获取指定汉化组的所有作品集信息
// 这个函数拆分了整体逻辑，方便重试和调试
func (c *Crawler) SyncWorksetsOfTeam(team *dbmodel.Team) error {
	slog.Info("开始从尨译获取作品集信息",
		"team_name", team.Name, "team_id", team.Id)

	// 从尨译获取作品集信息，由于作品集不大，因此选择直接全部读入
	projSets, err := c.apiClient.
		GetAllProjSets(team.MoetranId)
	if err != nil {
		slog.Error("获取尨译作品集信息失败", "error", err)
		return err
	}

	slog.Info("获取尨译作品集信息成功", "length", len(projSets))

	// 将获取到的作品集信息转换成数据库模型
	worksets, err := transformer.ProjSetsToWorksets(team, projSets)
	if err != nil {
		slog.Error("转换作品集信息失败", "error", err)
		return err
	}

	slog.Info("转换作品集信息成功", "length", len(worksets))

	// 将转换后的作品集信息存入数据库
	if err := c.worksetsRepo.BulkUpsert(worksets); err != nil {
		slog.Error("批量插入作品集信息到数据库失败", "error", err)
		return err
	}

	slog.Info("批量插入作品集信息到数据库成功", "length", len(worksets))

	slog.Info("开始处理各个作品集中的作品")

	// 然后开始轮次处理各个作品集中的作品
	for _, workset := range worksets {
		if err := c.SyncWorksOfWorkset(team, workset); err != nil {
			slog.Error("处理作品集失败",
				"workset_name", workset.Name, "workset_id", workset.Id, "error", err)
			continue // 继续处理下一个作品集，避免直接退出
		}
	}

	slog.Info("所有作品集和作品处理完成")

	return nil
}

// SyncWorksOfWorkset 从尨译获取指定作品集的所有作品信息
// 这个函数拆分了整体逻辑，方便重试和调试
func (c *Crawler) SyncWorksOfWorkset(team *dbmodel.Team, workset *dbmodel.Workset) error {
	slog.Info("开始处理作品集",
		"workset_name", workset.Name, "workset_id", workset.Id)

	// 每个作品集下可能有多页作品，因此需要分页获取
	for page := 1; ; page++ {
		// 从尨译获取当前页的作品信息
		moeProjs, err := c.apiClient.
			GetPartProjs(team.MoetranId, workset.MoetranId, page)
		if err != nil {
			slog.Error("从尨译获取作品信息失败", "error", err, "page", page)
			return err
		}

		if len(moeProjs) == 0 {
			// 如果当前页没有作品信息，则说明已经是最后一页了
			slog.Info("当前页没有作品信息，可能是最后一页", "page", page)
			break
		}

		slog.Info("从尨译获取作品信息成功", "length", len(moeProjs), "page", page)

		// 将获取到的作品信息转换成数据库模型
		projects, err := transformer.
			ProjsToProjects(moeProjs, workset)
		if err != nil {
			slog.Error("转换作品信息失败", "error", err, "page", page)
			return err
		}

		// 将转换后的作品信息存入数据库
		if err := c.projsRepo.BulkUpsert(projects); err != nil {
			slog.Error("批量插入作品信息到数据库失败", "error", err, "page", page)
			return err
		}

		slog.Info("批量插入作品信息到数据库成功", "length", len(projects), "page", page)

		// 如果当前页的作品数量小于预设的每页作品数量，则确定已经是最后一页了
		// 这是对 len == 0 的补充判断
		if len(moeProjs) < PROJ_PAGE_SIZE {
			break // 退出当前作品集的处理循环
		}

		// 为了防止请求过快导致被限速，增加延时
		time.Sleep(50 * time.Millisecond)
	}

	slog.Info("当前作品集已处理完毕", "workset_name", workset.Name, "workset_id", workset.Id)

	return nil
}
