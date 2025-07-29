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
	Id uint `json:"id"`
	// 昵
	Nickname string `json:"nickname"`
	// 邮箱
	Email string `json:"email,omitempty"`
	// QQ 号
	QQNumber int `json:"qq_number,omitempty"`
	// 是否是 panel 管理员
	IsAdmin bool `json:"is_admin"`
	// 补充备注
	Remark string `json:"remark,omitempty"`
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
// @Success	 	200 {object} UserInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/users/me [get]
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

// Detail godoc
//
// @Summary 	获取用户详情
// @Description 获取指定用户的详细信息，包括 ID、昵称、QQ 等
// @Param 		id path uint true "用户 ID"
//
// @Tags 		user
// @Produce 	json
// @Success	 	200 {object} UserInfo
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/api/users/{id} [get]
func (h *UserHandler) Detail(ctx iris.Context) {
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
	user, err := h.UserService.GetUserDetail(userId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(ErrorResponse{
			Error:  "获取特定用户详情失败",
			Detail: err.Error(),
		})
		return
	}

	// 将服务层的 UserInfo 转换为 DTO
	ctx.JSON(UserInfo{
		Id:       user.Id,
		Nickname: user.Nickname,
		Email:    user.Email,
		QQNumber: user.QQNumber,
		IsAdmin:  user.IsAdmin,
		Remark:   user.Remark,
	})
}
