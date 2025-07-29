package handlers

import (
	"errors"
	"poplargrid/internal/api_server/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// InvitationInfo 定义了邀请的基本信息 DTO
type InvitationInfo struct {
	// 邀请 ID
	Id uint `json:"id"`
	// 邀请时间
	Time string `json:"time"`
	// 邀请者成员 ID
	InvitorMemberId uint `json:"invitor_member_id"`
	// 邀请者成员昵称
	InvitorNickname string `json:"invitor_nickname"`
	// 接收者成员 ID
	InviteeMemberId uint `json:"invitee_member_id"`
	// 接收者成员昵称
	InviteeNickname string `json:"invitee_nickname"`
	// 目标项目 ID
	TargetProjectId uint `json:"target_project_id"`
	// 目标项目名称
	TargetProjectName string `json:"target_project_name"`
	// 目标分工掩码
	TargetLaborMask LaborMask `json:"target_labor_mask"`
}

// CreateInvitationRequest 定义了创建邀请的请求参数
type CreateInvitationRequest struct {
	// 邀请者成员 ID
	InvitorMemberId uint `json:"invitor_member_id" binding:"required"`
	// 接收者成员 ID
	InviteeMemberId uint `json:"invitee_member_id" bindging:"required"`
	// 目标项目 ID
	TargetProjectId uint `json:"target_project_id" binding:"required"`
	// 目标分工掩码
	TargetLaborMask services.LaborMask `json:"target_labor_mask" binding:"required"`
}

// ProcessInvitationRequest 定义了处理邀请的请求参数
type ProcessInvitationRequest struct {
	// 邀请 ID
	InvitationId uint `json:"invitation_id" binding:"required"`
	// 处理者成员 ID
	ProcessorMemberId uint `json:"processor_member_id" binding:"required"`
	// 是否接受邀请
	Accept bool `json:"accept" binding:"required"`
}

// RouteInvitationHandler 注册邀请相关的路由
func RouteInvitationHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/invitations").
		Handle(new(InvitationHandler))
}

// InvitationHandler 处理邀请相关的请求
type InvitationHandler struct {
	invitationService services.InvitationService
}

// BeforeActivation 在控制器激活前注册路由
func (h *InvitationHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
	b.Handle("POST", "/", "Create")
	b.Handle("PUT", "/{id:uint}", "Process")
}

// List godoc
// @Summary 	获取邀请列表
// @Description 获取当前用户的邀请列表，通过参数的选择来确定到底是返回邀请者的邀请还是被邀请者的邀请
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		invitor_member_id query int false "邀请人的成员 ID，与被邀请人的成员 ID 只能二选一"
// @Param 		invitee_member_id query int false "被邀请人的成员 ID，与邀请人的成员 ID 只能二选一"
// @Param 		target_project_id query int false "目标项目 ID，如果需要筛选特定项目的邀请，可以使用此参数"
//
// @Tags 		invitation
// @Produce 	json
// @Success 	200 {object} []InvitationInfo
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/invitations [get]
func (h *InvitationHandler) List(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	// 获取查询参数
	invitorMemberId := ctx.URLParamIntDefault("invitor_member_id", 0)
	inviteeMemberId := ctx.URLParamIntDefault("invitee_member_id", 0)
	targetProjectId := ctx.URLParamIntDefault("target_project_id", 0)

	// 检查参数的组合是否合法
	if err := checkInvitationListQueryParams(invitorMemberId, inviteeMemberId, targetProjectId); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "查询参数组合非法",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务获取邀请列表
	invitations, err := h.invitationService.GetInvitations(&services.InvitationListParams{
		Offset:           (pageSerial - 1) * pageSize,
		Limit:            pageSize,
		InvitorMemberId:  uint(invitorMemberId),
		InviteeMemberId:  uint(inviteeMemberId),
		TargetProjectId:  uint(targetProjectId),
		CurrentMemberIds: memberIds,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取邀请列表失败",
			Detail: err.Error(),
		})
		return
	}

	// 将邀请列表转换为 DTO
	invitationInfos := make([]InvitationInfo, 0, len(invitations))

	for _, inv := range invitations {
		invitationInfos = append(invitationInfos, InvitationInfo{
			Id:                inv.Id,
			Time:              inv.Time.Format(time.DateTime),
			InvitorMemberId:   inv.InvitorMemberId,
			InvitorNickname:   inv.InvitorNickname,
			InviteeMemberId:   inv.InviteeMemberId,
			InviteeNickname:   inv.InviteeNickname,
			TargetProjectId:   inv.TargetProjectId,
			TargetProjectName: inv.TargetProjectTitle,
			TargetLaborMask:   inv.TargetLaborMask,
		})
	}

	ctx.JSON(invitationInfos)
}

// Create godoc
// @Summary 	创建邀请
// @Description 创建一个新的邀请记录
//
// @Accept 		json
// @Param 		body_params body CreateInvitationRequest true "邀请的基本信息"
//
// @Tags 		invitation
// @Produce 	json
// @Success 	200 {object} SuccessResponse "创建成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/invitations [post]
func (h *InvitationHandler) Create(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	var req CreateInvitationRequest

	// 绑定请求体参数
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "无效的请求参数",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务创建邀请
	if err := h.invitationService.CreateInvitation(&services.CreateInvitationParams{
		InvitorMemberId:  req.InvitorMemberId,
		InviteeMemberId:  req.InviteeMemberId,
		TargetProjectId:  req.TargetProjectId,
		TargetLaborMask:  req.TargetLaborMask,
		CurrentMemberIds: memberIds,
	}); err != nil {
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

// Process godoc
// @Summary 	处理邀请
// @Description 接受或拒绝邀请
//
// @Accept 		json
// @Param 		id path uint true "邀请 ID"
// @Param 		body_params body ProcessInvitationRequest true "处理邀请的请求参数"
//
// @Tags 		invitation
// @Produce 	json
// @Success 	200 {object} SuccessResponse "处理成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/invitations/{id} [put]
func (h *InvitationHandler) Process(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	var req ProcessInvitationRequest

	// 绑定请求体参数
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "无效的请求参数",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务处理邀请
	if err := h.invitationService.ProcessInvitation(&services.ProcessInvitationParams{
		InvitationId:      req.InvitationId,
		ProcessorMemberId: req.ProcessorMemberId,
		Accept:            req.Accept,
		CurrentMemberIds:  memberIds,
	}); err != nil {
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

// checkInvitationListQueryParams 检查邀请列表查询参数的合法性
func checkInvitationListQueryParams(invitorMemberId, inviteeMemberId, targetProjectId int) error {
	if invitorMemberId <= 0 && inviteeMemberId <= 0 {
		return errors.New("invitor_member_id 和 invitee_member_id 至少需要一个大于 0 以作为筛选条件")
	}
	if invitorMemberId > 0 && inviteeMemberId > 0 {
		return errors.New("invitor_member_id 和 invitee_member_id 不能同时大于 0 作为筛选条件")
	}
	if targetProjectId <= 0 {
		return errors.New("target_project_id 必须大于 0")
	}

	return nil
}
