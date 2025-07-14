package handlers

import (
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
	party.Router.Get("/list", handler.UserListPage)
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
		myParty.Router.Get("/projects", myHandler.MyProjects)
	}
}

// UserHandler 处理用户相关的请求
type UserHandler struct {
}

// UserListPage godoc
// @Summary 	获取用户列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_asc | id_desc，其他输入无效"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} []dtos.UserBasic
// @Router 		/user/list [get]
func (h *UserHandler) UserListPage(ctx iris.Context) {
}

// UserDetail godoc
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、用户名、头像等
// @Param 		user_id query int true "用户 ID"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} dtos.UserDetail
// @Router 		/user/detail [get]
func (h *UserHandler) UserDetail(ctx iris.Context) {
}

// UserMyHandler 处理与当前用户信息相关的路由
type UserMyHandler struct {
}

// MyDetail godoc
// @Summary 	获取当前用户的详细信息
// @Description 获取当前登录用户的详细信息，包括 ID、用户名、头像等
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} dtos.UserDetail
// @Router 		/user/my/detail [get]
func (h *UserMyHandler) MyDetail(ctx iris.Context) {
}

// MyTeams godoc
// @Summary 	获取当前用户的汉化组列表
// @Description 获取当前登录用户所属的所有汉化组列表
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.TeamBasic
// @Router 		/user/my/teams [get]
func (h *UserMyHandler) MyTeams(ctx iris.Context) {
}

// MyProjects godoc
// @Summary 	获取当前用户的项目列表
// @Description 获取当前登录用户参与的所有项目列表
// @Tags 		user_my
// @Produce 	json
// @Success	 	200 {object} []dtos.ProjectBasic
// @Router 		/user/my/projects [get]
func (h *UserMyHandler) MyProjects(ctx iris.Context) {
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
