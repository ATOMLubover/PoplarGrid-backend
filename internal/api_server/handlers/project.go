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

	// 尨译项目的信息，直接以 string 形式返回，后端不做解析
	MoetranProjectInfo string `json:"moetran_project_info,omitempty"`
}

// ProjectDetail 定义了项目的响应格式
type ProjectDetail struct {
	ProjectInfo *ProjectInfo `json:"project_info"`
	// 相关的分工，如果未指定将不会返回
	Labors []*LaborInfo `json:"labors,omitempty"`
}

// ProjectCreatedInfo 定义了创建项目后的返回信息
type ProjectCreatedInfo struct {
	ProjectId uint   // 项目 ID
	MoetranId string // 龙译 ID
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
// @Success     200 {object} FormatResponse[[]ProjectInfo]
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router      /api/projects [get]
func (h *ProjectHandler) List(ctx iris.Context) {
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	statusMaskSlice := ctx.URLParamSlice("status_mask")
	statusMasks, err := h.sliceStringToStatus(statusMaskSlice)
	if err != nil {
		wrapError(ctx, newHdlErr(ErrBadParams, fmt.Sprintf("无法转换 status_mask: %v", err)))
		return
	}

	// 注意默认为 0，代表不筛选 workset
	worksetId := ctx.URLParamIntDefault("workset_id", 0)

	// 注意默认为 0，代表不筛选 user
	memberId := ctx.URLParamIntDefault("member_id", 0)

	// 注意默认为 0，代表不筛选 index / legacy id
	index := ctx.URLParamIntDefault("index", 0)

	// 根据参数调用服务层获取数据
	projects, e := h.ProjectService.GetProjects(&services.ProjectListParams{
		Offset:      (pageSerial - 1) * pageSize,
		Limit:       pageSize,
		StatusMasks: statusMasks,
		MemberId:    uint(memberId),
		WorksetId:   uint(worksetId),
		Index:       uint(index),
	})
	if e != nil {
		wrapError(ctx, e)
		return
	}

	// 将 service 的格式转换为 DTO 格式
	projectInfos := make([]*ProjectInfo, 0, len(projects))

	for _, project := range projects {
		projectInfos = append(projectInfos, serviceProjectToDTO(project))
	}

	wrapSuccess(ctx, projectInfos)
}

// Detail godoc
// @Summary 	获取项目详情
// @Description 获取指定项目的详细信息，包括各类状态和分工信息
//
// @Param 		id path integer true "项目 ID"
//
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} FormatResponse[ProjectInfo]
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/projects/{id} [get]
func (h *ProjectHandler) Detail(ctx iris.Context) {
	// 读取上下文中的 user_id
	userId, err := ctx.Values().GetUint("user_id")
	if err != nil || userId <= 0 {
		wrapError(ctx, newHdlErr(ErrHeaderLackage, "未提取到有效的 user ID"))
		return
	}

	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无法获取有效的 project ID"))
		return
	}

	// 调用服务层获取数据
	projectInfo, e := h.ProjectService.GetProjectDetail(&services.ProjectDetailParams{
		ProjectId: projectId,
		UserId:    userId,
	})
	if e != nil {
		wrapError(ctx, e)
		return
	}

	// 将 service 的格式转换为 DTO 格式
	res := &ProjectDetail{
		ProjectInfo: serviceProjectToDTO(projectInfo),
	}

	// 转换分工信息
	if projectInfo.Labors != nil {
		res.Labors = make([]*LaborInfo, len(projectInfo.Labors))

		// 将分工信息转换为 DTO 格式
		for i, labor := range projectInfo.Labors {
			res.Labors[i] = &LaborInfo{
				MemberId:  labor.MemberId,
				Nickname:  labor.Nickname,
				LaborMask: labor.LaborMask,
			}
		}
	}

	wrapSuccess(ctx, res)
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
// @Success	 	200 {object} FormatResponse[ProjectCreatedInfo] "创建成功"
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/projects [post]
func (h *ProjectHandler) Create(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		wrapError(ctx, newHdlErr(ErrHeaderLackage, "未提取到有效 member IDs"))
		return
	}

	// 将请求体绑定到 CreateProjectRequest 结构体
	var request CreateProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无效的请求参数"))
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
		wrapError(ctx, err)
		return
	}

	wrapSuccess(ctx, &ProjectCreatedInfo{
		ProjectId: info.ProjectId,
		MoetranId: info.MoetranId,
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
// @Success	 	200 {object} StringFormatResponse "删除成功"
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/projects/{id} [delete]
func (h *ProjectHandler) Delete(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		wrapError(ctx, newHdlErr(ErrHeaderLackage, "未提取到有效的 member IDs"))
		return
	}

	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无法获取有效的 project ID"))
		return
	}

	// 获取请求体参数
	var request DeleteProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无效的请求参数"))
		return
	}

	// 调用服务层删除项目
	if err := h.ProjectService.DeleteProject(&services.DeleteProjectParams{
		UsingMemberId:    request.OperatorMemberId,
		CurrentMemberIds: memberIds,

		ProjectId: projectId,
	}); err != nil {
		wrapError(ctx, err)
		return
	}

	wrapSuccess(ctx, "项目删除成功")
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
// @Success	 	200 {object} StringFormatResponse "更新成功"
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/projects/{id} [patch]
func (h *ProjectHandler) Update(ctx iris.Context) {
	// 读取上下文中的 member_ids
	memberIds, ok := ctx.Values().Get("member_ids").(map[uint]struct{})
	if !ok {
		wrapError(ctx, newHdlErr(ErrHeaderLackage, "未提取到有效的 member IDs"))
		return
	}

	// 从 path 参数中获取项目 ID
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无法获取有效的 project ID"))
		return
	}

	// 读取请求体
	var request UpdateProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无效的请求参数"))
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
		wrapError(ctx, err)
		return
	}

	wrapSuccess(ctx, "项目更新成功")
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

	return projectInfo
}
