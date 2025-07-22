package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteTeamHandler 注册汉化组信息相关的路由
func RouteTeamHandler(root *mvc.Application) {
	// 创建 TeamHandler 实例
	handler := &TeamHandler{}

	// 注册路由组
	teamParty := root.
		Party("/teams").
		Handle(handler)

	// 注册路由与方法的映射（类型安全）
	teamParty.Router.Get("", handler.TeamListPage)
	teamParty.Router.Get("/{id:uint}/members", handler.MemberListPage)
}

// TeamHandler 处理汉化组信息相关的请求
type TeamHandler struct {
	TeamService services.TeamService
}

// TeamListPage godoc
// @Summary 	获取当前用户的汉化组列表分页
// @Description 根据分页参数获取汉化组列表，支持分页和排序。当列表为空时，会返回 null 而不是空数组。
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []dtos.TeamBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/teams [get]
func (h *TeamHandler) TeamListPage(ctx iris.Context) {
	// 从上下文中获取当前用户 ID
	userId, err := ctx.Values().GetUint("user_id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	// 调用服务层获取数据
	teams, err := h.TeamService.GetBasicPageByUserId(userId, pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取汉化组列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(teams)
}

// MemberListPage godoc
// @Summary 	获取汉化组成员列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param		member_nickname query string false "要模糊搜索的成员（部分）昵称"
// @Param  		member_qq query int false "要搜索的成员的 QQ 号，使用 interger 加速查找"
// @Param 		id path int true "所属汉化组 ID"
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []dtos.MemberBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/teams/{id}/members [get]
func (h *TeamHandler) MemberListPage(ctx iris.Context) {
	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	teamId, err := ctx.Params().GetUint("id")
	if err != nil || teamId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 team_id",
		})
		return
	}

	// 获取昵称残片条件，默认为空代表不查找
	nickname := ctx.URLParam("member_nickname")

	// 获取 QQ 号条件，默认为 -1 代表不查找
	qqNumber := ctx.URLParamInt32Default("member_qq", -1)

	// 调用服务层获取数据
	members, err := h.TeamService.GetMemberBasicPageWithParams(
		teamId, pageSerial, pageSize, nickname, int(qqNumber))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取特定汉化组的成员基础信息列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(members)
}
