// 跨域中间件
package middleware

import (
	"KazeFrame/internal/config"

	"github.com/gin-gonic/gin"
)

// 跨域中间件, 允许的请求源列表在config.yml定义
func CORS() gin.HandlerFunc {
	corsConfig := config.GetConfig().CORS
	allowedOrigins := make(map[string]bool)
	for _, origin := range corsConfig {
		allowedOrigins[origin] = true
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins[origin] || allowedOrigins["*"] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(200)
			}
		}
		c.Next()
	}
}
