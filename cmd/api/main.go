package main

import (
	"fmt"
	"path/filepath"
	"poplargrid/internal/api_server/config"
	"poplargrid/internal/api_server/handler"
	"poplargrid/internal/api_server/repository"
	"poplargrid/internal/api_server/route"
	"poplargrid/internal/api_server/services"
	"poplargrid/internal/shared/configutil"
	"strconv"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "poplargrid/docs/api" // 引入生成的 Swagger 文档
)

// @title PoplarGrid API 服务器
// @version 1.0
// @description PoplarGrid 服务器的 API 处理服务器，负责协调表格相关业务
// @host localhost:8071
func main() {
	// 先加载配置文件
	cfg := LoadConfig("api_config.yaml", "yaml")

	// 构造底层 Iris 应用
	app := NewIrisApp(cfg)

	// 初始化 Swagger
	InitSwagger(app, cfg)

	// 初始化 MVC 应用
	InitMvcApp(NewMvcApp(app, cfg), cfg)

	// 开始启动监听
	app.Listen(":" + strconv.Itoa(cfg.Server.Port))
}

// ========= 在 main 函数初始化加载流程中出现错误直接 panic ==========

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

// NewIrisApp 构造底层 Iris 应用，并辅以简单的配置
func NewIrisApp(cfg *config.Config) *iris.Application {
	// 创建一个新的 Iris 应用实例
	app := iris.New()

	// 设置应用的配置
	switch cfg.Server.Mode {
	case "debug":
		{
			app.Logger().SetLevel("debug")

			// 启用 request ID 中间件
			app.UseRouter(requestid.New())
			// 启用 recovery 中间件
			app.UseRouter(recover.New())
		}
	case "release":
		{
			// 在 release 模式下设置日志级别为 error
			app.Logger().SetLevel("error")

			// 启用 recovery 中间件
			app.UseRouter(recover.New())
		}
	default:
		{
			panic("未知的运行模式，请检查配置文件")
		}
	}

	return app
}

// InitSwagger 初始化 Swagger 文档和 UI
func InitSwagger(app *iris.Application, cfg *config.Config) {
	if cfg.Server.Mode != "debug" {
		// 在非 debug 模式下不启用 Swagger
		return
	}

	// 注册 Swagger UI 和文档路由
	// 明确处理 /swagger 路径，重定向到 /swagger/index.html，
	// 此处使用 302 Found 来避免浏览器缓存问题
	app.Get("/swagger", func(ctx iris.Context) {
		ctx.Redirect("/swagger/index.html", iris.StatusFound)
	})

	// 注册 Swagger UI 的 handler
	app.Get("/swagger/{any:path}",
		swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
			// 指定获取 Swagger 文档的 URL
			c.URL = "/swagger/doc.json" // Modified to relative path
		}))
}

// NewMvcApp 构造一个封装底层 Iris 应用的 MVC 应用
func NewMvcApp(irisApp *iris.Application, cfg *config.Config) *mvc.Application {
	// 直接包装整个根路由
	mvcApp := mvc.New(irisApp)

	// 设置 MVC 应用的配置
	mvcApp.Register(cfg)

	return mvcApp
}

// InitMvcApp 专门负责 MVC 应用的复杂初始化
func InitMvcApp(root *mvc.Application, cfg *config.Config) {
	// 初始化所有 repositories
	dbCtx := NewDatabase(cfg)
	relTables := repository.NewRelationTables(dbCtx)
	memberRepo := repository.NewMembersRepo(dbCtx)
	projRepo := repository.NewProjectsRepo(dbCtx)

	// 初始化所有 services
	memberSrv := services.NewMemberService(memberRepo)
	projSrv := services.NewProjectService(projRepo)

	// 注册 repositories 和 services 到 MVC 应用
	root.Register(dbCtx, relTables, memberRepo, projRepo,
		memberSrv, projSrv)

	// 添加 /member 子路由组
	route.ConfigureMemberRoutes(root)
	// 添加 /project 子路由组
	{
		projHandler := root.Party("/project")

		// 注册 ProjectHandler
		projHandler.Handle(new(handler.ProjectHandler))
	}
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
