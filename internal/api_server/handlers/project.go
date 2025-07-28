package handlers

import (
	"fmt"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Status 是在 handler 层为了方便使用而定义的项目状态类型
// @Description 项目状态使用位掩码表示，与 LaborMask 类似
// @Description 统一每两位表示一个状态，0b00 表示未开始，0b01 表示进行中，0b10 表示已完成，0b11 表示在搜索中忽略
// @Description 第 0~1 位表示翻译状态
// @Description 第 2~3 位表示校对状态
// @Description 第 4~5 位表示嵌字状态
// @Description 第 6~7 位表示嵌字审核状态
// @Description 第 8~9 位表示发布状态
type Status = services.ProjectStatus

// CreateProjectRequest 定义了创建项目的请求体
type CreateProjectRequest struct {
	// 申请者的成员 ID
	ApplicantMemberId uint `json:"applicant_member_id"`

	// 项目标题，必填
	Title string `json:"title" binding:"required"`
	// 作品集 ID，必填
	WorksetId uint `json:"workset_id" binding:"required"`
	// 项目描述
	Description string `json:"description"`

	// 是否允许自动加入，默认为 false
	AllowAutoJoin bool `json:"allow_auto_join"`
	// 是否隐藏项目，默认为 false
	IsHidden bool `json:"is_hidden"`
}

// UpdateProjectRequest 定义了更新项目的请求体
type UpdateProjectRequest struct {
	// 操作成员的 ID，即进行更新这个操作的成员 ID
	OperatorMemberId uint `json:"operator_member_id"`

	// 项目标题
	Title string `json:"title"`
	// 项目描述
	Description string `json:"description"`
	// 项目状态，使用位掩码表示，不使用的未应当使用 0b11 填充
	// @example 0b11111011
	Status Status `json:"status"`
}

// DeleteProjectRequest 定义了删除项目的请求体
type DeleteProjectRequest struct {
	// 操作成员的 ID，即进行删除这个操作的成员 ID
	OperatorMemberId uint `json:"operator_member_id"`
}

// LaborInfo 定义了成员的分工信息
type LaborInfo struct {
	MemberId  uint      `json:"member_id"`  // 成员 ID
	Nickname  string    `json:"nickname"`   // 成员昵称
	LaborMask LaborMask `json:"labor_mask"` // 分工掩码，使用
}

// ProjectInfo 定义了项目的基本信息
type ProjectInfo struct {
	// 项目 ID
	Id uint `json:"id"`
	// 项目标题
	Title string `json:"title"`
	// 项目描述，如果为空将不会返回
	Description string `json:"description,omitempty"`
	// 项目状态，使用位掩码表示
	Status Status `json:"status"`
	// 龙译 ID
	MoetranId string `json:"moetran_id"`
	// 作品集 ID
	WorksetId uint `json:"workset_id"`
	// 作品集索引
	WorksetIndex uint `json:"workset_index"`
	// Legacy ID，如果未指定将不会返回
	LegacyId uint `json:"legacy_id,omitempty"`
	// 是否允许自动加入，如果未指定将不会返回
	AllowAutoJoin bool `json:"allow_auto_join,omitempty"`
	// 相关的分工，如果未指定将不会返回
	Labors []*LaborInfo `json:"labors,omitempty"`
}

// RouteProjectHandler 注册项目相关的路由
func RouteProjectHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/projects").
		Handle(new(ProjectHandler))
}

// ProjectHandler 处理项目相关的请求
type ProjectHandler struct {
	ProjectService services.ProjectService
}

// BeforeActivation 在控制器激活前注册路由
func (p *ProjectHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
	b.Handle("GET", "/{id:uint}", "ProjectDetail")
	b.Handle("POST", "/", "Create")
	b.Handle("DELETE", "/{id:uint}", "Delete")
	b.Handle("PATCH", "/{id:uint}", "Update")
}

// List godoc
// @Summary     获取项目列表分页
// @Description 根据作品集 ID 获取项目列表，支持分页、排序和状态筛选。当列表为空时，会返回 null 而不是空数组。
//
// @Param       page_serial query integer false "页码，默认值为 1"
// @Param       page_size query integer false "每页数量，默认值为 10"
// @Param       status_mask query integer false "项目状态（位掩码），用于复合查询，可以有多个 status，默认不筛选"
// @Param       member_id query integer false "成员 ID，默认不筛选成员"
// @Param       workset_id query integer false "项目所属的作品集 ID，默认不筛选"
// @Param       index query integer false "项目的 workset index 或 legacy ID，默认不筛选"
//
// @Tags        project
// @Produce     json
// @Success     200 {object} []ProjectInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router      /api/projects [get]
func (h *ProjectHandler) List(ctx iris.Context) {
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	statusMaskSlice := ctx.URLParamSlice("status_mask")
	statusMasks, err := h.sliceStringToStatus(statusMaskSlice)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "无法转换状态参数",
			Detail: err.Error(),
		})
		return
	}

	// 注意默认为 0，代表不筛选 workset
	worksetId := ctx.URLParamIntDefault("workset_id", 0)

	// 注意默认为 0，代表不筛选 user
	memberId := ctx.URLParamIntDefault("member_id", 0)

	// 注意默认为 0，代表不筛选 index / legacy id
	index := ctx.URLParamIntDefault("index", 0)

	// 根据参数调用服务层获取数据
	projects, err := h.ProjectService.GetProjects(&services.ProjectListParams{
		Offset:      (pageSerial - 1) * pageSize,
		Limit:       pageSize,
		StatusMasks: statusMasks,
		MemberId:    uint(memberId),
		WorksetId:   uint(worksetId),
		Index:       uint(index),
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取项目列表失败",
			Detail: err.Error(),
		})
		return
	}

	// 将 service 的格式转换为 DTO 格式
	projectInfos := make([]*ProjectInfo, 0, len(projects))

	for _, project := range projects {
		projectInfos = append(projectInfos, serviceProjectToDTO(project))
	}

	ctx.JSON(projectInfos)
}

