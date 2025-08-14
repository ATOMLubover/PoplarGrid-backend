package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12/mvc"
)

// FormatResponse 定义了统一的响应格式
type FormatResponse[T any] struct {
	ErrorCode int    `json:"error_code"`     // 错误码，不为 0 时表示发生错误
	Message   string `json:"message"`        // 错误信息或成功消息
	Data      T      `json:"data,omitempty"` // 成功时返回的数据，可能不携带
}

// StringFormatResponse 定义了字符串格式的响应
type StringFormatResponse FormatResponse[string]

// RouteAPIHandler 注册 API 相关的路由
func RouteAPIHandler(root *mvc.Application, tokenFactory services.AuthTokenFactory) {
	// 创建 /api 路由组
	api := root.Party("/api")

	// 设置全局中间件
	api.Router.Use(NewUserInfoExtractMiddleware(tokenFactory))

	// 注册各个 handler 的路由
	RouteUserHandler(api)
	RouteTeamHandler(api)
	RouteWorksetHandler(api)
	RouteMemberHandler(api)
	RouteInvitationHandler(api)
	RouteApplicationHandler(api)
	RouteCrawlerHandler(api)
}
