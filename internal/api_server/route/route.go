package route

import (
	"poplargrid/internal/api_server/handler"

	"github.com/kataras/iris/v12/mvc"
)

// ConfigureMemberRoutes 配置 MemberHandler 相关的路由
func ConfigureMemberRoutes(app *mvc.Application) {
	// 创建 /member 子路由组
	memberParty := app.Party("/member")

	// 绑定指定路由（使用类型安全的方法）
	memberHandlerInst := new(handler.MemberHandler)
	memberParty.Handle(memberHandlerInst)
	{
		// 注册测试路由
		memberParty.Router.Get("/list",
			memberHandlerInst.MemberListPage)
	}
}
