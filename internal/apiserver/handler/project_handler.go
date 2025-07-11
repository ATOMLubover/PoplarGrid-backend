package handler

import (
	"poplargrid/internal/apiserver/dto"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteProjectHandler 注册项目相关的路由
func RouteProjectHandler(app *mvc.Application) {
	// 创建 ProjectHandler 实例
	handler := &ProjectHandler{}

	// 注册路由组
	party := app.Party("/project")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Get("/stats", handler.ProjectStats)
	party.Router.Get("/list", handler.ProjectListPage)
	party.Router.Get("/detail", handler.ProjectDetail)
}

// ProjectHandler 处理项目相关的请求
type ProjectHandler struct {
}

// ProjectStats godoc
// @Summary 	获取项目统计信息
// @Description 获取所有项目的统计信息，包括总数、翻译进行/完成、校对进行/完成、嵌字进行/完成，审核进行/完成、发布完成对应数量等
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dto.ProjectStats
// @Router 		/project/stats [get]
func (h *ProjectHandler) ProjectStats(ctx iris.Context) {

}

// ProjectListPage godoc
// @Summary 	获取项目列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		pageSerial query int false "页码，默认值为 1"
// @Param 		pageSize query int false "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_desc | update_desc，其他输入无效"
// @Param 		status query string false "项目状态，默认为 all，支持 all | translating | translated | prooving | prooved | lettering | lettered | published，其他输入无效"
// @Param 		tag_id query string false "特定标签，默认为空，表示不筛选"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} []dto.ProjectBasic
// @Router 		/project/list [get]
func (h *ProjectHandler) ProjectListPage(ctx iris.Context) {
	// 解析查询参数
	var params dto.ProjectListQueryParams

	if err := ctx.ReadQuery(&params); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":  "查询参数解析错误",
			"detail": err.Error(),
		})
	}

	// 校验查询参数
	if err := params.Validate(); err != nil {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":  "查询参数无效",
			"detail": err.Error(),
		})
		return
	}
}

// ProjectDetail godoc
// @Summary 	获取项目详情
// @Description 获取指定项目的详细信息，包括翻译、校对、嵌字、审核、发布等状态
// @Param 		id query string true "项目 ID"
// @Tags 		project
// @Produce 	json
// @Success	 	200 {object} dto.ProjectDetail
// @Router 		/project/{id} [get]
func (h *ProjectHandler) ProjectDetail(ctx iris.Context) {

}
