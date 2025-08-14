package services

import (
	"poplargrid/internal/shared/jwtutils"
	"time"
)

// AuthToken 定义了 Poplar Grid 的 JWT 结构
type AuthToken struct {
	UserId     uint   // UserId 是用户的唯一标识符
	MemberIds  []uint // MemberIds 是用户在各个团队中的成员 ID 列表
	MoetranJwt string // MoetranJwt 是用户的 Moetran JWT
}

// AuthTokenFactory 定义了生成和解析 AuthToken 的接口
type AuthTokenFactory interface {
	// GenerateToken 生成一个新的 AuthToken 签名字符串
	GenerateToken(token *AuthToken) (string, Err)
	// ParseToken 解析一个 AuthToken 字符串，返回 AuthToken 对象
	ParseToken(tokenString string) (*AuthToken, Err)
}

// authTokenFactoryImpl 实现了 AuthTokenFactory 接口
type authTokenFactoryImpl struct {
	factory        *jwtutils.CliamsFactory[AuthToken] // 使用 jwtutils 的 CliamsFactory 生成和解析 JWT
	secretKey      string                             // 用于签名和验证的密钥
	expirationTime time.Duration                      // token 的过期时间
}

// NewAuthTokenFactory 创建一个新的 AuthToken 工厂
func NewAuthTokenFactory(secretKey string, expirationTime time.Duration) AuthTokenFactory {
	return &authTokenFactoryImpl{
		factory:        jwtutils.NewClaimsFactory[AuthToken](expirationTime, secretKey),
		secretKey:      secretKey,
		expirationTime: expirationTime,
	}
}

// GenerateToken 实现 AuthTokenFactory 接口的 GenerateToken 方法
func (f *authTokenFactoryImpl) GenerateToken(token *AuthToken) (string, Err) {
	tokenStr, err := f.factory.GenToken(*token)
	if err != nil {
		return "", newSrvError(ErrTokenGenerationFailure, err.Error())
	}
	return tokenStr, nil
}

// ParseToken 实现 AuthTokenFactory 接口的 ParseToken 方法
func (f *authTokenFactoryImpl) ParseToken(tokenString string) (*AuthToken, Err) {
	claims, err := f.factory.ParseToken(tokenString)
	if err != nil {
		return nil, newSrvError(ErrTokenGenerationFailure, err.Error())
	}
	return claims, nil
}
