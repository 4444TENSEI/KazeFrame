// 站点访问日志中间件-录入到数据库表
package middleware

import (
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/pkg/util"
	"time"

	"KazeFrame/internal/model"

	"github.com/gin-gonic/gin"
)

// 接口访问日志中间件, 同时还能够捕捉运行时的数据库错误使用zap记录到本地日志文件
// 将会录入到数据库表的内容查看model.RequestLog结构体
func LoggerMiddleware() gin.HandlerFunc {
	// 从配置载入令牌密钥
	jwtKey := config.GetConfig().Token.JwtKey
	logger := config.GetLogger()
	return func(c *gin.Context) {
		var logUid, logUsername, logRoleLevel string
		// 如果用户未登录则不进行Token校验, 直接标记为未登录
		// 如果是已登录的用户, 每一次请求服务都进行Token校验, 略微造成性能影响, 验证为了防止伪造
		userAccessToken, err := c.Cookie("access_token")
		if err != nil || userAccessToken == "" {
			logUid = ""
			logUsername = "未登录"
			logRoleLevel = ""
		} else {
			checkedToken, err := util.VerifyToken(jwtKey, userAccessToken)
			if err != nil {
				logUid = ""
				logUsername = "未登录"
				logRoleLevel = ""
			}
			// Token校验通过, 取Token中的用户ID、昵称、权限等级录入到日志表
			logUid = checkedToken.TkUID
			logUsername = checkedToken.TkUsername
			logRoleLevel = checkedToken.TkRoleLevel
		}
		requestMethod := c.Request.Method
		requestRoute := c.Request.URL.Path
		c.Next()
		responeCode := c.Writer.Status()
		userLog := model.RequestLog{
			UID:           logUid,
			Username:      logUsername,
			RoleLevel:     logRoleLevel,
			RequestIP:     c.ClientIP(),
			RequestTime:   time.Now(),
			RequestMethod: requestMethod,
			RequestRoute:  requestRoute,
			ResponeCode:   responeCode,
		}
		if err := dao.RequestLogRepo.Create(userLog); err != nil {
			// 全局捕捉数据库服务运行时错误
			logger.Errorf("服务端数据库炸了, %v", err)
			//  此时还可以顺便推送服务器运行错误的消息到你的webhook消息推送服务接口, 例如pushplus
			//  http.Get("你的pushplus webhook地址")
			return
		}
	}
}
