package services

import (
	"errors"
	"log/slog"
	"poplargrid/internal/api_server/apiclient"
	"poplargrid/internal/shared/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginParams 定义了登录请求的参数
type LoginParams struct {
	Email       string // 账号邮箱
	Password    string // 账号密码
	Captcha     string // 验证码
	CaptchaInfo string // 验证码信息
}

// AuthService 定义了鉴权服务的接口
type AuthService interface {
	// Login 用户登录，返回用户信息并生成 JWT token
	Login(params *LoginParams) (info *UserInfo, poplarToken, moetranToken string, err error)
	// Register 用户注册，返回用户信息和 JWT token
	// Register(params *RegisterParams) (*UserInfo, string, error)
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

// Login 实现 AuthService 接口的 Login 方法
func (s *authServiceImpl) Login(params *LoginParams) (*UserInfo, string, string, error) {
	if params.Email == "" ||
		params.Password == "" ||
		params.Captcha == "" ||
		params.CaptchaInfo == "" {
		return nil, "", "", errors.New("缺少必需的登录参数")
	}

	// 验证用户的邮箱和密码
	userSpec := &models.UserSpec{
		Email: &params.Email,
	}
	userFields := &models.UserFields{
		Id:           true,
		Email:        true,
		Nickname:     true,
		PasswordHash: true,
		QQNumber:     true,
		IsAdmin:      true,
	}

	user, err := models.GetUser().SelectFirst(s.handle, userSpec, userFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Login 查询没有结果", slog.String("email", params.Email))
			return nil, "", "", errors.New("用户不存在或密码错误")
		}
		s.logger.Error("Login 查询用户失败", slog.Any("error", err))
		return nil, "", "", errors.New("查询用户失败")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(params.Password),
	); err != nil {
		// 密码哈希不匹配，返回错误
		s.logger.Warn("Login 密码不匹配", slog.String("email", params.Email))
		return nil, "", "", errors.New("密码错误")
	}

	// 调用尨译进行登录验证
	moetranJWT, err := s.apiClient.Login(&apiclient.LoginParams{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	})
	if err != nil {
		s.logger.Error("Login 调用龙译 API 登录失败", slog.Any("error", err))
		return nil, "", "", errors.New("调用龙译 API 登录失败")
	}

	// 更新用户的 MoetranJwt
	if err := models.GetUser().Update(s.handle, &models.User{
		BaseModel: models.BaseModel{
			Id: user.Id,
		},
		MoetranJwt: moetranJWT,
	}); err != nil {
		s.logger.Error("Login 更新用户 MoetranJwt 失败", slog.Any("error", err))
		return nil, "", "", errors.New("更新用户 MoetranJwt 失败")
	}

	// 查找用户对应的 member IDs 用以加入到 JWT token 中
	memberSpec := &models.MemberSpec{
		UserId: &user.Id,
	}
	memberFields := &models.MemberFields{
		Id: true,
	}

	members, err := models.GetMember().SelectMany(
		s.handle, memberSpec, memberFields,
		// 这里不需要 offset 和 limit，因为需要获取用户的所有成员 ID
		nil, nil)
	if err != nil {
		s.logger.Error("Login 查询用户成员失败", slog.Any("error", err))
		return nil, "", "", errors.New("查询用户成员信息失败")
	}

	// 生成 JWT token
	authToken := AuthToken{
		UserId:    uint(user.Id),
		MemberIds: make([]uint, 0, len(members)),
	}

	for _, member := range members {
		authToken.MemberIds = append(authToken.MemberIds, uint(member.Id))
	}

	poplarJWT, err := s.tokenFactory.GenerateToken(&authToken)
	if err != nil {
		s.logger.Error("Login 生成 JWT token 失败", slog.Any("error", err))
		return nil, "", "", errors.New("生成 JWT token 失败")
	}

	// 返回用户信息
	userInfo := &UserInfo{
		Id:       uint(user.Id),
		Nickname: user.Nickname,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
	if user.QQNumber != nil {
		userInfo.QQNumber = *user.QQNumber
	}

	return userInfo, poplarJWT, moetranJWT, nil
}
