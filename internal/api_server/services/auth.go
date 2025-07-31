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

// LoginResult 定义了登录的结果
type LoginResult struct {
	UserInfo   UserInfo // 绑定后的用户信息
	PoplarJWT  string   // PoplarGrid 的 token
	MoetranJWT string   // 龙译的 JWT
}

// RegisterParams 定义了注册请求的参数
type RegisterParams struct {
	Email    string // 账号邮箱
	Password string // 账号密码
	Nickname string // 昵称
	VCode    string // 龙译发送的邮箱验证码
	QQNumber *int   // QQ 号，可选
}

// RegisterResult 定义了绑定的结果
// 其本质上与 LoginResult 相同
type RegisterResult LoginResult

// BindParams 定义了绑定请求的参数
// 其本质上与 LoginParams 相同
type BindParams LoginParams

// BindResult 定义了绑定的结果
// 其本质上与 LoginResult 相同
type BindResult LoginResult

// AuthService 定义了鉴权服务的接口
type AuthService interface {
	// Login 用户登录，返回用户信息并生成 JWT
	Login(params *LoginParams) (*LoginResult, error)
	// Register 用户注册，返回用户信息和 JWT
	Register(params *RegisterParams) (*RegisterResult, error)
	// Bind 绑定已有的龙译账号到 PoplarGrid 用户
	Bind(params *BindParams) (*BindResult, error)
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
func (s *authServiceImpl) Login(params *LoginParams) (*LoginResult, error) {
	if params.Email == "" ||
		params.Password == "" ||
		params.Captcha == "" ||
		params.CaptchaInfo == "" {
		return nil, errors.New("缺少必需的登录参数")
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
		return nil, errors.New("调用龙译 API 登录失败")
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
			return nil, errors.New("用户不存在或密码错误")
		}
		s.logger.Error("Login 查询用户失败", slog.Any("error", err))
		return nil, errors.New("查询用户失败")
	}

	// 如果哈希密码不为空，验证密码是否正确
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(params.Password),
	); err != nil {
		// 密码哈希不匹配，返回错误
		s.logger.Warn("Login 密码不匹配", slog.String("email", params.Email))
		return nil, errors.New("密码错误")
	}

	// 更新用户的 MoetranJwt
	if err := models.GetUser().Update(s.handle, &models.User{
		BaseModel: models.BaseModel{
			Id: user.Id,
		},
		MoetranJwt: moetranJWT,
	}); err != nil {
		s.logger.Error("Login 更新用户 MoetranJwt 失败", slog.Any("error", err))
		return nil, errors.New("更新用户 MoetranJwt 失败")
	}

	// 查找用户对应的 member IDs 用以加入到 JWT 中
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
		return nil, errors.New("查询用户成员信息失败，但是尨译账号登录成功")
	}

	// 生成 PoplarGrid 的 token
	authToken := &AuthToken{
		UserId:     uint(user.Id),
		MemberIds:  make([]uint, 0, len(members)),
		MoetranJwt: moetranJWT,
	}

	for _, member := range members {
		authToken.MemberIds = append(authToken.MemberIds, uint(member.Id))
	}

	poplarToken, err := s.tokenFactory.GenerateToken(authToken)
	if err != nil {
		s.logger.Error("Login 生成 JWT 失败", slog.Any("error", err))
		return nil, errors.New("生成 JWT 失败，但是尨译账号登录成功")
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

	return &LoginResult{
		UserInfo:   *userInfo,
		PoplarJWT:  poplarToken,
		MoetranJWT: moetranJWT,
	}, nil
}

