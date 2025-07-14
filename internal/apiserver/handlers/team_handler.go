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
	party.Router.Get("/member_list", handler.MemberListPage)
}

// TeamHandler 处理汉化组信息相关的请求
type TeamHandler struct {
	TeamService services.TeamService
}

// MemberListPage godoc
// @Summary 	获取汉化组成员列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		team_id query int true "所属汉化组 ID，必填"
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []dtos.MemberBasic
// @Router 		/team/member_list [get]
func (h *TeamHandler) MemberListPage(ctx iris.Context) {
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
	members, err := h.TeamService.GetMemberBasicPage(uint(teamId), pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取成员列表失败"})
		return
	}

	ctx.JSON(members)
}
