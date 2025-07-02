package main

import (
	"fmt"
	"path/filepath"
	"poplargrid/internal/gateway_server/config"
	"poplargrid/internal/gateway_server/routes"
	"poplargrid/internal/shared/configutil"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	_ "poplargrid/docs/gateway" // 引入生成的 Swagger 文档
)

// @title PoplarGrid 网关 API
// @version 1.0
// @description PoplarGrid 服务器的网关服务器，负责限流、转发和协调请求
// @host localhost:8080
func main() {
	cfg := LoadConfig("gateway_config.yaml", "yaml")

	app := NewIrisApp(cfg)
	mvcApp := NewMvcApp(app, cfg)

	// 初始化 MVC 应用
	InitMvcApp(mvcApp, cfg)

	// 开始启动监听
	app.Listen(":" + strconv.Itoa(cfg.Server.Port))
}

// ========= 在 main 函数初始化加载流程中出现错误直接 panic ==========

// 加载网关服务器的配置
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

// 构造底层 Iris 应用
func NewIrisApp(cfg *config.Config) *iris.Application {
	app := iris.Default()

	// 设置日志级别
	// 注意：在 release 模式下，日志级别设置为 error
	// 默认状态下，日志级别为 info
	switch cfg.Server.Mode {
	case "debug":
		{
			app.Logger().SetLevel("debug")
		}
	case "release":
		{
			app.Logger().SetLevel("error")
		}
	default:
		{
			panic("未知的运行模式，请检查配置文件")
		}
	}

	return app
}

// 构造一个封装底层 Iris 应用的 MVC 应用
func NewMvcApp(irisApp *iris.Application, cfg *config.Config) *mvc.Application {
	// 直接包装整个根路由
	mvcApp := mvc.New(irisApp)

	// 设置 MVC 应用的配置
	mvcApp.Register(cfg)

	return mvcApp
}

// 由于初始化 MVC 应用较为复杂，单独分出一个函数
func InitMvcApp(root *mvc.Application, cfg *config.Config) {
	// 添加 /transfer 子路由组
	routes.ConfigureTransferRoutes(root)
}
