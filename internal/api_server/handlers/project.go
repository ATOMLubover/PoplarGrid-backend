package handlers

import (
	"fmt"
	"poplargrid/internal/api_server/dtos"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

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
	// 注册 GET /projects
	b.Handle("GET", "/", "ProjectListPage")
	// 注册 GET /projects/{id:uint}
	b.Handle("GET", "/{id:uint}", "ProjectDetail")
}

// ProjectListPage godoc
// @Summary     获取项目列表分页
// @Description 根据作品集 ID 获取项目列表，支持分页、排序和状态筛选。当列表为空时，会返回 null 而不是空数组。
// @Param       page_serial query integer false "页码，默认值为 1"
// @Param       page_size query integer false "每页数量，默认值为 10"
// @Param       sort query integer false "排序方式，0：按 ID 倒序，1：按 updated_at 倒序，默认按 ID 倒序"
// @Param       status query integer false "项目状态（位掩码），用于复合查询，默认不筛选查询"
// @Param       member_id query integer false "成员 ID，默认不筛选成员"
// @Param       workset_id query integer false "项目所属的作品集 ID，默认不筛选作品集"
// @Param       index query integer false "项目的索引或 legacy ID，默认不筛选"
// @Tags        project
// @Produce     json
// @Success     200 {object} []dtos.ProjectBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /projects [get]
func (h *ProjectHandler) ProjectListPage(ctx iris.Context) {
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	statusesSlice := ctx.URLParamSlice("status")
	statuses, err := h.sliceStringToStatus(statusesSlice)
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
		Offset:    (pageSerial - 1) * pageSize,
		Limit:     pageSize,
		Status:    statuses,
		MemberId:  uint(memberId),
		WorksetId: uint(worksetId),
		Index:     uint(index),
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取项目列表失败",
			Detail: err.Error(),
		})
		return
	}
	ctx.JSON(projects)

}

// ProjectDetail godoc
// @Summary 	获取项目详情
// @Description 获取指定项目的详细信息，包括翻译、校对、嵌字、审核、发布等状态
// @Param 		id path integer true "项目 ID"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectDetail
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id} [get]
func (h *ProjectHandler) ProjectDetail(ctx iris.Context) {
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

	ctx.JSON(project)
}

// ProjectProcHandler 处理与项目进度动态推进有关的路由
type ProjectProcHandler struct {
	ProjectService services.ProjectService
}

// BeforeActivation 在控制器激活前注册路由
func (p *ProjectProcHandler) BeforeActivation(b mvc.BeforeActivation) {
	// 注册 POST /projects
	b.Handle("POST", "/", "Create")
	// 注册 DELETE /projects/{id:uint}
	b.Handle("DELETE", "/{id:uint}", "Delete")
	// 注册 PATCH /projects/{id:uint}
	b.Handle("PATCH", "/{id:uint}", "UpdateInfo")
}

// Create godoc
// @Summary 	创建项目
// @Description 创建一个新的项目
// @Accept      application/json
// @Param       body_params body dtos.CreateProjectRequest true "创建项目的请求体"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectCreatedInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects [post]
func (h *ProjectProcHandler) Create(ctx iris.Context) {
	// 读取上下文中的 user_id
	userId, err := ctx.Values().GetUint("user_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "未提取到有效 user_id",
		})
		return
	}

	// 将请求体绑定到 CreateProjectRequest 结构体
	var request dtos.CreateProjectRequest

	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "请求体格式错误"})
		return
	}

	// 调用服务层创建项目
	info, err := h.ProjectService.CreateProject(&dtos.CreateProjectParams{
		Title:         request.Title,
		Description:   request.Description,
		WorksetId:     request.WorksetId,
		CreatorUserId: userId,
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

	ctx.JSON(info)
}

// Delete godoc
// @Summary 	删除项目
// @Description 删除指定的项目，需提供项目 ID
// @Param       id path integer true "项目 ID"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "删除成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id} [delete]
func (h *ProjectProcHandler) Delete(ctx iris.Context) {
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 调用服务层删除项目
	if err := h.ProjectService.DeleteProject(projectId); err != nil {
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
// @Accept      application/json
// @Param       id path integer true "项目 ID"
// @Param       body_params body dtos.UpdateProjectRequest true "更新的项目信息"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "更新成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id} [patch]
func (h *ProjectProcHandler) UpdateInfo(ctx iris.Context) {
	projectId, err := ctx.Params().GetUint("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	var request dtos.UpdateProjectRequest
	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "请求体格式错误",
		})
		return
	}

	// 调用服务层更新项目
	if err := h.ProjectService.UpdateProject(projectId, &request); err != nil {
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