// Register 实现 AuthService 接口的 Register 方法
func (s *authServiceImpl) Register(params *RegisterParams) (*RegisterResult, error) {
	if params.Email == "" ||
		params.Password == "" ||
		params.Nickname == "" ||
		params.VCode == "" {
		return nil, errors.New("缺少必需的注册参数")
	}

	// 先向尨译发送注册请求
	moetranJWT, err := s.apiClient.Register(&apiclient.RegisterParams{
		Email:    params.Email,
		Password: params.Password,
		Name:     params.Nickname,
		VCode:    params.VCode,
	})
	if err != nil {
		s.logger.Error("Register 调用龙译 API 注册失败", slog.Any("error", err))
		return nil, errors.New("调用龙译 API 注册失败")
	}

	// 创建用户
	user := &models.User{
		Email:    params.Email,
		Nickname: params.Nickname,
		QQNumber: params.QQNumber,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Register 哈希密码失败", slog.Any("error", err))
		return nil, errors.New("哈希密码失败，但是尨译账号注册成功，可以登录")
	}

	user.PasswordHash = string(hashedPassword)

	if err := models.GetUser().Insert(s.handle, user); err != nil {
		s.logger.Error("Register 创建用户失败", slog.Any("error", err))
		return nil, errors.New("注册用户失败，但是尨译账号注册成功，可以登录")
	}

	// 生成 PoplarGrid 的 token
	authToken := &AuthToken{
		UserId:     uint(user.Id),
		MoetranJwt: moetranJWT,
	}

	poplarToken, err := s.tokenFactory.GenerateToken(authToken)
	if err != nil {
		s.logger.Error("Register 生成 JWT 失败", slog.Any("error", err))
		return nil, errors.New("生成 JWT 失败，但是尨译账号注册成功，可以登录")
	}

	userInfo := &UserInfo{
		Id:       uint(user.Id),
		Nickname: user.Nickname,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
	if user.QQNumber != nil {
		userInfo.QQNumber = *user.QQNumber
	}

	return &RegisterResult{
		UserInfo:   *userInfo,
		PoplarJWT:  poplarToken,
		MoetranJWT: moetranJWT,
	}, nil
}

// Bind 实现 AuthService 接口的 Bind 方法
func (s *authServiceImpl) Bind(params *BindParams) (*BindResult, error) {
	if params.Email == "" ||
		params.Password == "" ||
		params.Captcha == "" ||
		params.CaptchaInfo == "" {
		return nil, errors.New("缺少必需的绑定参数")
	}

	// 调用尨译进行登录验证
	moetranJWT, err := s.apiClient.Login(&apiclient.LoginParams{
		Email:       params.Email,
		Password:    params.Password,
		Captcha:     params.Captcha,
		CaptchaInfo: params.CaptchaInfo,
	})
	if err != nil {
		s.logger.Error("Bind 调用龙译 API 登录失败", slog.Any("error", err))
		return nil, errors.New("调用龙译 API 登录失败")
	}

	// 获取龙译用户的信息
	userInfo, err := s.apiClient.GetUserInfo(moetranJWT)
	if err != nil {
		s.logger.Error("Bind 调用龙译 API 获取用户信息失败",
			slog.Any("error", err),
			slog.Any("message", userInfo.Error))
		return nil, errors.New("调用龙译 API 获取用户信息失败")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Bind 哈希密码失败", slog.Any("error", err))
		return nil, errors.New("哈希密码失败，但是尨译账号登录成功")
	}

	// 将信息绑定进入本地数据库
	newUser := &models.User{
		Nickname:     userInfo.Response.Name,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
		MoetranId:    userInfo.Response.ID,
		MoetranJwt:   moetranJWT,
	}

	if err := models.GetUser().Insert(s.handle, newUser); err != nil {
		s.logger.Error("Bind 绑定用户失败", slog.Any("error", err))
		return nil, errors.New("绑定用户失败，但是尨译账号登录成功")
	}

	// 生成 PoplarGrid 的 token
	authToken := &AuthToken{
		UserId:     uint(newUser.Id),
		MoetranJwt: moetranJWT,
	}
	poplarToken, err := s.tokenFactory.GenerateToken(authToken)
	if err != nil {
		s.logger.Error("Bind 生成 JWT 失败", slog.Any("error", err))
		return nil, errors.New("生成 JWT 失败，但是尨译账号登录成功")
	}

	r := &BindResult{
		UserInfo: UserInfo{
			Id:       uint(newUser.Id),
			Nickname: newUser.Nickname,
			Email:    newUser.Email,
			IsAdmin:  newUser.IsAdmin,
		},
		PoplarJWT:  poplarToken,
		MoetranJWT: moetranJWT,
	}
	if newUser.QQNumber != nil {
		r.UserInfo.QQNumber = *newUser.QQNumber
	}

	return r, nil
}
