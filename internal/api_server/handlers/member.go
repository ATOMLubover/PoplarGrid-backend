package handlers

import (
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// LaborMask 定义了分工的掩码类型
// @Description 分工掩码与成员在汉化组中的职责相关，使用位掩码表示不同的职责
// @Description 第 0 位表示监制/负责人
// @Description 第 1 位表示图源
// @Description 第 2 位表示美工
// @Description 第 3 位表示翻译
// @Description 第 4 位表示校对
// @Description 第 5 位表示嵌字
// @Description 第 6 位表示嵌字审核
// @Description 第 7 位表示发布
type LaborMask = services.LaborMask

// MemberInfo 定义了用户在团队中的基本信息 DTO
type MemberInfo struct {
	// 成员 ID
	Id uint `json:"id"`
	// 对应的用户信息，可能为空
	User *UserInfo `json:"user,omitempty"`
	// 对应的汉化组信息，可能为空
	Team *TeamInfo `json:"team,omitempty"`
	// 在组内的职责（掩码格式）
	Role LaborMask `json:"role"`
}

// RouteMemberHandler 注册成员相关的路由
func RouteMemberHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/members").
		Handle(new(MemberHandler))
}

// MmeberHandler 处理成员相关的请求
type MemberHandler struct {
	memberService services.MemberService
}

// BeforeActivation 在控制器激活前注册路由
func (h *MemberHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/", "List")
}

// List godoc
// @Summary 	获取指定汉化组的成员列表
// @Description 获取指定汉化组的成员列表，注意当列表为空，会返回 null 而不是空数组
//
// @Param 		page_serial query int false "页码，默认值为 1"
// @Param 		page_size query int false "每页数量，默认值为 10"
// @Param 		team_id query int true "汉化组 ID"
//
// @Tags 		member
// @Produce 	json
// @Success 	200 {object} []MemberInfo
// @Failure 	400 {object} ErrorResponse "无效的请求参数"
// @Failure 	500 {string} string "服务器内部错误"
//
// @Router 		/api/members [get]
func (h *MemberHandler) List(ctx iris.Context) {
	pageSerial := ctx.URLParamIntDefault("page_serial", 1)
	pageSize := ctx.URLParamIntDefault("page_size", 10)

	// 调用服务获取成员列表
	members, err := h.memberService.GetMembers(&services.MemberListParams{
		Offset: (pageSerial - 1) * pageSize,
		Limit:  pageSize,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "获取成员列表失败",
		})
		return
	}

	// 将成员列表转换为 DTO
	memberInfos := make([]*MemberInfo, 0, len(members))

	for _, member := range members {
		memberInfos = append(memberInfos, &MemberInfo{
			Id: member.Id,
			User: &UserInfo{
				Id:       member.User.Id,
				Nickname: member.User.Nickname,
			},
			Team: &TeamInfo{
				Id:   member.Team.Id,
				Name: member.Team.Name,
			},
			Role: member.Role,
		})
	}

	ctx.JSON(members)
}
