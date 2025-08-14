package handlers

import (
	"fmt"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// UserInfo 定义了用户的基本信息 DTO
type UserInfo struct {
	// 用户 ID
	ID uint `json:"id"`
	// 昵称
	Nickname string `json:"nickname"`
	// 邮箱
	Email string `json:"email,omitempty"`
	// QQ 号
	QQNumber int `json:"qq_number,omitempty"`
	// 是否是 panel 管理员
	IsAdmin bool `json:"is_admin"`
	// 补充备注
	Remark string `json:"remark,omitempty"`
	// 龙译 ID
	MoetranId string `json:"moetran_id,omitempty"`
	// 龙译 JWT
	MoetranJwt string `json:"moetran_jwt,omitempty"`
}

// UserDetail 定义了用户详情的 DTO
type UserDetail struct {
	// 用户的基本信息
	User *UserInfo `json:",inline"`
	// 在各个汉化组中的成员信息
	Members []MemberInfo `json:"members,omitempty"`
}

// RouteUserHandler 注册用户相关的路由
func RouteUserHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/users").
		Handle(new(UserHandler))
}

// UserHandler 处理用户相关的请求
type UserHandler struct {
	UserService services.UserService
}

// BeforeActivation 在控制器激活前注册路由和中间件
func (h *UserHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/me", "MyDetail")
	b.Handle("GET", "/{id:uint}", "Detail")

	// 注册控制器特定的中间件
	b.Router().Use(NewCheckUserIdMiddleware())
}

// MyDetail godoc
//
// @Summary 	利用 cookie，获取当前用户的详情
// @Description 获取指定用户的详细信息，本质上是一次重定向到 /users/{id}/detail
//
// @Tags 		user
// @Produce 	json
// @Success	 	307
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
//
// @Router 		/api/users/me [get]
func (h *UserHandler) MyDetail(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无法获取有效的 user_id"))
		return
	}

	// 直接重定向到 UserDetail 方法
	url := fmt.Sprintf("/users/%d/detail", userId)
	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

// Detail godoc
//
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、昵称、QQ 等
// @Param 		id path uint true "用户 ID"
//
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} FormatResponse[UserInfo]
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/api/users/{id} [get]
func (h *UserHandler) Detail(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无法获取有效的 user_id"))
		return
	}

	// 调用服务层获取用户详情
	user, e := h.UserService.GetUserDetail(userId)
	if e != nil {
		wrapError(ctx, e)
		return
	}

	// 将服务层的 UserInfo 转换为 DTO
	res := &UserDetail{
		User: &UserInfo{
			ID:         user.ID,
			Nickname:   user.Nickname,
			Email:      user.Email,
			QQNumber:   user.QQNumber,
			IsAdmin:    user.IsAdmin,
			Remark:     user.Remark,
			MoetranId:  user.MoetranId,
			MoetranJwt: user.MoetranJwt,
		},
	}

	// 如果用户有成员信息，则填充
	if user.Members != nil {
		res.Members = make([]MemberInfo, len(user.Members))

		for i, member := range user.Members {
			res.Members[i] = MemberInfo{
				Id: member.Id,
				Team: &TeamInfo{
					Id:   member.Team.Id,
					Name: member.Team.Name,
				},
				Role: member.Role,
			}
		}
	}

	wrapSuccess(ctx, res)
}
