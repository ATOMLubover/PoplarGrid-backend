package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteTeamHandler 注册汉化组信息相关的路由
func RouteTeamHandler(root *mvc.Application) {
	// 注册路由组
	root.Party("/teams").
		Handle(new(TeamHandler))
}

// TeamHandler 处理汉化组信息相关的请求
type TeamHandler struct {
	TeamService services.TeamService
}

// TeamInfo 定义了汉化组的基本信息 DTO
type TeamInfo struct {
	// 汉化组 ID
	// @example 114514
	Id uint `json:"id"`
	// 汉化组名称
	// @example 白杨汉化组
	Name string `json:"name"`
	// 汉化组对应的龙译 ID
	// @example 64aac7d91b3
	MoetranId string `json:"moetran_id,omitempty"`
	// 汉化组描述
	// @example 这是一个测试汉化组
	Description string `json:"description,omitempty"`
}

// BeforeActivation 在控制器激活前注册路由
func (t *TeamHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "TeamList")
}

// TeamList godoc
//
// @Summary 	获取汉化组列表
// @Description 注意当列表为空，会返回 null 而不是空数组
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		user_id query string true "当前用户的 ID"
//
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []TeamInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/teams [get]
func (h *TeamHandler) TeamList(ctx iris.Context) {
	// 从上下文获取当前用户 ID
	memberId, err := ctx.Values().GetUint("user_id")
	if err != nil || memberId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 member_id",
		})
		return
	}

	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	// 调用服务层获取数据
	teams, err := h.TeamService.GetTeams(
		&services.TeamListParams{
			UserId: memberId,
			Offset: (pageSerial - 1) * pageSize,
			Limit:  pageSize,
		})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取当前汉化组列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(teams)
}

// // MemberListPage godoc
// // @Summary 	获取汉化组成员列表分页
// // @Description 注意当列表为空，会返回 null 而不是空数组
// // @Param 		page_serial query int false "页码，默认值为 1"
// // @Param 		page_size query int false "每页数量，默认值为 10"
// // @Param		member_nickname query string false "要模糊搜索的成员（部分）昵称"
// // @Param  		member_qq query int false "要搜索的成员的 QQ 号，使用 interger 加速查找"
// // @Param 		id path int true "所属汉化组 ID"
// // @Tags 		team
// // @Produce 	json
// // @Success	 	200 {object} []dtos.MemberBasic
// // @Failure     400 {object} ErrorResponse "无效的请求参数"
// // @Failure     500 {object} ErrorResponse "服务器内部错误"
// // @Router 		/teams/{id}/members [get]
// func (h *TeamHandler) MemberListPage(ctx iris.Context) {
// 	// 获取分页参数
// 	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
// 	pageSize := ctx.URLParamIntDefault("page_size", 10)

// 	teamId, err := ctx.Params().GetUint("id")
// 	if err != nil || teamId <= 0 {
// 		ctx.StatusCode(iris.StatusBadRequest)
// 		ctx.JSON(ErrorResponse{
// 			Error: "无法获取有效的 team_id",
// 		})
// 		return
// 	}

// 	// 获取昵称残片条件，默认为空代表不查找
// 	nickname := ctx.URLParam("member_nickname")

// 	// 获取 QQ 号条件，默认为 -1 代表不查找
// 	qqNumber := ctx.URLParamInt32Default("member_qq", -1)

// 	// 调用服务层获取数据
// 	members, err := h.TeamService.GetMemberBasicPageWithParams(
// 		teamId, pageSerial, pageSize, nickname, int(qqNumber))
// 	if err != nil {
// 		ctx.StatusCode(iris.StatusInternalServerError)
// 		ctx.JSON(ErrorResponse{
// 			Error:  "获取特定汉化组的成员基础信息列表失败",
// 			Detail: err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(members)
// }

// // WorksetListPage godoc
// // @Summary 	获取特定汉化组的工作集列表分页，按 ID 倒序
// // @Description 注意当列表为空，会返回 null 而不是空数组
// // @Param 		page_serial query int false "页码，默认值为 1"
// // @Param 		page_size query int false "每页数量，默认值为 10"
// // @Param 		id query uint true "所属汉化组 ID"
// // @Tags 		team
// // @Produce 	json
// // @Success	 	200 {object} []dtos.WorksetBasic
// // @Failure     400 {object} ErrorResponse "无效的请求参数"
// // @Failure     500 {object} ErrorResponse "服务器内部错误"
// // @Router 		/teams/{id}/worksets [get]
// func (h *TeamHandler) WorksetListPage(ctx iris.Context) {
// 	// 获取分页参数
// 	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
// 	pageSize := ctx.URLParamIntDefault("page_size", 10)

// 	teamId, err := ctx.Params().GetUint("id")
// 	if err != nil || teamId <= 0 {
// 		ctx.StatusCode(iris.StatusBadRequest)
// 		ctx.JSON(ErrorResponse{
// 			Error: "team_id 必须是一个明确给出的正整数",
// 		})
// 		return
// 	}

// 	// 调用服务层获取数据
// 	worksets, err := h.WorksetService.GetBasicPage(teamId, pageSerial, pageSize)
// 	if err != nil {
// 		ctx.StatusCode(iris.StatusInternalServerError)
// 		ctx.JSON(ErrorResponse{
// 			Error:  "获取工作集列表失败",
// 			Detail: err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(worksets)
// }
