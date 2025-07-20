package handlers

import (
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteProjectHandler 注册项目相关的路由
func RouteProjectHandler(root *mvc.Application) {
	// 创建 handler 实例
	projectHandler := &ProjectHandler{}
	projectProcHandler := &ProjectProcHandler{}
	projectLaborHandler := &ProjectLaborHandler{}

	// 注册路由组和 handler
	projectParty := root.
		Party("/projects").
		Handle(projectHandler).
		Handle(projectProcHandler).
		Handle(projectLaborHandler)

	// // 注册 handler
	// projectParty.Handle(projectHandler)
	// projectParty.Handle(projectProcHandler)
	// projectParty.Handle(projectLaborHandler)

	// 注册路由与方法间的映射（静态安全）
	projectParty.Router.Get("", projectHandler.ProjectListPage)
	projectParty.Router.Get("/{id:uint}", projectHandler.ProjectDetail)

	projectParty.Router.Post("", projectProcHandler.Create)

	projectParty.Router.Delete("/{id:uint}", projectProcHandler.Delete)

	projectParty.Router.Patch("/{id:uint}", projectProcHandler.UpdateInfo)

	projectParty.Router.Get("/{id:uint}/labors", projectLaborHandler.LaborDivision)
}

// ProjectHandler 处理项目相关的请求
type ProjectHandler struct {
	ProjectService services.ProjectService
}

// ProjectListPage godoc
// @Summary     获取项目列表分页 (按作品集筛选)
// @Description 根据作品集 ID 获取项目列表，支持分页、排序和状态筛选。当列表为空时，会返回 null 而不是空数组。
// @Param       page_serial query integer false "页码，默认值为 1"
// @Param       page_size query integer false "每页数量，默认值为 10"
// @Param       sort query integer false "排序方式，0：按 ID 倒序，1：按 updated_at 倒序"
// @Param       status query integer false "项目状态（位掩码），用于复合查询，默认不筛选查询"
// @Param       workset_id query integer true "项目所属的作品集 ID"
// @Tags        project
// @Produce     json
// @Success     200 {object} []dtos.ProjectBasic
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router      /projects [get]
func (h *ProjectHandler) ProjectListPage(ctx iris.Context) {
	pageSerial := ctx.URLParamInt32Default("page_serial", 1)

	pageSize := ctx.URLParamInt32Default("page_size", 10)

	sort := ctx.URLParamInt32Default("sort", 0)

	// 注意默认为 PROJECT_STATUS_ALL，表示不筛选状态
	status := ctx.URLParamInt32Default("status", dtos.PROJECT_STATUS_ALL)

	worksetId := ctx.URLParamIntDefault("workset_id", 0)

	// 如果提供了 workset_id，则查询该作品集下的项目列表
	projects, err := h.ProjectService.GetBasicPageWithParams(uint(worksetId), int(pageSerial), int(pageSize), int(sort), dtos.ProjectOverallStatus(status))
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

// // ProjectListPageByUserId godoc
// // @Summary 	获取用户参与的项目列表分页，默认按照 ID 倒序排列
// // @Description 注意当列表为空，会返回 null 而不是空数组
// // @Param 		page_serial query integer false "页码，默认值为 1"
// // @Param 		page_size query integer false "每页数量，默认值为 10"
// // @Param 		user_id query integer true "用户 ID"
// // @Tags 		project
// // @Produce 	json
// // @Success	 	200 {object} []dtos.MyProjectBasic
// // @Failure     400 {object} ErrorResponse "无效的请求参数"
// // @Failure     500 {object} ErrorResponse "服务器内部错误"
// // @Router 		/projects/list_by_user [get]
// func (h *ProjectHandler) ProjectListPageByUserId(ctx iris.Context) {
// 	pageSerial := ctx.URLParamInt32Default("page_serial", 1)

// 	pageSize := ctx.URLParamInt32Default("page_size", 10)

// 	userId, err := ctx.URLParamInt("user_id")
// 	if err != nil || userId <= 0 {
// 		ctx.StatusCode(iris.StatusBadRequest)
// 		ctx.JSON(ErrorResponse{
// 			Error: "user_id 必须是一个明确给出的正整数",
// 		})
// 		return
// 	}

// 	// 调用服务层获取数据
// 	projects, err := h.ProjectService.GetBasicPageByUserId(uint(userId), int(pageSerial), int(pageSize))
// 	if err != nil {
// 		ctx.StatusCode(iris.StatusInternalServerError)
// 		ctx.JSON(ErrorResponse{
// 			Error:  "获取用户参与的项目列表失败",
// 			Detail: err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(projects)
// }

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
	projectId, err := ctx.URLParamInt("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "id 必须是一个明确给出的正整数"})
		return
	}

	// 调用服务层获取数据
	project, err := h.ProjectService.GetDetailById(uint(projectId))
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

// Create godoc
// @Summary 	创建项目
// @Description 创建一个新的项目
// @Accept      application/json
// @Param       body_params body dtos.CreateProjectRequest true "创建项目的请求体"
// @Tags 		project_proc
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
	info, err := h.ProjectService.CreateProject(&dtos.CreateProjectInfo{
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
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "删除成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id} [delete]
func (h *ProjectProcHandler) Delete(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 调用服务层删除项目
	if err := h.ProjectService.DeleteProjectById(uint(projectId)); err != nil {
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
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} SuccessResponse "更新成功"
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id} [patch]
func (h *ProjectProcHandler) UpdateInfo(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("id")
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
	if err := h.ProjectService.UpdateProjectInfo(uint(projectId), &request); err != nil {
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

// ProjectLaborHandler 处理项目角色相关的请求
type ProjectLaborHandler struct {
	ProjectService services.ProjectService
}

// LaborDivision godoc
// @Summary 	获取指定项目的分工信息
// @Description 获取指定项目的分工信息，包括成员的角色和状态
// @Param 		id path uint true "项目 ID"
// @Tags 		project_labor
// @Produce 	json
// @Success	 	200 {object} []dtos.LaborDivision
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
// @Router 		/projects/{id}/labors [get]
func (h *ProjectLaborHandler) LaborDivision(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("project_id")
	if err != nil || projectId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获得有效的 project_id",
		})
		return
	}

	// 调用服务层获取分工信息
	laborDivisions, err := h.ProjectService.GetLaborDivisionByProjectId(uint(projectId))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取项目分工信息失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(laborDivisions)
}
