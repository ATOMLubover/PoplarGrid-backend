package handlers

import (
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteLaborProcHandler 注册邀请、申请相关的路由
func RouteLaborProcHandler(root *mvc.Application) {
	// 创建 InvitationHandler 实例
	invHandler := &InvitationHandler{}

	// 注册路由组
	invParty := root.
		Party("/invitations").
		Handle(invHandler)

	// 注册路由与方法的映射（类型安全）
	invParty.Router.Get("", invHandler.InvitationListPage)
	invParty.Router.Post("", invHandler.CreateInvitation)

	// 创建 ApplicationHandler 实例
	appHandler := &ApplicationHandler{}

	// 注册路由组
	appParty := root.
		Party("/applications").
		Handle(appHandler)

	// 注册路由与方法的映射（类型安全）
	appParty.Router.Get("", appHandler.ApplicationListPage)
	appParty.Router.Post("", appHandler.CreateApplication)
}

// InvitationHandler 处理邀请相关的请求
type InvitationHandler struct {
	LaborService services.LaborService
}

// InvitationListPage godoc
// @Summary 获取当前用户的邀请（发出或者收到）列表，支持分页
// @Description 根据分页参数获取用户的邀请列表，支持分页和排序。当列表为空时，会返回 null 而不是空数组。
// @Param page_serial query int false "页码，默认值为 1"
// @Param page_size query int false "每页数量，默认值为 10"
// @Param kind query string true "邀请类型，0：发送的邀请，1：收到的邀请"
// @Tags invitation
// @Produce json
// @Success 200 {object} []dtos.InvitationBasic
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /invitations [get]
func (h *InvitationHandler) InvitationListPage(ctx iris.Context) {
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
		invitations, err = h.LaborService.GetInvitationSentByUserId(userId, pageSerial, pageSize)
	case "1": // 收到的邀请
		invitations, err = h.LaborService.GetInvitationRecievedByUserId(userId, pageSerial, pageSize)
	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 kind 参数",
		})
		return
	}

	ctx.JSON(invitations)
}

// CreateInvitation godoc
// @Summary 创建一个新的邀请
// @Description 创建一个新的邀请，邀请用户加入项目
// @Param body_params body dtos.InviteMemberRequest true "邀请参数"
// @Tags invitation
// @Produce json
// @Success 200 {object} SuccessResponse "邀请创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /invitations [post]
func (h *InvitationHandler) CreateInvitation(ctx iris.Context) {
	// 从上下文中获取当前用户 ID
	inviterId, err := ctx.Values().GetUint("user_id")
	if err != nil || inviterId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 inviter_id",
		})
		return
	}

	// 读取请求体中的邀请参数
	var req dtos.InviteMemberRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 调用服务层创建邀请
	if err := h.LaborService.CreateInvitation(
		inviterId, req.MemberId,
		req.ProjectId, req.InviteRole); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "创建邀请失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "邀请创建成功",
	})
}

// ProcessInvitation godoc
// @Summary 处理邀请
// @Description 处理邀请，接受或拒绝邀请
// @Param body_params body dtos.ProcessInvitationRequest true "处理邀请参数"
// @Param id path integer true "邀请 ID"
// @Tags invitation
// @Produce json
// @Success 200 {object} SuccessResponse "邀请处理成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /invitations/{id} [post]
func (h *InvitationHandler) ProcessInvitation(ctx iris.Context) {
	// 从上下文中获取当前用户 ID
	inviteeId, err := ctx.Values().GetUint("user_id")
	if err != nil || inviteeId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 invitee_id",
		})
		return
	}

	// 读取请求体中的处理邀请参数
	var req dtos.ProcessInvitationRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 如果路径 ID 与请求体中的邀请 ID 不匹配，则返回错误
	invitationId, err := ctx.Params().GetUint("id")
	if err != nil || invitationId <= 0 || invitationId != req.InvitationId {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 invitation_id",
		})
		return
	}

	// 调用服务层处理邀请
	status := 0
	switch req.Accept {
	case true:
		status = 1 // 接受邀请
	case false:
		status = 2 // 拒绝邀请
	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 accept 参数",
		})
		return
	}

	if err := h.LaborService.ProcessInvitation(
		inviteeId, req.InvitationId, status); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "处理邀请失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "邀请处理成功",
	})
}

// ApplicationHandler 处理申请相关的请求
type ApplicationHandler struct {
	LaborService services.LaborService
}

// ApplicationListPage godoc
// @Summary 获取当前用户的申请（发出或收到）列表，支持分页
// @Description 根据分页参数获取用户的申请列表，支持分页和排序。当列表为空
// @Param page_serial query int false "页码，默认值为 1"
// @Param page_size query int false "每页数量，默认值为 10"
// @Param kind query string true "申请类型，0：发出的申请，1：收到的申请"
// @Tags application
// @Produce json
// @Success 200 {object} []dtos.ApplicationBasic
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /applications [get]
func (h *ApplicationHandler) ApplicationListPage(ctx iris.Context) {
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
		// 获取发出的申请
		applications, err = h.LaborService.GetApplisSentByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取发出的申请列表失败",
				Detail: err.Error(),
			})
			return
		}
	case "1":
		// 获取收到的申请
		applications, err = h.LaborService.GetApplisRecievedByUserId(userId, pageSerial, pageSize)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(ErrorResponse{
				Error:  "获取收到的申请列表失败",
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

// CreateApplication godoc
// @Summary 创建一个新的申请
// @Description 创建一个新的申请，申请加入项目
// @Param body_params body dtos.ApplyProjectRequest true "申请参数"
// @Tags application
// @Produce json
// @Success 200 {object} SuccessResponse "申请创建成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(ctx iris.Context) {
	// 从上下文中获取当前用户 ID
	userId, err := ctx.Values().GetUint("user_id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 读取请求体中的申请参数
	var req dtos.ApplyProjectRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 调用服务层创建申请
	if err := h.LaborService.CreateApplication(userId, req.ProjectId, req.ApplyRole); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "创建申请失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "申请创建成功",
	})
}

// ProcessApplication godoc
// @Summary 处理申请
// @Description 处理申请，接受或拒绝申请
// @Param body_params body dtos.ProcessApplicationRequest true "处理申请参数"
// @Param id path integer true "申请 ID"
// @Tags application
// @Produce json
// @Success 200 {object} SuccessResponse "申请处理成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /applications/{id} [post]
func (h *ApplicationHandler) ProcessApplication(ctx iris.Context) {
	// 从上下文中获取当前用户 ID
	userId, err := ctx.Values().GetUint("user_id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 读取请求体中的处理申请参数
	var req dtos.ProcessApplicationRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 如果路径 ID 与请求体中的申请 ID 不匹配，则返回错误
	applicationId, err := ctx.Params().GetUint("id")
	if err != nil || applicationId <= 0 || applicationId != req.ApplicationId {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 application_id",
		})
		return
	}

	// 调用服务层处理申请
	status := 0
	switch req.Accept {
	case true:
		status = 1 // 接受申请
	case false:
		status = 2 // 拒绝申请
	default:
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的 accept 参数",
		})
		return
	}

	if err := h.LaborService.ProcessApplication(userId, req.ApplicationId, status); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "处理申请失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "申请处理成功",
	})
}
