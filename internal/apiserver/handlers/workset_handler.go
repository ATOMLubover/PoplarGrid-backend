package handlers

import (
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
}

// WorksetListPage godoc
// @Summary 	获取工作集列表分页，按 ID 倒序
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		team_id query int true "所属汉化组 ID"
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} []dtos.WorksetBasic
// @Router 		/workset/list [get]
func (h *WorksetHandler) WorksetListPage(ctx iris.Context) {
}
