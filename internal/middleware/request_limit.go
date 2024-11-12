// 频繁请求限制中间件
package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/pkg/util"
)

// 使用时传递2个参数: 请求上限次数+限制的时间段
// 这里顺便能捕捉到Redis服务运行错误, 通过zap自动记录到本地运行日志
// 使用示例(60分钟内同一IP最多请求200次): r.Use(middleware.IPRateLimiter(200, 60*time.Minute))
// 按照IP+访问次数的形式记录到Redis进行访问次数统计
func IPRateLimiter(max int, expires time.Duration) gin.HandlerFunc {
	logger := config.GetLogger()
	client := config.GetRedis()
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("%s%s:%s", cache.RequestCount, ip, c.Request.URL.Path)
		ctx := context.Background()
		allowed, err := client.SetNX(ctx, key, 0, expires).Result()
		if err != nil {
			util.Rsp(c, 500, "服务端缓存服务炸了, "+err.Error())
			// 记录Redis崩溃到本地日志文件
			logger.Error("服务端缓存服务炸了, ", err.Error())
			//  此时还可以顺便推送服务器运行错误的消息到你的webhook消息推送服务接口, 例如pushplus
			//  http.Get("你的pushplus webhook地址")
			return
		}
		if !allowed {
			current, err := client.Get(ctx, key).Int()
			if err != nil {
				if err == redis.Nil {
					c.Next()
				}
				util.Rsp(c, 500, "服务端缓存服务处理不过来, "+err.Error())
				return
			}
			if current >= max {
				util.Rsp(c, 429, fmt.Sprintf("请勿频繁操作, %d分钟内最多允许请求%d次", int(expires.Minutes()), max))
				return
			}
			if _, err := client.Incr(ctx, key).Result(); err != nil {
				util.Rsp(c, 429, "服务端缓存服务处理失败: "+err.Error())
				return
			}
		}
		c.Next()
	}
}
