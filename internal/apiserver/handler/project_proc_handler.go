package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteProjectProcHandler 注册项目进度动态推进相关的路由
func RouteProjectProcHandler(app *mvc.Application) {
	// 创建 ProjectProcHandler 实例
	handler := &ProjectProcHandler{}

	// 注册路由组
	party := app.Party("/proc")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
}

// ProjectProcHandler 处理与项目进度动态推进有关的路由
type ProjectProcHandler struct {
}

// Create godoc
// @Summary 	创建项目
// @Description 创建一个新的项目，请求体暂时未确定
// @Accept      application/json
// @Param       request body dto.CreateProjectRequest true "创建项目的请求体"
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
// @Param       id formData string true "项目 ID"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/delete [delete]
func (h *ProjectProcHandler) Delete(ctx iris.Context) {
}

// Update godoc
// @Summary 	更新项目
// @Description 更新指定的项目，需提供项目 ID 和更新的字段
// @Accept      multipart/form-data
// @Param       id formData string true "项目 ID"
// @Param       title formData string false "项目标题"
// @Param       description formData string false "项目描述"
// @Param       tags formData []string false "标签列表"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/update [put]
func (h *ProjectProcHandler) Update(ctx iris.Context) {
}

// UpdateStatus godoc
// @Summary 	更新项目状态
// @Description 更新指定项目的状态，需提供项目 ID 和新的状态（之所以没有 translating 之类的状态，是因为这是翻译等一加入就自动进入的）
// @Accept      multipart/form-data
// @Param       id formData string true "项目 ID"
// @Param       status formData string true "新的项目状态，只能是 translated | prooved | lettered | reviewed | published"
// @Tags 		project_proc
// @Produce 	json
// @Success	 	200 {object} map[string]string
// @Router 		/project/proc/status [put]
func (h *ProjectProcHandler) UpdateStatus(ctx iris.Context) {

}
