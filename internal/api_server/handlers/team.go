package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// TeamInfo 定义了汉化组的基本信息 DTO
type TeamInfo struct {
	// 汉化组 ID
	Id uint `json:"id"`
	// 汉化组名称
	Name string `json:"name"`
	// 汉化组对应的龙译 ID
	MoetranId string `json:"moetran_id,omitempty"`
	// 汉化组描述
	Description string `json:"description,omitempty"`
}

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

// BeforeActivation 在控制器激活前注册路由
func (t *TeamHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
}

// List godoc
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
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/teams [get]
func (h *TeamHandler) List(ctx iris.Context) {
	// 从上下文获取当前用户 ID
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
	teams, err := h.TeamService.GetTeams(&services.TeamListParams{
		UserId: userId,
		Offset: (pageSerial - 1) * pageSize,
		Limit:  pageSize,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "获取当前汉化组列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(teams)
}
