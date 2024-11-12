package user

import (
	"KazeFrame/internal/config"
	"KazeFrame/internal/model"
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 用户登录接口
func UserLogin(c *gin.Context) {
	var userData *model.User
	var loginPayload model.UserLoginPayload
	var err error
	// 从配置载入令牌的密钥和有效期限
	jwtKey := config.GetConfig().Token.JwtKey
	accessExp := config.GetConfig().Token.AccessExp
	refreshExp := config.GetConfig().Token.RefreshExp
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		// 使用gin自带的validator校验（这里已经在model结构体定义了）
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	if loginPayload.Email == "" && loginPayload.Username == "" {
		util.Rsp(c, 400, "邮箱和用户名不能同时为空")
		return
	}
	if loginPayload.Email != "" && loginPayload.Username != "" {
		util.Rsp(c, 400, "只能单独使用邮箱或用户名登录, 不能同时使用")
		return
	}
	// 校验用户是否已注册, 未注册无法登录, 如果存在则需要取出用户数据校验登陆密码
	if loginPayload.Email != "" {
		userData = service.IsExistingData("email", loginPayload.Email)
	} else if loginPayload.Username != "" {
		userData = service.IsExistingData("username", loginPayload.Username)
	}
	if userData == nil {
		util.Rsp(c, 404, "用户不存在, 先去注册叭")
		return
	}
	// 校验密码
	if !util.ComparePassword(loginPayload.Password, userData.Password) {
		util.Rsp(c, 401, "登录失败, 密码错误")
		return
	}
	// 创建Token存到Cookie, 内容包括用户ID、用户Username、用户等级RoleLevel、过期时间
	// 如果用户勾选了“三十天内免登录”, 则刷新令牌设置为30天, 否则使用配置文件中的默认值
	var accessExpFormatted, refreshExpFormatted time.Duration
	if loginPayload.RememberMe {
		refreshExpFormatted = 30 * 24 * time.Hour
	} else {
		refreshExpFormatted = time.Duration(refreshExp) * time.Second
	}
	accessExpFormatted = time.Duration(accessExp) * time.Second
	userIDStr := strconv.FormatUint(uint64(userData.ID), 10)
	roleLevelStr := strconv.Itoa(userData.RoleLevel)
	accessToken, refreshToken, err := util.CreateToken(jwtKey, accessExpFormatted, refreshExpFormatted, userIDStr, userData.Username, roleLevelStr)
	if err != nil {
		util.Rsp(c, 500, "登录失败, Token创建失败")
		return
	}
	// 登陆时设置登录Token、刷新Token到Cookie中
	httpOnly, secure := true, c.Request.TLS != nil
	c.SetCookie("ck_uid", userIDStr, int(refreshExpFormatted.Seconds()), "/", "", secure, httpOnly)
	c.SetCookie("access_token", *accessToken, accessExp, "/", "", secure, httpOnly)
	c.SetCookie("refresh_token", *refreshToken, int(refreshExpFormatted.Seconds()), "/", "", secure, httpOnly)
	util.Rsp(c, 200, "登录成功")
}
