package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteWorksetHandler 注册工作集相关的路由
func RouteWorksetHandler(root *mvc.Application) {
	// 创建 WorksetHandler 实例
	handler := &WorksetHandler{}

	// 注册路由组
	root.Party("/worksets").
		Handle(handler)
}

// WorksetHandler 处理工作集相关的请求
type WorksetHandler struct {
	WorksetService services.WorksetService
}

// BeforeActivation 在控制器激活前注册路由
func (w *WorksetHandler) BeforeActivation(b mvc.BeforeActivation) {
	// 注册 GET /worksets/{id:uint}/stats
	b.Handle("GET", "/{id:uint}/stats", "ProjectStats")
}

// ProjectStats godoc
// @Summary 	获取特定作品集项目统计信息
// @Description 获取所有项目的统计信息，包括总数、翻译、校对、嵌字，审核、发布对应数量等
// @Param 		id path int true "作品集 ID"
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectStats
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/worksets/{id}/stats [get]
func (h *WorksetHandler) ProjectStats(ctx iris.Context) {
	// 获取 workset_id 参数
	worksetId, err := ctx.Params().GetUint("id")
	if err != nil || worksetId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "workset_id 必须是一个明确给出的正整数",
		})
		return
	}

	// 调用服务层获取统计数据
	stats, err := h.WorksetService.GetWorksetStats(worksetId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ctx.JSON(ErrorResponse{
			Error:  "获取特定作品集的项目统计信息失败",
			Detail: err.Error(),
		}))
		return
	}

	ctx.JSON(stats)
}
