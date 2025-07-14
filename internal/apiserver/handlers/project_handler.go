package handlers

import (
	"poplargrid/internal/apiserver/dtos"
	"poplargrid/internal/apiserver/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteProjectHandler 注册项目相关的路由
func RouteProjectHandler(root *mvc.Application) {
	// 创建 ProjectHandler 实例
	projectHandler := &ProjectHandler{}

	// 注册路由组
	projectParty := root.Party("/project")

	// 注册 handler
	projectParty.Handle(projectHandler)

	// 注册路由与方法的映射（类型安全）
	projectParty.Router.Get("/stats", projectHandler.ProjectStats)
	projectParty.Router.Get("/list", projectHandler.ProjectListPage)
	projectParty.Router.Get("/published_list", projectHandler.ProjectListPublishedPage)
	projectParty.Router.Get("/detail", projectHandler.ProjectDetail)

	{ // 创建 ProjectProcHandler 实例
		procHandler := &ProjectProcHandler{}

		// 注册路由组
		procParty := projectParty.Party("/proc")

		// 注册 handler
		procParty.Handle(procHandler)

		// 注册路由与方法的映射（类型安全）
		procParty.Router.Post("/create", procHandler.Create)
		procParty.Router.Delete("/delete", procHandler.Delete)
		procParty.Router.Put("/update_info", procHandler.UpdateInfo)
		procParty.Router.Put("/update_status", procHandler.UpdateStatus)
	}

	{
		// 创建 ProjectLaborHandler 实例
		laborHandler := &ProjectLaborHandler{}

		// 注册路由组
		laborParty := projectParty.Party("/role")

		// 注册 handler
		laborParty.Handle(laborHandler)

		// 注册路由与方法的映射（类型安全）
		laborParty.Router.Post("/invite", laborHandler.Invite)
		laborParty.Router.Post("/accept_appli", laborHandler.AcceptAppli)
		laborParty.Router.Post("/refuse_appli", laborHandler.RefuseAppli)

		laborParty.Router.Post("/apply", laborHandler.Apply)
		laborParty.Router.Post("/accept", laborHandler.Accept)
		laborParty.Router.Post("/refuse", laborHandler.Refuse)
	}
}

// ProjectHandler 处理项目相关的请求
type ProjectHandler struct {
	ProjectService services.ProjectService
}

// ProjectStats godoc
// @Summary 	获取项目统计信息
// @Description 获取所有项目的统计信息，包括总数、翻译进行/完成、校对进行/完成、嵌字进行/完成，审核进行/完成、发布完成对应数量等
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectStats
// @Router 		/project/stats [get]
func (h *ProjectHandler) ProjectStats(ctx iris.Context) {

}

// ProjectListPage godoc
// @Summary     获取项目列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组；如果要单独查询已发布的项目列表，请使用 /project/published_list 接口
// @Param       page_serial query integer false "页码，默认值为 1" default(1)
// @Param       page_size query integer false "每页数量，默认值为 20" default(10)
// @Param       sort query integer false "排序方式，0：按 ID 倒序，1：按 updated_at 倒序"
// @Param       status query integer false "项目状态（位掩码），用于复合查询，默认不筛选查询"
// @Param       workset_id query integer true "项目所属的作品集 ID，必填"
// @Tags        project
// @Produce     json
// @Success     200 {object} []dtos.ProjectBasic
// @Router      /project/list [get]
func (h *ProjectHandler) ProjectListPage(ctx iris.Context) {
	pageSerial := ctx.URLParamInt32Default("page_serial", 1)

	pageSize := ctx.URLParamInt32Default("page_size", 10)

	sort := ctx.URLParamInt32Default("sort", 0)

	// 注意默认为 -1，表示不筛选状态，需要调用特殊的 service 函数处理
	status := ctx.URLParamInt32Default("status", -1)

	worksetId, err := ctx.URLParamInt("workset_id")
	if err != nil || worksetId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "workset_id 必须是一个明确给出的正整数"})
		return
	}

	// 调用服务层获取数据
	projects, err := h.ProjectService.GetBasicPageWithParams(uint(worksetId), int(pageSerial), int(pageSize), int(sort), dtos.ProjectOverallStatus(status))
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取项目列表失败"})
		return
	}

	ctx.JSON(projects)
}

