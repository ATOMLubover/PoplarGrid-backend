package handlers

import (
	"poplargrid/internal/api_server/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// LoginParams 定义了登录的请求参数
type LoginParams struct {
	// Email 是龙译账号的邮箱
	Email string `json:"email" binding:"required"`
	// Password 是龙译账号的密码
	Password string `json:"password" binding:"required"`
	// Captcha 是验证码
	Captcha string `json:"captcha" binding:"required"`
	// CaptchaInfo 是验证码信息
	CaptchaInfo string `json:"captcha_info" binding:"required"`
}

// LoginResponse 定义了登录的响应结果
type LoginResponse struct {
	// MoetranJWT 是登录成功后返回的 JWT token
	MoetranJWT string `json:"token"`
	// User 是登录成功后返回的用户信息
	User UserInfo `json:"user"`
}

// RouteAuthHandler 注册鉴权相关的路由
func RouteAuthHandler(root *mvc.Application) {
	// 注册路由组和 handler
	root.Party("/auth").
		Handle(new(AuthHandler))
}

// AuthHandler 处理鉴权相关的请求
// 包括用户登录、注册等
type AuthHandler struct {
	AuthService services.AuthService
}

// BeforeActivation 在控制器激活前注册路由
func (h *AuthHandler) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/login", "Login")
	// b.Handle("POST", "/register", "Register")
}

// Login godoc
// @Summary 	登录账号
// @Description 登录 PoplarGrid 以及龙译的账号，返回龙译的 JWT token
//
// @Accept      json
// @Param 		params body LoginParams true "登录参数"
//
// @Tags 		auth
// @Produce 	json
// @Success	 	200 {object} LoginResponse
// @Failure     400 {object} ErrorResponse "无效的请求参数"
// @Failure     500 {object} ErrorResponse "服务器内部错误"
//
// @Router 		/auth/login [post]
func (h *AuthHandler) Login(ctx iris.Context) {
	var params LoginParams

	if err := ctx.ReadJSON(&params); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error: "无效的请求参数",
		})
		return
	}

	// 调用服务层进行登录处理
	user, poplarToken, moetranToken, err := h.AuthService.Login(&services.LoginParams{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(ErrorResponse{
			Error:  "登录失败",
			Detail: err.Error(),
		})
		return
	}

	// 为响应添加 Set-Cookie 头部
	ctx.SetCookie(&iris.Cookie{
		Name:     "poplar_token",
		Value:    poplarToken,
		Expires:  time.Now().Add(48 * time.Hour),
		HttpOnly: true,
		// TODO: 测试环境不启用 HTTPS
		// Secure:   true,
	})

	ctx.JSON(LoginResponse{
		MoetranJWT: moetranToken,
		User: UserInfo{
			Id:       user.Id,
			Nickname: user.Nickname,
			Email:    user.Email,
			QQNumber: user.QQNumber,
			IsAdmin:  user.IsAdmin,
		},
	})
}
