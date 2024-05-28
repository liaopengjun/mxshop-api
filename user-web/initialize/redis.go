package initialize

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mxshop-api/user-web/global"

	"go.uber.org/zap"
)

// Redis 初始化redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       global.ServerConfig.RedisInfo.DB,
		PoolSize: global.ServerConfig.RedisInfo.PoolSize,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.Redis = client
	}
	return
}
