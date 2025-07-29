package handlers

import (
	"errors"
	"poplargrid/internal/api_server/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// ApplicationInfo 定义了申请的基本信息 DTO
type ApplicationInfo struct {
	// 申请 ID
	Id uint `json:"id"`
	// 申请时间
	Time string `json:"time"`
	// 申请者成员 ID
	ApplicantMemberId uint `json:"applicant_member_id"`
	// 申请者成员昵称
	ApplicantNickname string `json:"applicant_nickname"`
	// 处理者成员 ID
	ProcessorMemberId uint `json:"processor_member_id"`
	// 处理者成员昵称
	ProcessorNickname string `json:"processor_nickname"`
	// 目标项目 ID
	TargetProjectId uint `json:"target_project_id"`
	// 目标项目名称
	TargetProjectName string `json:"target_project_name"`
	// 申请分工掩码
	TargetLaborMask LaborMask `json:"target_labor_mask"`
}

// CreateApplicationRequest 定义了创建申请的请求参数
type CreateApplicationRequest struct {
	// 申请者成员 ID
	ApplicantMemberId uint `json:"applicant_member_id" binding:"required"`
	// 目标项目 ID
	TargetProjectId uint `json:"target_project_id" binding:"required"`
	// 申请分工掩码
	TargetLaborMask services.LaborMask `json:"target_labor_mask" binding:"required"`
}

// ProcessApplicationRequest 定义了处理申请的请求参数
type ProcessApplicationRequest struct {
	// 申请 ID
	ApplicationId uint `json:"application_id" binding:"required"`
	// 处理者成员 ID
	ProcessorMemberId uint `json:"processor_member_id" binding:"required"`
	// 是否接受申请
	Accept bool `json:"accept" binding:"required"`
}

// RouteApplicationHandler 注册申请相关的路由
func RouteApplicationHandler(root *mvc.Application) {
	root.Party("/applications").
		Handle(new(ApplicationHandler))
}

// ApplicationHandler 处理申请相关的请求
type ApplicationHandler struct {
	applicationService services.ApplicationService
}

// BeforeActivation 在控制器激活前注册路由
func (h *ApplicationHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
	b.Handle("POST", "/", "Create")
	b.Handle("PUT", "/{id:uint}", "Process")
}

// List godoc
// @Summary 	获取申请列表
// @Description 获取当前用户的申请列表，根据参数的选择来确定到底是返回申请者的申请还是处理者的申请
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param       applicant_member_id query int false "申请人的成员 ID，与处理者的成员 ID 只能二选一"
// @Param       processor_member_id query int false "处理者的成员 ID，与申请人的成员 ID 只能二选一"
// @Param       target_project_id query int false "目标项目 ID，如果需要筛选特定项目的申请，可以使用此参数"
//
// @Tags 		application
// @Produce 	json
// @Success 	200 {object} []ApplicationInfo
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/applications [get]
func (h *ApplicationHandler) List(ctx iris.Context) {
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
	applicantMemberId := ctx.URLParamIntDefault("applicant_member_id", 0)
	processorMemberId := ctx.URLParamIntDefault("processor_member_id", 0)
	targetProjectId := ctx.URLParamIntDefault("target_project_id", 0)

	// 检查参数的组合是否合法
	if err := checkApplicationListQueryParams(applicantMemberId, processorMemberId, targetProjectId); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "查询参数组合非法",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务获取申请列表
	applications, err := h.applicationService.GetApplications(&services.ApplicationListParams{
		Offset:            (pageSerial - 1) * pageSize,
		Limit:             pageSize,
		ApplicantMemberId: uint(applicantMemberId),
		ProcessorMemberId: uint(processorMemberId),
		TargetProjectId:   uint(targetProjectId),
		CurrentMemberIds:  memberIds,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取申请列表失败",
			Detail: err.Error(),
		})
		return
	}

	// 将申请列表转换为 DTO
	applicationInfos := make([]ApplicationInfo, 0, len(applications))

	for _, app := range applications {
		applicationInfos = append(applicationInfos, ApplicationInfo{
			Id:                app.Id,
			Time:              app.Time.Format(time.DateTime),
			ApplicantMemberId: app.ApplicantMemberId,
			ApplicantNickname: app.ApplicantNickname,
			ProcessorMemberId: app.ProcessorMemberId,
			ProcessorNickname: app.ProcessorNickname,
			TargetProjectId:   app.TargetProjectId,
			TargetProjectName: app.TargetProjectTitle,
			TargetLaborMask:   app.TargetLaborMask,
		})
	}

	ctx.JSON(applicationInfos)
}

// Create godoc
// @Summary 	创建申请
// @Description 创建一个新的申请记录
//
// @Accept 		json
// @Param 		body_params body CreateApplicationRequest true "申请的基本信息"
//
// @Tags 		application
// @Produce 	json
// @Success 	200 {object} SuccessResponse "创建成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/applications [post]
func (h *ApplicationHandler) Create(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	var req CreateApplicationRequest

	// 绑定请求体参数
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "无效的请求参数",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务创建申请
	if err := h.applicationService.CreateApplication(&services.CreateApplicationParams{
		ApplicantMemberId: req.ApplicantMemberId,
		TargetProjectId:   req.TargetProjectId,
		TargetLaborMask:   req.TargetLaborMask,
		CurrentMemberIds:  memberIds,
	}); err != nil {
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

// Process godoc
// @Summary 	处理申请
// @Description 接受或拒绝申请
//
// @Accept 		json
// @Param 		id path uint true "申请 ID"
// @Param 		body_params body ProcessApplicationRequest true "处理申请的请求参数"
//
// @Tags 		application
// @Produce 	json
// @Success 	200 {object} SuccessResponse "处理成功"
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/applications/{id} [put]
func (h *ApplicationHandler) Process(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	var req ProcessApplicationRequest

	// 绑定请求体参数
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "无效的请求参数",
			Detail: err.Error(),
		})
		return
	}

	// 调用服务处理申请
	if err := h.applicationService.ProcessApplication(&services.ProcessApplicationParams{
		ApplicationId:     req.ApplicationId,
		ProcessorMemberId: req.ProcessorMemberId,
		Accept:            req.Accept,
		CurrentMemberIds:  memberIds,
	}); err != nil {
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

// checkApplicationListQueryParams 检查申请列表查询参数的合法性
func checkApplicationListQueryParams(applicantMemberId, processorMemberId, targetProjectId int) error {
	if applicantMemberId <= 0 && processorMemberId <= 0 {
		return errors.New("applicant_member_id 和 processor_member_id 至少需要一个大于 0 以作为筛选条件")
	}
	if applicantMemberId > 0 && processorMemberId > 0 {
		return errors.New("applicant_member_id 和 processor_member_id 不能同时大于 0 作为筛选条件")
	}
	if targetProjectId <= 0 {
		return errors.New("target_project_id 必须大于 0")
	}

	return nil
}
