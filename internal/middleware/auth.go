// 权限认证中间件, 校验Token中权限等级
package middleware

import (
	"strconv"

	"KazeFrame/internal/config"
	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 权限认证中间件, 使用JWT校验登陆令牌中的权限等级
func RoleAuth(requiredRoleLevel int) gin.HandlerFunc {
	jwtKey := config.GetConfig().Token.JwtKey
	return func(c *gin.Context) {
		// 判断请求者是否携带登录令牌
		tokenString, err := c.Cookie("refresh_token")
		if err != nil || tokenString == "" {
			util.Rsp(c, 401, 4104)
			c.Abort()
			return
		}
		// 校验登陆令牌中的权限等级
		checkedToken, err := util.VerifyToken(jwtKey, tokenString)
		if err != nil {
			util.Rsp(c, 403, 4105)
			c.Abort()
			return
		}
		intRoleLevel, _ := strconv.Atoi(checkedToken.TkRoleLevel)
		if intRoleLevel < requiredRoleLevel {
			util.Rsp(c, 403, 4102)
			c.Abort()
			return
		}
		c.Next()
	}
}
