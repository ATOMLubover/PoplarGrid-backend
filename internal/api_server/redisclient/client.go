package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 接口定义了 Redis 客户端的使用方法
type RedisClient interface {
}

// redisClientImpl 是 RedisClient 的实现
type redisClientImpl struct {
	client *redis.Client
}

// NewRedisClient 创建一个新的 RedisClient 实例
func NewRedisClient(
	addr string,
	password string,
	maxConn int,
	minIdle int,
	maxRetries int,
) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		PoolSize:     maxConn,
		MinIdleConns: minIdle,
		MaxRetries:   maxRetries,
		PoolTimeout:  0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 此处 ping 一次服务器，保证连接正常
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &redisClientImpl{
		client: client,
	}, nil
}
