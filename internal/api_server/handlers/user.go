package handlers

import (
	"fmt"
	"poplargrid/internal/api_server/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

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

// UserInfo 定义了用户的基本信息 DTO
type UserInfo struct {
	// 用户 ID
	// @example 123456
	Id uint `json:"id"`
	// 昵称
	// @example [influ3nza]翻校
	Nickname string `json:"nickname"`
	// 邮箱
	// @example 1919810@163.com
	Email string `json:"email,omitempty"`
	// QQ 号
	// @example 123456789
	QqNumber string `json:"qq_number,omitempty"`
	// 是否是 panel 管理员
	// @example true
	PoplarIsAdmin bool `json:"poplar_is_admin"`
	// 补充备注
	// @example 这是一个测试用户
	Remark string `json:"remark,omitempty"`
	// 在各个汉化组中的成员信息
	Members []MemberInfo `json:"members,omitempty"`
}

// BeforeActivation 在控制器激活前注册路由和中间件
func (h *UserHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/me", "MyDetail")
	b.Handle("GET", "/{id:uint}/detail", "UserDetail")

	// 注册控制器特定的中间件
	b.Router().Use(NewCheckUserIdMiddleware())
}

// UserDetail godoc
//
// @Summary 	利用 cookie，获取当前用户的详情
// @Description 获取指定用户的详细信息，本质上是一次重定向到 /users/{id}/detail
//
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} UserDetail
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/users/me [get]
func (h *UserHandler) MyDetail(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 直接重定向到 UserDetail 方法
	url := fmt.Sprintf("/users/%d/detail", userId)
	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

// UserDetail godoc
//
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、昵称、QQ 等
// @Param 		id path uint true "用户 ID"
//
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} UserDetail
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/users/{id}/detail [get]
func (h *UserHandler) UserDetail(ctx iris.Context) {
	// 从路径参数获取用户 ID
	userId, err := ctx.Params().GetUint("id")
	if err != nil || userId <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无法获取有效的 user_id",
		})
		return
	}

	// 调用服务层获取用户详情
	userDetail, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取特定用户详情失败",
			Detail: err.Error(),
		})
		return
	}

	ctx.JSON(userDetail)
}
