package cache

import (
	"github.com/redis/go-redis/v9"
	"nba-remake/internal/config"
)

func NewCache(redisConf *config.RedisConfig) *redis.Client {
	// 创建 Redis 客户端
	return redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})
}