// ProjectListPublished godoc
// @Summary 	获取已发布的项目列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query integer true "列表查询参数"
// @Param 		page_size query integer true "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_asc | id_desc | poplarity_desc，其他输入无效"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} []dtos.ProjectBasic
// @Router 		/project/published_list [get]
func (h *ProjectHandler) ProjectListPublishedPage(ctx iris.Context) {
}

// ProjectDetail godoc
// @Summary 	获取项目详情
// @Description 获取指定项目的详细信息，包括翻译、校对、嵌字、审核、发布等状态
// @Param 		id query string true "项目 ID"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dtos.ProjectDetail
// @Router 		/project/{id} [get]
func (h *ProjectHandler) ProjectDetail(ctx iris.Context) {

}

// ProjectProcHandler 处理与项目进度动态推进有关的路由
type ProjectProcHandler struct {
}

// Create godoc
// @Summary 	创建项目
// @Description 创建一个新的项目，请求体暂时未确定
// @Accept      application/json
// @Param       request body dtos.CreateProjectRequest true "创建项目的请求体"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]any
// @Router 		/project/proc/create [post]
func (h *ProjectProcHandler) Create(ctx iris.Context) {

}

// Delete godoc
// @Summary 	删除项目
// @Description 删除指定的项目，需提供项目 ID
// @Accept      multipart/form-data
// @Param       body_params body dtos.DeleteProjectRequest true "要删除项目的信息"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/delete [delete]
func (h *ProjectProcHandler) Delete(ctx iris.Context) {
}

// Update godoc
// @Summary 	更新项目
// @Description 更新指定的项目，需提供项目 ID
// @Accept      multipart/form-data
// @Param       body_params body dtos.UpdateProjectRequest true "更新的项目信息"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/update_info [put]
func (h *ProjectProcHandler) UpdateInfo(ctx iris.Context) {
}

// UpdateStatus godoc
// @Summary 	更新项目状态
// @Description 更新指定项目的状态，需提供项目 ID 和新的状态，调用一次只允许更新一个状态
// @Accept      multipart/form-data
// @Param       body_params body dtos.UpdateProjectStatusRequest true "更新项目状态的信息"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/update_status [put]
func (h *ProjectProcHandler) UpdateStatus(ctx iris.Context) {
}

// ProjectLaborHandler 处理项目角色相关的请求
type ProjectLaborHandler struct {
}

// InviteMember godoc
// @Summary     邀请成员加入项目
// @Description 邀请成员加入项目，需提供成员的 ID、项目 ID 和邀请职位
// @Accept      application/json
// @Param       body_params body dtos.InviteMemberRequest true "邀请成员加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/invite [post]
func (h *ProjectLaborHandler) Invite(ctx iris.Context) {
}

// Accept godoc
// @Summary     接受邀请加入项目
// @Description 接受邀请加入项目，需提供邀请的 ID
// @Accept      application/json
// @Param       body_params body dtos.AcceptInvitationRequest true "接受邀请加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/accept [post]
func (h *ProjectLaborHandler) Accept(ctx iris.Context) {
}

// Refuse godoc
// @Summary     拒绝邀请加入项目
// @Description 拒绝邀请加入项目，需提供邀请的 ID
// @Accept      application/json
// @Param       body_params body dtos.RefuseInvitationRequest true "拒绝邀请加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/refuse [post]
func (h *ProjectLaborHandler) Refuse(ctx iris.Context) {
}

// Apply godoc
// @Summary     申请加入项目
// @Description 申请加入项目，需提供项目 ID 和申请的职位
// @Accept      application/json
// @Param       body_params body dtos.ApplyProjectRequest true "申请加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/apply [post]
func (h *ProjectLaborHandler) Apply(ctx iris.Context) {
}

// AcceptAppli godoc
// @Summary     接受申请加入项目
// @Description 接受申请加入项目，需提供申请的 ID
// @Accept      application/json
// @Param       body_params body dtos.AcceptApplicationRequest true "接受申请加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/accept_appli [post]
func (h *ProjectLaborHandler) AcceptAppli(ctx iris.Context) {
}

// RefuseAppli godoc
// @Summary     拒绝申请加入项目
// @Description 拒绝申请加入项目，需提供申请的 ID
// @Accept      application/json
// @Param       body_params body dtos.RefuseApplicationRequest true "拒绝申请加入项目的请求体"
// @Tags        project_labor
// @Produce     json
// @Success     200
// @Router      /project/labor/refuse_appli [post]
func (h *ProjectLaborHandler) RefuseAppli(ctx iris.Context) {

}
