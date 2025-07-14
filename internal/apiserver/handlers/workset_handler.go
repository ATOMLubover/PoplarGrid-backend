package handlers

import (
	"poplargrid/internal/apiserver/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteWorksetHandler 注册工作集相关的路由
func RouteWorksetHandler(root *mvc.Application) {
	// 创建 WorksetHandler 实例
	handler := &WorksetHandler{}

	// 注册路由组
	party := root.Party("/workset")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Get("/list", handler.WorksetListPage)
}

// WorksetHandler 处理工作集相关的请求
type WorksetHandler struct {
	WorksetService services.WorksetService
}

// WorksetListPage godoc
// @Summary 	获取特定汉化组的工作集列表分页，按 ID 倒序
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		team_id query int true "所属汉化组 ID"
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} []dtos.WorksetBasic
// @Router 		/workset/list [get]
func (h *WorksetHandler) WorksetListPage(ctx iris.Context) {
	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)

	pageSize := ctx.URLParamIntDefault("page_size", 10)

	teamId, err := ctx.URLParamInt("team_id")
	if err != nil || teamId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "team_id 必须是一个明确给出的正整数"})
		return
	}

	// 调用服务层获取数据
	worksets, err := h.WorksetService.GetBasicPage(pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取工作集列表失败"})
		return
	}

	ctx.JSON(worksets)
}

// ProjectStats godoc
// @Summary 	获取特定作品集项目统计信息
// @Description 获取所有项目的统计信息，包括总数、翻译、校对、嵌字，审核、发布对应数量等
// @Param 		workset_id query int true "作品集 ID，必填"
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectStats
// @Router 		/workset/stats [get]
func (h *WorksetHandler) ProjectStats(ctx iris.Context) {
	// 获取 workset_id 参数
	worksetId, err := ctx.URLParamInt("workset_id")
	if err != nil || worksetId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "workset_id 必须是一个明确给出的正整数"})
		return
	}

	// 调用服务层获取统计数据
	stats, err := h.WorksetService.GetProjectStats(uint(worksetId))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取项目总体统计信息失败"})
		return
	}

	ctx.JSON(stats)
}
