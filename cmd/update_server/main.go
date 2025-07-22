package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"poplargrid/internal/shared/configutil"
	"poplargrid/internal/shared/logutils"
	"poplargrid/internal/update_server/config"
	"poplargrid/internal/update_server/crawler"
	"poplargrid/internal/update_server/repository"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Update Server 负责从尨译拉取最新信息并更新到 PoplarGrid 的数据库中
func main() {
	// 先加载配置文件
	LoadConfig("update_config.yaml", "yaml")

	// 创建日志记录器
	InitLogger()

	// 创建数据库上下文
	dbCtx := NewDatabase()

	// 创建各个 repo
	projectsRepo := repository.NewProjectsRepo(dbCtx)
	worksetsRepo := repository.NewWorksetsRepo(dbCtx)
	teamsRepo := repository.NewTeamsRepo(dbCtx)
	membersRepo := repository.NewMembersRepo(dbCtx)
	usersRepo := repository.NewUsersRepo(dbCtx)

	// 创建 API 客户端
	apiClient := NewApiClient()

	// 创建 Crawler 实例
	crawler := crawler.NewCrawler(
		projectsRepo, worksetsRepo, teamsRepo, membersRepo, usersRepo,
		apiClient)

	// 建立停止信号通道
	shutdownChan := make(chan os.Signal, 1)
	// 监听系统信号以便优雅地关闭
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// 启动爬虫
	crawler.Start()

	sig := <-shutdownChan
	fmt.Printf("接收到信号: %v，正在关闭 Update Server...\n", sig)

	// 停止爬虫
	crawler.Stop()
	// 停止数据库连接
	db, err := dbCtx.DB()
	if err != nil {
		fmt.Printf("无法获取数据库连接: %v\n", err)
		fmt.Println("Update Server 非正常停止")
		return
	}
	db.Close()

	fmt.Println("Update Server 已停止")
}

// LoadConfig 加载网关服务器的配置
func LoadConfig(relPath string, cfgType string) {
	if err := config.LoadConfig(relPath, cfgType); err != nil {
		panic(fmt.Sprintf("无法加载配置: %v", err))
	}
}

// InitLogger 设置全局日志记录器
func InitLogger() {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("无法获取配置，日志记录器无法初始化")
	}

	// 暂时不根据 cfg 配置使用不同的日志记录器
	lgr := logutils.NewLogger(nil)

	// 因为项目体量小，最终决定还是直接使用全局 slog.Logger
	// 这样可以避免在每个模块中都传递日志记录器
	slog.SetDefault(lgr)
}

// NewDatabase 创建一个新的数据库上下文
func NewDatabase() *gorm.DB {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("无法获取配置，日志记录器无法初始化")
	}

	// 根据 cfg 来创建对应数据库连接
	switch cfg.Database.Type {
	case "postgresql":
		{
			// 使用 PostgreSQL 数据库
			dsn := fmt.Sprintf(
				`host=%s port=%d user=%s password=%s
                    dbname=%s sslmode=%s connect_timeout=%d`,
				cfg.Database.Host, cfg.Database.Port,
				cfg.Database.User, configutil.LoadEnvVariable(cfg.Database.PwdEnvVar, "TPOW2483137020#"),
				cfg.Database.DbName, cfg.Database.SslEnabled, cfg.Database.ConnectTimeout)

			// 连接到 PostgreSQL 数据库
			context, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				panic(fmt.Errorf("无法连接到 PostgreSQL 数据库: %v", err))
			}

			return context
		}

	default:
		// 目前只支持 PostgreSQL
		panic(fmt.Sprintf("不支持的数据库类型: %s", cfg.Database.Type))
	}
}

// NewApiClient 创建一个新的 API 客户端
func NewApiClient() *crawler.ApiClient {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("无法获取配置，日志记录器无法初始化")
	}

	// 使用默认的尨译 API 基础 URL
	c := crawler.NewApiClient(cfg.Api.BaseUrl)

	// 添加上 Authorization Token
	c.ModifyAuthToken(cfg.Api.AuthToken)

	return c
}
