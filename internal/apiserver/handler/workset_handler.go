package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// WorksetHandler 处理工作集相关的请求
type WorksetHandler struct {
}

// BeforeActivation 注册 WorksetHandler 的路由与方法的映射
func (h *WorksetHandler) BeforeActivation(b mvc.BeforeActivation) {
	// 注册工作集列表分页的路由
	b.Handle(iris.MethodGet, "/list", "WorksetListPage")
}

// WorksetListPage godoc
// @Summary 	获取工作集列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		pageSerial query int false "页码，默认值为 1"
// @Param 		pageSize query int false "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_asc | id_desc | poplarity_desc，其他输入无效"
// @Tags 		workset
// @Produce 	json
// @Success	 	200 {object} []dto.WorksetBasic
// @Router 		/workset/list [get]
func (h *WorksetHandler) WorksetListPage(ctx iris.Context) {
}
