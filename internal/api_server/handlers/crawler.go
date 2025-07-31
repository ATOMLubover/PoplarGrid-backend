package handlers

import (
	"poplargrid/internal/api_server/services"

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
	// 从上下文中获取 moetranAuth
	moetranAuth := ctx.Values().GetString("moetranAuth")
	if moetranAuth == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效的 moetran_auth",
		})
		return
	}

	// 调用服务层进行自动更新
	if err := h.CrawlerService.AutoUpdateAll(moetranAuth); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "自动更新失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(SuccessResponse{
		Message: "自动更新成功",
	})
}
