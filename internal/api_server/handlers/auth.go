package handlers

import (
	"poplargrid/internal/api_server/config"
	"poplargrid/internal/api_server/services"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// BindParams 定义了绑定的请求参数
type BindParams struct {
	// Email 是龙译账号的邮箱
	Email string `json:"email" binding:"required"`
	// Password 是龙译账号的密码
	Password string `json:"password" binding:"required"`
	// Captcha 是验证码
	Captcha string `json:"captcha" binding:"required"`
	// CaptchaInfo 是验证码信息
	CaptchaInfo string `json:"captcha_info" binding:"required"`
}

// BindResponse 定义了绑定的响应结果
type BindResponse struct {
	// MoetranJWT 是登录成功后返回的 JWT
	MoetranJWT string `json:"token"`
	// User 是登录成功后返回的用户信息
	User UserInfo `json:"user"`
	// Members 是用户在各个汉化组中的成员信息
	Members []MemberInfo `json:"members,omitempty"`
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
	b.Handle("POST", "/bind", "Bind")
}

// Bind godoc
// @Summary 	绑定龙译账号
// @Description 绑定已有的龙译账号到 PoplarGrid 用户
//
// @Accept      json
// @Param 		body_params body BindParams true "绑定参数"
//
// @Tags 		auth
// @Produce 	json
// @Success	 	200 {object} FormatResponse[BindResponse]
// @Failure     400 {object} StringFormatResponse "无效的请求参数"
// @Failure     500 {string} string "服务器内部错误"
//
// @Router 		/auth/bind [post]
func (h *AuthHandler) Bind(ctx iris.Context) {
	var params BindParams

	if err := ctx.ReadJSON(&params); err != nil {
		wrapError(ctx, newHdlErr(ErrParamsLackage, "无效的请求参数"))
		return
	}

	// 调用服务层进行绑定处理
	result, err := h.AuthService.Bind(&services.BindParams{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	})
	if err != nil {
		wrapError(ctx, err)
		return
	}

	// 为响应添加 Set-Cookie 头部
	expiresAt := time.Now().Add(7 * time.Hour) // 默认为 7 天
	if cfg := config.GetConfig(); cfg != nil && cfg.Server.CookieLifetime > 0 {
		expiresAt = time.Now().Add(time.Duration(cfg.Server.CookieLifetime) * time.Second)
	}

	ctx.SetCookie(&iris.Cookie{
		Name:     "poplar_token",
		Value:    result.PoplarJWT,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
		// TODO: 测试环境不启用 HTTPS
		// Secure:   true,
	})

	wrapSuccess(ctx, BindResponse{
		User: UserInfo{
			ID:       result.UserInfo.ID,
			Nickname: result.UserInfo.Nickname,
			Email:    result.UserInfo.Email,
			QQNumber: result.UserInfo.QQNumber,
			IsAdmin:  result.UserInfo.IsAdmin,
		},
		MoetranJWT: result.MoetranJWT,
	})
}
