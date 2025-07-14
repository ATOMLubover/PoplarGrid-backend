package handlers

import (
	"poplargrid/internal/apiserver/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteUserHandler 注册用户相关的路由
func RouteUserHandler(root *mvc.Application) {
	// 创建 UserHandler 实例
	handler := &UserHandler{}

	// 注册路由组
	party := root.Party("/user")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Get("/detail", handler.UserDetail)

	{
		// 创建 UserMyHandler 实例
		myHandler := &UserMyHandler{}

		// 注册路由组
		myParty := party.Party("/my")

		// 注册 handler
		myParty.Handle(myHandler)

		// 注册路由与方法的映射（类型安全）
		myParty.Router.Get("/detail", myHandler.MyDetail)
		myParty.Router.Get("/teams", myHandler.MyTeams)
		myParty.Router.Get("/projects", myHandler.MyProjectList)
		myParty.Router.Get("/invitations_sent", myHandler.MyInvitationsSent)
		myParty.Router.Get("/invitations_received", myHandler.MyInvitationsReceived)
	}
}

// UserHandler 处理用户相关的请求
type UserHandler struct {
	UserService services.UserService
}

// UserDetail godoc
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、昵称等
// @Param 		user_id query int true "用户 ID"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} dtos.UserDetail
// @Router 		/user/detail [get]
func (h *UserHandler) UserDetail(ctx iris.Context) {
	// 从查询参数获取用户 ID
	userId, err := ctx.URLParamInt("user_id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "无效的用户 ID"})
		return
	}

	// 调用服务层获取用户详情
	userDetail, err := h.UserService.GetUserDetail(uint(userId))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取用户详情失败"})
		return
	}

	ctx.JSON(userDetail)
}

// UserMyHandler 处理与当前用户信息相关的路由
type UserMyHandler struct {
	UserService    services.UserService
	TeamService    services.TeamService
	ProjectService services.ProjectService
}

// MyDetail godoc
// @Summary 	获取当前用户的详细信息
// @Description 获取当前登录用户的详细信息，包括 ID、用户名、头像等
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} dtos.UserDetail
// @Router 		/user/my/detail [get]
func (h *UserMyHandler) MyDetail(ctx iris.Context) {
	// 从上下文获取当前用户 ID
	userId, err := ctx.Values().GetInt("user_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "无效的用户 ID"})
		return
	}

	// 调用服务层获取当前用户详情
	userDetail, err := h.UserService.GetUserDetail(uint(userId))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取用户详情失败"})
		return
	}

	ctx.JSON(userDetail)
}

// MyTeams godoc
// @Summary 	获取当前用户的汉化组列表
// @Description 获取当前登录用户所属的所有汉化组列表
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.TeamBasic
// @Router 		/user/my/teams [get]
func (h *UserMyHandler) MyTeams(ctx iris.Context) {
	// 从上下文获取当前用户 ID
	userId, err := ctx.Values().GetInt("user_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "无效的用户 ID"})
		return
	}

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
	teams, err := h.TeamService.GetTeamBasicPageByUserId(uint(userId), pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "获取汉化组列表失败"})
		return
	}

	// 返回结果
	ctx.JSON(teams)
}

// MyProjectList godoc
// @Summary 	获取当前用户的项目列表
// @Description 获取当前登录用户参与的所有项目列表，按照 ID 倒序排列
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.ProjectBasic
// @Router 		/user/my/project_list [get]
func (h *UserMyHandler) MyProjectList(ctx iris.Context) {
	// 从上下文获取当前用户 ID
	userId, err := ctx.Values().GetInt("user_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "无效的用户 ID"})
		return
	}

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
	projects, err := h.ProjectService.GetBasicPageByUserId(uint(userId), pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "获取项目列表失败"})
		return
	}

	ctx.JSON(projects)
}

// MyInvitationsSent godoc
// @Summary 	获取当前用户的邀请列表
// @Description 获取当前登录用户发出的所有邀请列表，仅包括项目邀请
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.InvitationBasic
// @Router 		/user/my/invitations_sent [get]
func (h *UserMyHandler) MyInvitationsSent(ctx iris.Context) {
}

// MyInvitationsReceived godoc
// @Summary 	获取当前用户的收到的邀请列表
// @Description 获取当前登录用户收到的所有邀请列表，仅包括项目邀请
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.InvitationBasic
// @Router 		/user/my/invitations_received [get]
func (h *UserMyHandler) MyInvitationsReceived(ctx iris.Context) {
}
