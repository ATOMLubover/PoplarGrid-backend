package handlers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// RouteMemberHandler 注册成员相关的路由
func RouteMemberHandler(root *mvc.Application) {
	// 创建 MemberHandler 实例
	handler := &MemberHandler{}

	// 注册路由组
	party := root.Party("/member")

	// 注册 handler
	party.Handle(handler)

	// 注册路由与方法的映射（类型安全）
	party.Router.Get("/list", handler.MemberListPage)
}

// MemberHandler 处理成员相关的请求
type MemberHandler struct {
}

// MemberListPage godoc
// @Summary 	获取成员列表分页
// @Description 注意当列表为空，会返回 null 而不是空数组
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		sort query string false "排序方式，默认值为 id_desc，支持 id_asc | id_desc，其他输入无效"
// @Param 		team_id query int true "所属汉化组 ID"
// @Tags 		member
// @Produce 	json
// @Success	 	200 {object} []dtos.MemberBasic
// @Router 		/member/list [get]
func (h *MemberHandler) MemberListPage(ctx iris.Context) {
}
