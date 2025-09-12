package services

import (
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BindParams 定义了绑定请求的参数
type BindParams struct {
	Email       string // 账号邮箱
	Password    string // 账号密码
	Captcha     string // 验证码
	CaptchaInfo string // 验证码信息
}

// BindResult 定义了绑定的结果
type BindResult struct {
	UserInfo   UserInfo // 绑定后的用户信息
	PoplarJWT  string   // PoplarGrid 的 token
	MoetranJWT string   // 龙译的 JWT
}

// AuthService 定义了鉴权服务的接口
type AuthService interface {
	// // Login 用户登录，返回用户信息并生成 JWT
	// Login(params LoginParams) (*LoginResult, Error)
	// // Register 用户注册，返回用户信息和 JWT
	// Register(params RegisterParams) (*RegisterResult, Error)
	// Bind 绑定已有的龙译账号到 PoplarGrid 用户
	Bind(params *BindParams) (*BindResult, Err)
}

// authServiceImpl 实现了 AuthService 接口
type authServiceImpl struct {
	tokenFactory AuthTokenFactory
	handle       *gorm.DB
	apiClient    apiclient.ApiClient
	logger       *slog.Logger
}

// NewAuthService 创建一个新的 AuthService 实例
func NewAuthService(
	factory AuthTokenFactory,
	hdl *gorm.DB,
	apiClient apiclient.ApiClient,
	lgr *slog.Logger,
) AuthService {
	return &authServiceImpl{
		tokenFactory: factory,
		handle:       hdl,
		apiClient:    apiClient,
		logger:       lgr,
	}
}

// Bind 实现 AuthService 接口的 Bind 方法
func (s *authServiceImpl) Bind(params *BindParams) (*BindResult, Err) {
	if params.Email == "" ||
		params.Password == "" ||
		params.Captcha == "" ||
		params.CaptchaInfo == "" {
		return nil, ErrParamsLackage
	}

	// 调用尨译进行登录验证
	moetranJWT, err := s.apiClient.Login(&apiclient.LoginParams{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	})
	if err != nil {
		s.logger.Error("Bind 调用龙译 API 登录失败", slog.Any("Error", err))
		return nil, ErrMoetranAPIFailure
	}

	// 获取龙译用户的信息
	userInfo, err := s.apiClient.GetUserInfo(moetranJWT)
	if err != nil {
		s.logger.Error("Bind 调用龙译 API 获取用户信息失败",
			slog.Any("Error", err),
			slog.Any("message", userInfo.Error))
		return nil, ErrMoetranAPIFailure
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Bind 哈希密码失败", slog.Any("Error", err))
		return nil, ErrMoetranAPIFailure
	}

	// 将信息绑定进入本地数据库
	newUser := &models.User{
		Nickname:     userInfo.Response.Name,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
		MoetranId:    userInfo.Response.ID,
		MoetranJwt:   moetranJWT,
	}

	if err := s.handle.Model(newUser).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "moetran_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"nickname"}),
		}).
		Create(newUser); err != nil {
		s.logger.Error("Bind 绑定用户失败",
			slog.Any("Error", err))
		return nil, ErrDatabaseFailure
	}

	// 生成 PoplarGrid 的 token
	authToken := &AuthToken{
		UserId:     uint(newUser.Id),
		MoetranJwt: moetranJWT,
	}
	poplarToken, err := s.tokenFactory.GenerateToken(authToken)
	if err != nil {
		s.logger.Error("Bind 生成 JWT 失败", slog.Any("Error", err))
		return nil, ErrTokenGenerationFailure
	}

	r := &BindResult{
		UserInfo: UserInfo{
			ID:       uint(newUser.Id),
			Nickname: newUser.Nickname,
			Email:    newUser.Email,
			IsAdmin:  newUser.IsAdmin,
		},
		PoplarJWT:  poplarToken,
		MoetranJWT: moetranJWT,
	}
	if newUser.QQNumber != nil {
		r.UserInfo.QQNumber = int(*newUser.QQNumber)
	}

	return r, nil
}
