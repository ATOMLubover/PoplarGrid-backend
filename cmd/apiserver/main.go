package main

import (
	"fmt"
	"poplargrid/internal/apiserver/config"
	"strconv"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/kataras/iris/v12/mvc"

	_ "poplargrid/docs/api" // 引入生成的 Swagger 文档
)

// @title PoplarGrid 核心 API 服务器
// @version 1.0
// @description PoplarGrid API 服务器，提供项目、成员等操作功能
// @host localhost:8081
func main() {
	// 先加载全局配置
	LoadConfig("api_config.yaml", "yaml")

	// 初始化 Iris 应用实例
	app := InitIrisApp()

	// 初始化 Swagger
	ApplySwagger(app)
	// 初始化 MVC
	ApplyMvc(app)

	// 开始启动监听
	app.Listen(":" + strconv.Itoa(config.GetConfig().Server.Port))
}

// ========= 在 main 函数初始化加载流程中出现错误直接 panic ==========

// LoadConfig 加载网关服务器的配置
func LoadConfig(relPath string, cfgType string) *config.Config {
	// 先加载配置文件
	if err := config.LoadConfig(relPath, cfgType); err != nil {
		panic(fmt.Errorf("配置文件加载失败: %w", err))
	}

	return config.GetConfig()
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

func ApplySwagger(irisApp *iris.Application) {
	// 获取全局配置
	cfg := config.GetConfig()
	if cfg == nil {
		panic("配置未加载，应先调用 LoadConfig 函数")
	}

	if cfg.Server.Mode != "debug" {
		// 在非 debug 模式下不启用 Swagger
		return
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
}

// ApplyMvc 构造一个封装底层 Iris 应用的 MVC 应用
func ApplyMvc(irisApp *iris.Application) *mvc.Application {
	// 直接包装整个根路由
	mvcApp := mvc.New(irisApp)

	return mvcApp
}
