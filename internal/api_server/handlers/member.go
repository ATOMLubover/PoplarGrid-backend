package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12/mvc"
)

// MmeberHandler 处理成员相关的请求
type MemberHandler struct {
	memberService services.MemberService
}

// MemberInfo 定义了用户在团队中的基本信息 DTO
type MemberInfo struct {
	// 成员 ID
	// @example 123456
	Id uint `json:"id"`
	// 对应的用户信息
	User UserInfo `json:"user"`
	// 对应的汉化组信息
	// @example 114514
	Team TeamInfo `json:"team"`
	// 在组内的职责（掩码格式）
	// @example 5
	Role services.RoleMask `json:"role"`
}

// BeforeActivation 在控制器激活前注册路由
func (h *MemberHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "MemberList")
}

// MemberList godoc
//
// @Summary 	获取指定用户的成员列表
// @Description 获取指定用户在各个汉化组中的成员信息
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		user_id query uint true "用户 ID"