// Detail godoc
// @Summary 	获取项目详情
// @Description 获取指定项目的详细信息，包括各类状态和分工信息
//
// @Param 		id path integer true "项目 ID"
//
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} ProjectInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/projects/{id} [get]
func (h *ProjectHandler) Detail(ctx iris.Context) {
	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 调用服务层获取数据
	project, err := h.ProjectService.GetProjectDetail(projectId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取单个项目详情失败",
			Detail: err.Error(),
		})
		return
	}

	// 将 service 的格式转换为 DTO 格式
	ctx.JSON(serviceProjectToDTO(project))
}

// Create godoc
// @Summary 	创建项目
// @Description 创建一个新的项目
//
// @Accept      json
// @Param       body_params body CreateProjectRequest true "创建项目的请求体"
//
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "创建成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/projects [post]
func (h *ProjectHandler) Create(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	// 将请求体绑定到 CreateProjectRequest 结构体
	var request CreateProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "请求体格式错误",
		})
		return
	}

	// 调用服务层创建项目
	info, err := h.ProjectService.CreateProject(&services.CreateProjectParams{
		CreatorMemberId:  request.ApplicantMemberId,
		CurrentMemberIds: memberIds,

		Title:         request.Title,
		Description:   request.Description,
		WorksetId:     request.WorksetId,
		AllowAutoJoin: request.AllowAutoJoin,
		IsHidden:      request.IsHidden,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "创建项目失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: info.Message,
		Detail: &iris.Map{
			"project_id": info.ProjectId,
			"moetran_id": info.MoetranId,
		},
	})
}

// Delete godoc
// @Summary 	删除项目
// @Description 删除指定的项目，需提供项目 ID
//
// @Param       id path integer true "项目 ID"
//
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "删除成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/projects/{id} [delete]
func (h *ProjectHandler) Delete(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 获取请求体参数
	var request DeleteProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "请求体格式错误",
		})
		return
	}

	// 调用服务层删除项目
	if err := h.ProjectService.DeleteProject(&services.DeleteProjectParams{
		UsingMemberId:    request.OperatorMemberId,
		CurrentMemberIds: memberIds,

		ProjectId: projectId,
	}); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "删除项目失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "项目删除成功",
	})
}

// Update godoc
// @Summary 	更新项目
// @Description 更新指定的项目，需提供项目 ID
//
// @Accept      json
// @Param       id path integer true "项目 ID"
// @Param       body_params body UpdateProjectRequest true "更新的项目信息"
//
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "更新成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/projects/{id} [patch]
func (h *ProjectHandler) Update(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 member_ids",
		})
		return
	}

	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 读取请求体
	var request UpdateProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "请求体格式错误",
		})
		return
	}

	// 调用服务层更新项目
	if err := h.ProjectService.UpdateProject(&services.UpdateProjectParams{
		UsingMemberId:    request.OperatorMemberId,
		CurrentMemberIds: memberIds,

		ProjectId:   projectId,
		Title:       request.Title,
		Description: request.Description,
	}); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "更新项目失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(SuccessResponse{
		Message: "项目更新成功",
	})
}

// sliceStringToStatus 将字符串切片转换为无符号整数切片
func (*ProjectHandler) sliceStringToStatus(slice []string) ([]services.ProjectStatus, error) {
	result := make([]services.ProjectStatus, 0, len(slice))
	for _, str := range slice {
		var num uint
		if _, err := fmt.Sscanf(str, "%d", &num); err != nil {
			return nil, fmt.Errorf("无法转换字符串 '%s' 为无符号整数: %w", str, err)
		}
		result = append(result, services.ProjectStatus(num))
	}
	return result, nil
}

// serviceProjectToDTO 将服务层的 ProjectInfo 转换为 DTO 格式
func serviceProjectToDTO(project *services.ProjectInfo) *ProjectInfo {
	if project == nil {
		return nil
	}

	projectInfo := &ProjectInfo{
		Id:            project.Id,
		Title:         project.Title,
		Description:   project.Description,
		Status:        project.Status,
		MoetranId:     project.MoetranId,
		WorksetId:     project.WorksetId,
		WorksetIndex:  project.WorksetIndex,
		LegacyId:      project.LegacyId,
		AllowAutoJoin: project.AllowAutoJoin,
	}

	if project.Labors != nil {
		projectInfo.Labors = make([]*LaborInfo, 0, len(project.Labors))
		// 将分工信息转换为 DTO 格式
		for _, labor := range project.Labors {
			projectInfo.Labors = append(projectInfo.Labors, &LaborInfo{
				MemberId:  labor.MemberId,
				Nickname:  labor.Nickname,
				LaborMask: labor.LaborMask,
			})
		}
	}

	return projectInfo
}
