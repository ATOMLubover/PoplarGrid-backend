package jwtutils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 不对外暴露 JWT 的实现细节
type customClaims[T any] struct {
	jwt.RegisteredClaims
	Inner T
}

// JWT 的声明工厂
// T 是要封装的认证数据结构
type CliamsFactory[T any] struct {
	ExpirationTime time.Duration
	SecretKey      string
}

// 构造 JWT 声明工厂
func NewClaimsFactory[ClaimItem any](
	expiration time.Duration,
	secretKey string,
) *CliamsFactory[ClaimItem] {
	return &CliamsFactory[ClaimItem]{
		ExpirationTime: expiration,
		SecretKey:      secretKey,
	}
}

// 封装声明为签名令牌字符串
// 使用 HS256 算法签名
func (f *CliamsFactory[T]) GenToken(item T) (string, error) {
	now := time.Now()
	expireAt := now.Add(f.ExpirationTime)

	claims := customClaims[T]{
		Inner: item,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(f.SecretKey))
}

// 解析签名令牌
func (f *CliamsFactory[T]) ParseToken(signed string) (*T, error) {
	token, err := jwt.ParseWithClaims(
		signed,
		&customClaims[T]{},
		func(token *jwt.Token) (any, error) {
			return []byte(f.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*customClaims[T])
	if !ok || !token.Valid {
		return nil, errors.New("token无法解析")
	}

	return &claims.Inner, err
}
