package user

import (
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/internal/service"

	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 更新个人资料
func UpdateProfile(c *gin.Context) {
	// 从配置载入令牌密钥
	jwtKey := config.GetConfig().Token.JwtKey
	var updateProfilePayload model.UserUpdatePayload
	// 验证用户Cookie中的登录令牌
	userAccessToken, err := c.Cookie("access_token")
	if err != nil || userAccessToken == "" {
		util.Rsp(c, 401, 4101)
		service.ClearToken(c)
		return
	}
	checkedToken, err := util.VerifyToken(jwtKey, userAccessToken)
	if err != nil {
		util.Rsp(c, 403, 4103)
		service.ClearToken(c)
		return
	}
	if err := c.ShouldBindJSON(&updateProfilePayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	updates := make(map[string]interface{})
	if updateProfilePayload.Nickname != "" {
		updates["nickname"] = updateProfilePayload.Nickname
	}
	if updateProfilePayload.AvatarUrl != "" {
		updates["avtar_url"] = updateProfilePayload.AvatarUrl
	}
	if updateProfilePayload.BackgroundUrl != "" {
		updates["background_url"] = updateProfilePayload.BackgroundUrl
	}
	if updateProfilePayload.Signature != "" {
		updates["signature"] = updateProfilePayload.Signature
	}
	if updateProfilePayload.Gender != "" {
		updates["gender"] = updateProfilePayload.Gender
	}
	if len(updates) == 0 {
		util.Rsp(c, 400, "个人资料没有任何更新")
		return
	}
	// 使用校验过后的Token中的UID更新个人资料
	if err := dao.UserRepo.UpdateByField("id", checkedToken.TkUID, updates); err != nil {
		util.Rsp(c, 500, "个人资料更新失败"+err.Error())
		return
	}
	util.Rsp(c, 200, "个人资料更新成功")
}
