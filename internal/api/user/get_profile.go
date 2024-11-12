// 登录后请求此接口以获取个人基本信息, 存入客户端方便给客户端使用
// 假如用户后期更新个人资料过后, 客户端只需请求此资料接口获取最新的个人资料即可
package user

import (
	"KazeFrame/internal/config"
	"KazeFrame/internal/model"
	"KazeFrame/internal/service"

	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 获取个人资料接口
func GetProfile(c *gin.Context) {
	// 通过IsExistingData函数存入用户数据库资料到内存中, 方便后续查找
	var userData *model.User
	// 校验客户端token并从中获取用户ID用来查找个人资料
	jwtKey := config.GetConfig().Token.JwtKey
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		util.Rsp(c, 401, 4101)
		service.ClearToken(c)
		return
	}
	checkedToken, err := util.VerifyToken(jwtKey, accessToken)
	if err != nil {
		util.Rsp(c, 403, 4103)
		service.ClearToken(c)
		return
	}
	// 使用验证过的Token中的ID获取用户在数据表中的可公开个人信息存入Cookie方便前端调用
	userData = service.IsExistingData("id", checkedToken.TkUID)
	httpOnly, secure := true, c.Request.TLS != nil
	c.SetCookie("ck_nickname", userData.Nickname, 0, "/", "", secure, httpOnly)
	c.SetCookie("ck_signature", userData.Signature, 0, "/", "", secure, httpOnly)
	c.SetCookie("ck_avatar_url", userData.AvatarUrl, 0, "/", "", secure, httpOnly)
	c.SetCookie("ck_background_url", userData.BackgroundUrl, 0, "/", "", secure, httpOnly)
	// 调用封装的通用函数防止返回敏感信息, 传入查询方式(模糊fuzzy或精确exact)并指定要查询的数据库字段和查询参数
	service.UserProfileResponse(c, "exact", "id", checkedToken.TkUID)
}
