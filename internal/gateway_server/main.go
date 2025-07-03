package main

import (
	"fmt"
	"path/filepath"
	"poplargrid/internal/gateway_server/config"
	"poplargrid/internal/gateway_server/routes"
	"poplargrid/internal/shared/configutil"
	"strconv"

	"github.com/iris-contrib/swagger/swaggerFiles"
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/kataras/iris/v12/mvc"

	_ "poplargrid/docs/gateway" // 引入生成的 Swagger 文档
)

// @title PoplarGrid 网关服务器
// @version 1.0
// @description PoplarGrid 服务器的网关服务器，负责限流、转发和协调请求
// @host localhost:8080
func main() {
	cfg := LoadConfig("gateway_config.yaml", "yaml")

	app := NewIrisApp(cfg)

	// 初始化 Swagger
	InitSwagger(app, cfg)

	// 初始化 MVC 应用
	mvcApp := NewMvcApp(app, cfg)
	InitMvcApp(mvcApp, cfg)

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
	// 添加 /transfer 子路由组
	routes.ConfigureTransferRoutes(root)
}
