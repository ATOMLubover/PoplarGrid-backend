package handler

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// TeamHandler 处理汉化组信息相关的请求
type TeamHandler struct {
}

// BeforeActivation 注册 TeamHandler 的路由与方法的映射
func (h *TeamHandler) BeforeActivation(b mvc.BeforeActivation) {
	// 注册汉化组列表分页的路由
	b.Handle(iris.MethodGet, "/list", "TeamListPage")
}

// TeamListPage godoc
// @Summary 	获取汉化组列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		pageSerial query int false "页码，默认值为 1"
// @Param 		pageSize query int false "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_asc | id_desc，其他输入无效"
// @Tags 		team
// @Produce 	json
// @Success	 	200 {object} []dto.TeamBasic
// @Router 		/team/list [get]
func (h *TeamHandler) TeamListPage(ctx iris.Context) {

}
