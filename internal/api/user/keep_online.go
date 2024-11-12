package user

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/internal/service"
	"time"

	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 刷新访问令牌, 维持登陆和在线状态
func KeepOnline(c *gin.Context) {
	// 从配置载入令牌的密钥和有效期限
	jwtKey := config.GetConfig().Token.JwtKey
	accessExp := config.GetConfig().Token.AccessExp
	refreshExp := config.GetConfig().Token.RefreshExp
	// 从cookie中获取刷新令牌
	userRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		util.Rsp(c, 401, 4101)
		service.ClearToken(c)
		return
	}
	// 验证刷新令牌
	checkedToken, err := util.VerifyToken(jwtKey, userRefreshToken)
	if err != nil {
		util.Rsp(c, 403, 4103)
		service.ClearToken(c)
		return
	}
	// 从配置文件读取时间进行创建
	var accessExpFormatted, refreshExpFormatted time.Duration
	accessExpFormatted = time.Duration(accessExp) * time.Second
	refreshExpFormatted = time.Duration(refreshExp) * time.Second
	// 生成新的访问令牌, 这里的下划线是刷新令牌, 不让它自己刷自己
	newAccessToken, _, err := util.CreateToken(jwtKey, accessExpFormatted, refreshExpFormatted, checkedToken.TkUID, checkedToken.TkUsername, checkedToken.TkRoleLevel)
	if err != nil {
		util.Rsp(c, 500, "生成新的访问令牌出错")
		return
	}
	// 将新的访问令牌设置到cookie中
	// 是否需要同时续期刷新令牌根据自身需求吧, 这里只刷新登陆令牌
	// c.SetCookie("refresh_token", *newRefreshToken, cfg.Token.RefreshExp, "/", "", secure, httpOnly)
	c.SetCookie("access_token", *newAccessToken, accessExp, "/", "", false, false)
	// 更新Redis中的用户在线状态
	if err := cache.SetUserOnline(checkedToken.TkUID); err != nil {
		util.Rsp(c, 500, "更新在线状态出错")
		return
	}
	util.Rsp(c, 200, "欢迎光临")
}
