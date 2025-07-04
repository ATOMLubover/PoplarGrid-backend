package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"poplargrid/internal/api_server/config"
	"poplargrid/internal/shared/configutil"
	"poplargrid/internal/shared/logutils"
	"poplargrid/internal/update_server/crawler"
	"poplargrid/internal/update_server/persistence"
	"poplargrid/internal/update_server/transformer"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Update Server 负责从尨译拉取最新信息并更新到 PoplarGrid 的数据库中
func main() {
	// 先加载配置文件
	cfg := LoadConfig("update_config.yaml", "yaml")

	// 创建日志记录器
	logger := NewLogger(cfg)

	// 创建数据库上下文
	dbCtx := NewDatabase(cfg)

	// 创建各个 repo
	worksRepo := persistence.NewWorksRepo(dbCtx)
	worksetsRepo := persistence.NewWorksetsRepo(dbCtx, logger)
	teamsRepo := persistence.NewTeamsRepo(dbCtx)

	// 创建 transformer
	transformer := transformer.NewTransformer()

	// 创建 API 客户端
	apiClient := NewApiClient(cfg)

	// 创建 Crawler 实例
	crawler := crawler.NewCrawler(
		worksRepo, worksetsRepo, teamsRepo,
		apiClient,
		transformer,
		logger,
	)

	// 建立停止信号通道
	stopChan := make(chan os.Signal, 1)
	// 监听系统信号以便优雅地关闭
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// 启动爬虫
	crawler.Start()

	sig := <-stopChan
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
func LoadConfig(relPath string, cfgType string) *config.Config {
	// 先检查配置文件路径
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		panic(fmt.Sprintf("无法获取配置文件的绝对路径: %v", err))
	}

	var cfg config.Config
	if err := configutil.LoadConfig(&cfg, absPath, cfgType); err != nil {
		panic(fmt.Sprintf("无法加载配置: %v", err))
	}

	fmt.Printf("已加载配置文件：%v\n", cfg)

	// 在 debug 模式下打印内存中的结构体检查
	if cfg.Server.Mode == "debug" {
		fmt.Printf("已加载配置结构体：%+v\n", cfg)
	}

	return &cfg
}

// NewLogger 创建一个新的日志记录器
func NewLogger(cfg *config.Config) *slog.Logger {
	// 暂时不根据 cfg 配置使用不同的日志记录器
	lgr := logutils.NewLogger(nil)

	return lgr
}

// NewDatabase 创建一个新的数据库上下文
func NewDatabase(cfg *config.Config) *gorm.DB {
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
func NewApiClient(cfg *config.Config) *crawler.ApiClient {
	// 使用默认的尨译 API 基础 URL
	c := crawler.NewApiClient(cfg.Api.BaseUrl)

	// 添加上 Authorization Token
	c.ModifyAuthToken(cfg.Api.AuthToken)

	return c
}
