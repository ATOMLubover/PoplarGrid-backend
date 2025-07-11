package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteProjectRoleHandler 注册项目角色相关的路由
func RouteProjectRoleHandler(app *mvc.Application) {
	// 创建 ProjectRoleHandler 实例
	handler := &ProjectRoleHandler{}

	// 注册路由组
	party := app.Party("/role")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Post("/invite", handler.Invite)
	party.Router.Post("/accept_appli", handler.AcceptAppli)
	party.Router.Post("/refuse_appli", handler.RefuseAppli)

	party.Router.Post("/apply", handler.Apply)
	party.Router.Post("/accept", handler.Accept)
	party.Router.Post("/refuse", handler.Refuse)
}

// ProjectRoleHandler 处理项目角色相关的请求
type ProjectRoleHandler struct {
}

// InviteMember godoc
// @Summary 	邀请成员加入项目
// @Description 邀请成员加入项目，需提供成员的 ID
// @Accept      multipart/form-data
// @Param 		member_id formData string true "成员的 ID"
// @Param       invite_role formData string true "邀请加入的职位"
// @Param       project_id formData string true "项目 ID"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/invite [post]
func (h *ProjectRoleHandler) Invite(ctx iris.Context) {

}

// Accept godoc
// @Summary 	接受邀请加入项目
// @Description 接受邀请加入项目，需提供邀请的 ID
// @Accept      multipart/form-data
// @Param       invitition_id formData string true "邀请的 ID"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/accept [post]
func (h *ProjectRoleHandler) Accept(ctx iris.Context) {
}

// Refuse godoc
// @Summary 	拒绝邀请加入项目
// @Description 拒绝邀请加入项目，需提供邀请的 ID
// @Accept      multipart/form-data
// @Param       invitition_id formData string true "邀请的 ID"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/refuse [post]
func (h *ProjectRoleHandler) Refuse(ctx iris.Context) {
}

// Apply godoc
// @Summary 	申请加入项目
// @Description 申请加入项目，需提供项目 ID 和申请的职位
// @Accept      multipart/form-data
// @Param       project_id formData string true "项目 ID"
// @Param       apply_role formData string true "申请加入的职位"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/apply [post]
func (h *ProjectRoleHandler) Apply(ctx iris.Context) {
}

// AcceptAppli godoc
// @Summary 	接受申请加入项目
// @Description 接受申请加入项目，需提供申请的 ID
// @Accept      multipart/form-data
// @Param       application_id formData string true "申请的 ID"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/accept_appli [post]
func (h *ProjectRoleHandler) AcceptAppli(ctx iris.Context) {
}

// RefuseAppli godoc
// @Summary 	拒绝申请加入项目
// @Description 拒绝申请加入项目，需提供申请的 ID
// @Accept      multipart/form-data
// @Param       application_id formData string true "申请的 ID"
// @Tags 		project_role
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/role/refuse_appli [post]
func (h *ProjectRoleHandler) RefuseAppli(ctx iris.Context) {
}
