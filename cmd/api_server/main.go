package main

import (
	"fmt"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/api_server/config"
	"poplargrid/internal/api_server/handlers"
	"poplargrid/internal/api_server/repos"
	"poplargrid/internal/api_server/services"
	"poplargrid/internal/shared/configutil"
	"poplargrid/internal/shared/logutils"
	"strconv"
	"time"

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

// @title PoplarGrid 核心 API 服务器
// @version 1.0
// @description PoplarGrid API 服务器，提供项目、成员等操作功能
// @host localhost:8081
func main() {
	// 先加载全局配置
	LoadConfig("api_config.yaml", "yaml")

	// 初始化 Iris 应用实例，初始化 Swagger ，并初始化 MVC
	ApplyMvc(ApplySwagger(InitIrisApp())).
		// 开始启动监听
		Listen(":" + strconv.Itoa(config.GetConfig().Server.Port))
}

// ========= 在 main 函数初始化加载流程中出现错误直接 panic ==========

// LoadConfig 加载网关服务器的配置
func LoadConfig(relPath string, cfgType string) {
	// 先加载配置文件
	if err := config.LoadConfig(relPath, cfgType); err != nil {
		panic(fmt.Errorf("配置文件加载失败: %w", err))
	}
}

// InitLogger 设置全局日志记录器
func InitLogger(cfg *config.Config) {
	// 暂时不根据 cfg 配置使用不同的日志记录器
	lgr := logutils.NewLogger(nil)

	// 因为项目体量小，最终决定还是直接使用全局 slog.Logger
	// 这样可以避免在每个模块中都传递日志记录器
	slog.SetDefault(lgr)
}

// InitIrisApp 生成一个根据 cfg 调教过的 Iris 应用实例
func InitIrisApp() *iris.Application {
	// 获取全局配置
	cfg := config.GetConfig()
	if cfg == nil {
		panic("配置未加载，应先调用 LoadConfig 函数")
	}

	// 创建一个新的 Iris 应用实例
	app := iris.New()

	// 设置应用的配置
	switch cfg.Server.Mode {
	case "readonly":
		{
			// 用于前端检查 API 文档的模式
			app.Logger().SetLevel("info")
		}
	case "debug":
		{
			app.Logger().SetLevel("debug")
			app.UseRouter(requestid.New())
			app.UseRouter(recover.New())
		}
	case "release":
		{
			app.Logger().SetLevel("error")
			app.UseRouter(recover.New())
		}
	default:
		panic("未知的运行模式，请检查配置文件")
	}

	return app
}

func ApplySwagger(irisApp *iris.Application) *iris.Application {
	// 获取全局配置
	cfg := config.GetConfig()
	if cfg == nil {
		panic("配置未加载，应先调用 LoadConfig 函数")
	}

	if cfg.Server.Mode == "release" {
		// 在 release 模式下不启用 Swagger
		return irisApp
	}

	// 注册 Swagger UI 和文档路由
	// 明确处理 /swagger 路径，重定向到 /swagger/index.html，
	// 此处使用 302 Found 来避免浏览器缓存问题
	irisApp.Get("/swagger", func(ctx iris.Context) {
		ctx.Redirect("/swagger/index.html", iris.StatusFound)
	})

	// 注册 Swagger UI 的 handler
	irisApp.Get("/swagger/{any:path}",
		swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
			// 指定获取 Swagger 文档的 URL
			c.URL = "/swagger/doc.json" // Modified to relative path
		}))

	return irisApp
}

// InitializeDatabse 初始化数据库连接
func InitDatabase() *gorm.DB {
	cfg := config.GetConfig()
	if cfg == nil {
		panic("配置未加载，应先调用 LoadConfig 函数")
	}

	if cfg.Server.Mode == "readonly" {
		// 只读模式下不连接数据库
		return nil
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

// ApplyMvc 构造一个封装底层 Iris 应用的 MVC 应用
func ApplyMvc(irisApp *iris.Application) *iris.Application {
	// 获取全局配置
	cfg := config.GetConfig()
	if cfg == nil {
		panic("配置未加载，应先调用 LoadConfig 函数")
	}

	if cfg.Server.Mode == "readonly" {
		// 只读模式下不启用 MVC
		return irisApp
	}

	// 直接包装整个根路由
	mvcApp := mvc.New(irisApp)

	// 获取各个 repo 的实例
	handle := InitDatabase()

	userRepo := repos.NewUserRepo(handle)
	projRepo := repos.NewProjectRepo(handle)
	// teamRepo := repos.NewTeamRepo(handle)
	worksetRepo := repos.NewWorksetRepo(handle)
	appliRepo := repos.NewAppliRepo(handle)
	invitationRepo := repos.NewInvitationRepo(handle)
	laborRepo := repos.NewLaborRepo(handle)
	memberRepo := repos.NewTeamMemberRepo(handle)
	materialView := repos.NewMaterialView(handle)

	// 注册龙译 API Client
	apiClient := apiclient.NewApiClient(cfg.Api.BaseUrl, *slog.Default())

	// 注册各个 service 的依赖
	mvcApp.Register(
		services.NewLaborService(invitationRepo, appliRepo, slog.Default()),
		services.NewUserService(userRepo),
		services.NewProjectService(projRepo, laborRepo, userRepo, apiClient, slog.Default()),
		services.NewTeamService(memberRepo, slog.Default()),
		services.NewWorksetService(worksetRepo, materialView, slog.Default()),
	)

	// 注册中间件
	irisApp.Use(
		// 跨域控制
		handlers.NewCorsMiddleware(
			cfg.Server.CorsOrigins,
			cfg.Server.CorsMethods,
			cfg.Server.CorsHeaders,
			cfg.Server.CorsWithCredentials,
			time.Duration(cfg.Server.CorsMaxAge)*time.Second,
		),
		// cookie 和请求头预处理
		handlers.NewUserInfoExtractMiddleware(),
	)

	// 注册路由处理器
	handlers.RouteUserHandler(mvcApp)
	handlers.RouteProjectHandler(mvcApp)
	handlers.RouteWorksetHandler(mvcApp)
	handlers.RouteTeamHandler(mvcApp)
	handlers.RouteLaborProcHandler(mvcApp)

	return irisApp
}
