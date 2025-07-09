package handler

import (
	"poplargrid/internal/api_server/dto"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// MemberHandler 处理成员相关的请求
type MemberHandler struct {
	MemberService *services.MemberService
}

// BeforeActivation 注册 MemberHandler 的路由与方法的映射
func (h *MemberHandler) BeforeActivation(b mvc.BeforeActivation) {
	// 注册成员列表分页的路由
	b.Handle(iris.MethodGet, "/list", "MemberListPage")
}

// MemberListPage godoc
// @Summary 	获取成员详细信息的列表分页，如果想要获取全部内容可以设置一个很大的 pageSize
// @Description 注意当列表为空，会返回 null 费不是空数组
// @Param 		pageSerial query int false "页码，默认值为 1"
// @Param 		pageSize query int false "每页数量，默认值为 10"
// @Tags 		member
// @Produce 	json
// @Success	 	200 {object} []dtos.MemberFullInfo
// @Router 		/member/list [get]
func (h *MemberHandler) MemberListPage(ctx iris.Context) {
	// 获取分页参数
	pageSerial := ctx.URLParamIntDefault("pageSerial", 1)
	pageSize := ctx.URLParamIntDefault("pageSize", 10)

	// 确保分页参数合法
	if pageSerial < 1 {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":      "pageSerial非法",
			"pageSerial": pageSerial,
		})
		return
	}
	if pageSize < 1 {
		ctx.StopWithJSON(iris.StatusBadRequest, iris.Map{
			"error":      "pageSize非法",
			"pageSerial": pageSerial,
		})
	}

	// 处理分页参数为 offset 和 limit
	offset := (pageSerial - 1) * pageSize
	limit := pageSize

	// 调用服务获取成员列表
	members, err := h.MemberService.GetMemberFullList(offset, limit)
	if err != nil {
		ctx.StopWithJSON(iris.StatusInternalServerError, iris.Map{
			"error":   "获取成员列表失败",
			"details": err.Error(),
		})
		return
	}

	// 成功返回成员详情列表
	var responseSlice []*dto.MemberFullInfo

	for _, member := range members {
		responseSlice = append(responseSlice, &dto.MemberFullInfo{
			Nickname:      member.Nickname,
			Email:         member.Email,
			MoetranId:     member.MoetranId,
			PoplarIsAdmin: member.PoplarIsAdmin,
			Labors:        member.Labors,
			Remark:        member.Remark,
			LastActive:    member.LastActive.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(responseSlice)
}
