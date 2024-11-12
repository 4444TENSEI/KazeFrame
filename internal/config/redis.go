// Redis数据库配置
package config

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// 初始化Redis连接实例
func InitRedis() error {
	globalRedis = redis.NewClient(&redis.Options{
		Addr:     GetConfig().Redis.Address,
		Password: GetConfig().Redis.Password,
		DB:       GetConfig().Redis.Database,
	})
	ctx := context.Background()
	_, err := globalRedis.Ping(ctx).Result()
	if err != nil {
		return errors.Wrap(err, "Redis连接失败, 猜你没开")
	}
	return nil
}
