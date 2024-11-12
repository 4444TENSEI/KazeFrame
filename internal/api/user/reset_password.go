package user

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"context"
	"fmt"

	"KazeFrame/internal/model"

	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 用户找回密码接口
func UserForgetPassword(c *gin.Context) {
	// 校验配置文件中的邮件服务开关"enable"是否开启
	if !config.GetConfig().Email.Enable {
		util.Rsp(c, 400, 4502)
		return
	}
	var forgetPayload model.UserForgetPayload
	if err := c.ShouldBindJSON(&forgetPayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	// 检查邮箱是否已存在
	if _, err := dao.UserRepo.FindByFieldExact("email", forgetPayload.Email); err != nil {
		util.Rsp(c, 422, 4602)
		return
	}
	// 验证邮箱验证码
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", cache.EmailForgetCaptchaKey, forgetPayload.Email)
	code, err := config.GetRedis().Get(ctx, key).Result()
	if err == redis.Nil {
		util.Rsp(c, 400, 4401)
		return
	}
	// 校验Redis中的验证码
	if err != nil || code != forgetPayload.ForgetPswCaptcha {
		util.Rsp(c, 400, 4400)
		return
	}
	// 验证码正确, 进行密码加密并更新到数据库
	hashedPassword, err := util.BcryptPassword(forgetPayload.NewPassword)
	if err != nil {
		util.Rsp(c, 500, "密码加密错误")
		return
	}
	newPassword := map[string]interface{}{"password": hashedPassword}
	// 更新用户密码
	if err := dao.UserRepo.UpdateByField("email", forgetPayload.Email, newPassword); err != nil {
		util.Rsp(c, 500, "密码更新失败")
		return
	}
	// 密码更新成功后, 删除Redis中的旧验证码
	if err := config.GetRedis().Del(ctx, key).Err(); err != nil {
		util.Rsp(c, 500, "删除验证码失败")
		return
	}
	util.Rsp(c, 200, "密码找回成功, 请使用新密码登录")
}
