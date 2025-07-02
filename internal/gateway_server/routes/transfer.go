package routes

import (
	"poplargrid/internal/gateway_server/handlers"

	"github.com/kataras/iris/v12/mvc"
)

// 处理 app，添加 /transfer 子路由组
func ConfigureTransferRoutes(app *mvc.Application) {
	// 创建 /transfer 子路由组
	transferParty := app.Party("/transfer")

	// 绑定指定路由（使用类型安全的方法）
	transferHandlerInst := new(handlers.TransferHandler)
	transferParty.Handle(transferHandlerInst)
	{
		// 注册 ping 路由
		transferParty.Router.Get("/try_ping", transferHandlerInst.HealthCheck)
	}
}
