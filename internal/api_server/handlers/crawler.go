package handlers

import (
	"poplargrid/internal/api_server/config"
	"poplargrid/internal/api_server/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteCrawlerHandler 注册爬虫相关的路由
func RouteCrawlerHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/crawler").
		Handle(new(CrawlerHandler))
}

// CrawlerHandler 处理爬虫相关的请求
type CrawlerHandler struct {
	CrawlerService services.CrawlerService
	MemberService  services.MemberService
}

// BeforeActivation 在控制器激活前注册路由
func (h *CrawlerHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/auto-update-all", "AutoUpdateAll")
}

// AutoUpdateAll godoc
// @Summary 	更新汉化组的项目信息
// @Description 自动递归地更新当前用户所有汉化组的项目信息，可能会需要较长时间
//
// @Tags 		crawler
// @Produces 	json
// @Success     200 {object} SuccessResponse "更新成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {string} string "服务器内部错误"
//
// @Router /api/crawler/auto-update-all [post]
func (h *CrawlerHandler) AutoUpdateAll(ctx iris.Context) {
	// 从上下文中获取 moetranAuth 和 userID
	moetranAuth := ctx.Values().GetString("moetran_jwt")
	if moetranAuth == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效的 moetran_auth",
		})
		return
	}

	userID, err := ctx.Values().GetUint("user_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取用户 ID",
		})
		return
	}

	// 调用服务层进行自动更新
	if err := h.CrawlerService.AutoUpdateAll(userID, moetranAuth); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "自动更新失败",
			Detail: err.Error(),
		})
		return
	}

	// 随后更新 memberIDs 的 Cookie
	newToken, err := h.MemberService.UpdateToken(userID, moetranAuth)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "拉取尨译信息成功，更新 Token 失败",
			Detail: err.Error(),
		})
		return
	}

	// 设置新的 Token Cookie
	// 为响应添加 Set-Cookie 头部
	expiresAt := time.Now().Add(7 * time.Hour) // 默认为 7 天
	if cfg := config.GetConfig(); cfg != nil && cfg.Server.CookieLifetime <= 0 {
		expiresAt = time.Now().Add(time.Duration(cfg.Server.CookieLifetime) * time.Second)
	}

	ctx.SetCookie(&iris.Cookie{
		Name:     "poplar_token",
		Value:    newToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
		// TODO: 测试环境不启用 HTTPS
		// Secure:   true,
	})

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(SuccessResponse{
		Message: "自动更新成功",
	})
}
