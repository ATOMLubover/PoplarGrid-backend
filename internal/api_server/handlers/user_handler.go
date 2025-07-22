package handlers

import (
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteUserHandler 注册用户相关的路由
func RouteUserHandler(root *mvc.Application) {
	// 创建 UserHandler 和 UserMyHandler 实例
	userHandler := &UserHandler{}

	// 注册路由组和 handler
	userParty := root.
		Party("/users").
		Handle(userHandler)

	// 注册中间件，检查用户 ID 是否匹配当前登录用户
	userParty.Router.Use(NewCheckUserIdMiddleware())

	// 注册路由与方法的映射（类型安全）
	userParty.Router.Get("/{id:uint}/teams", userHandler.TeamListPage)
	userParty.Router.Get("/{id:uint}/invitations", userHandler.InvitationListPage)
	userParty.Router.Get("/{id:uint}/detail", userHandler.UserDetail)
	userParty.Router.Get("/{id:uint}/projects", userHandler.ProjectListPage)
}

// UserHandler 处理用户相关的请求
type UserHandler struct {
	UserService    services.UserService
	TeamService    services.TeamService
	ProjectService services.ProjectService
	InvAppService  services.LaborService
}

// func (h *UserHandler) BeforeActivation(b mvc.BeforeActivation) {
// 	b.Handle("GET", "/{id:uint}/teams", "TeamListPage")
// 	b.Handle("GET", "/{id:uint}/invitations", "InvitationListPage")
// 	b.Handle("GET", "/{id:uint}/detail", "UserDetail")
// 	b.Handle("GET", "/{id:uint}/projects", "ProjectListPage")

// 	// 注册控制器特定的中间件
// 	b.Router().Use(NewCheckUserIdMiddleware())
// }

// UserDetail godoc
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、昵称、QQ 等
// @Param 		id path uint true "用户 ID"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} dtos.UserDetail
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/users/{id}/detail [get]
func (h *UserHandler) UserDetail(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 调用服务层获取用户详情
	userDetail, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取特定用户详情失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(userDetail)
}

// TeamListPage godoc
// @Summary 	获取当前用户参与的汉化组列表
// @Description 获取指定用户参与的所有汉化组列表，支持分页和排序
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		id path uint true "用户 ID"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} []dtos.TeamBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/users/{id}/teams [get]
func (h *UserHandler) TeamListPage(ctx iris.Context) {
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

	// 调用服务层获取用户参与的汉化组列表
	teams, err := h.TeamService.GetBasicPageByUserId(userId, pageSerial, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取用户参与的汉化组列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(teams)
}

// ProjectListPage godoc
// @Summary 	获取用户参与的项目列表
// @Description 获取指定用户参与的所有项目列表，支持分页和复合查询
// @Param       page_serial query integer false "页码，默认值为 1"
// @Param       page_size query integer false "每页数量，默认值为 10"
// @Param       status query integer false "项目状态（位掩码），用于复合查询，默认不筛选查询"
// @Param       workset_id query integer false "项目所属的作品集 ID，默认不筛选作品集"
// @Param 		id path uint true "用户 ID"
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} []dtos.ProjectBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/users/{id}/projects [get]
func (h *UserHandler) ProjectListPage(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 user_id",
		})
		return
	}

	// 获取分页参数
	pageSerial := ctx.URLParamInt32Default("page_serial", 1)
	pageSize := ctx.URLParamInt32Default("page_size", 10)

	// 注意默认为 PROJECT_STATUS_ALL，表示不筛选状态
	status := ctx.URLParamInt32Default("status", dtos.PROJECT_STATUS_ALL)

	// 注意默认为 0，代表不筛选 workset
	worksetId := ctx.URLParamIntDefault("workset_id", 0)

	// 调用服务层获取数据
	projects, err := h.ProjectService.GetBasicPageWithParams(
		uint(worksetId), int(pageSerial), int(pageSize),
		dtos.SORT_ID_DESC, userId, dtos.ProjectOverallStatus(status))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取用户参与的项目列表失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(projects)
}

// InvitationListPage godoc
// @Summary 获取当前用户的邀请（发出或者收到）列表，支持分页
// @Description 根据分页参数获取用户的邀请列表，支持分页和排序\n当列表为空时，会返回 null 而不是空数组\n如果 id 不是当前登录的用户 ID，则返回 400 错误
// @Param page_serial query int false "页码，默认值为 1"
// @Param page_size query int false "每页数量，默认值为 10"
// @Param id path uint true "用户 ID"
// @Param kind query string true "邀请类型，0：发送的邀请，1：收到的邀请"
// @Tags invitation
// @Produce json
// @Success 200 {object} []dtos.InvitationBasic
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /users/{id}/invitations [get]
func (h *UserHandler) InvitationListPage(ctx iris.Context) {
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
	var invitations []*dtos.InvitationBasic

	switch ctx.URLParam("kind") {
	case "0": // 发送的邀请
		invitations, err = h.InvAppService.GetInvitationSentByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取用户发送的邀请列表失败",
				Detail: err.Error(),
			})
			return
		}

	case "1": // 收到的邀请
		invitations, err = h.InvAppService.GetInvitationRecievedByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取用户收到的邀请列表失败",
				Detail: err.Error(),
			})
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 kind 参数",
		})
		return
	}

	ctx.JSON(invitations)
}

// ApplicationSentListPage godoc
// @Summary 获取当前用户的申请（发出或收到）列表，支持分页
// @Description 根据分页参数获取用户的申请列表，支持分页和排序\n当列表为空时，会返回 null 而不是空数组\n如果 id 不是当前登录的用户 ID，则返回 400 错误
// @Param page_serial query int false "页码，默认值为 1"
// @Param page_size query int false "每页数量，默认值为 10"
// @Param id path uint true "用户 ID"
// @Param kind query string true "申请类型，0：发送的申请，1：收到的申请"
// @Tags application
// @Produce json
// @Success 200 {object} []dtos.ApplicationBasic
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /users/{id}/applications [get]
func (h *UserHandler) ApplicationSentListPage(ctx iris.Context) {
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
	var applications []*dtos.ApplicationBasic

	switch ctx.URLParam("kind") {
	case "0":
		// 获取发送的申请
		applications, err = h.InvAppService.GetApplisSentByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取用户发送的申请列表失败",
				Detail: err.Error(),
			})
			return
		}

	case "1":
		// 获取收到的申请
		applications, err = h.InvAppService.GetApplisRecievedByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取用户收到的申请列表失败",
				Detail: err.Error(),
			})
			return
		}

	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 kind 参数",
		})
		return
	}

	ctx.JSON(applications)
}
