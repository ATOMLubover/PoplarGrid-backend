package handlers

import (
	"poplargrid/internal/apiserver/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteTeamHandler 注册汉化组信息相关的路由
func RouteTeamHandler(root *mvc.Application) {
	// 创建 TeamHandler 实例
	handler := &TeamHandler{}

	// 注册路由组
	party := root.Party("/team")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Get("/list", handler.TeamListPage)
}

// TeamHandler 处理汉化组信息相关的请求
type TeamHandler struct {
	TeamService services.TeamService
}

// TeamListPage godoc
// @Summary 	获取汉化组列表分页，按 ID 倒序
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []dtos.TeamBasic
// @Router 		/team/list [get]
func (h *TeamHandler) TeamListPage(ctx iris.Context) {
	// 获取分页参数
	pageSerial, _ := ctx.URLParamInt("page_serial")
	if pageSerial <= 0 {
		pageSerial = 1
	}
	pageSize, _ := ctx.URLParamInt("page_size")
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务层获取数据
	teams, err := h.TeamService.GetBasicPage(pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "获取汉化组列表失败"})
		return
	}

	// 返回结果
	ctx.JSON(teams)
}
