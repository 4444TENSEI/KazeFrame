package user

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 用户注册接口
func UserRegisterPayload(c *gin.Context) {
	var registerPayload model.UserRegisterPayload
	if err := c.ShouldBindJSON(&registerPayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	if registerPayload.Email == "" && registerPayload.Username == "" {
		util.Rsp(c, 400, "邮箱和用户名不能同时为空")
		return
	}
	// 检查用户名是否已被注册
	existingUsername := service.IsExistingData("username", registerPayload.Username)
	if existingUsername != nil {
		util.Rsp(c, 409, "用户名已被注册, 请直接登录或找回密码")
		return
	}
	// 检查邮箱是否已存在
	existingEmail := service.IsExistingData("email", registerPayload.Email)
	if existingEmail != nil {
		util.Rsp(c, 409, "邮箱已被注册, 请直接登录或找回密码")
		return
	}

	// 根据配置判断是否需要校验邮箱验证码
	if config.GetConfig().Email.Enable {
		// 验证邮箱注册验证码
		ctx := context.Background()
		key := fmt.Sprintf("%s%s", cache.EmailRegisterCaptchaKey, registerPayload.Email)
		code, err := config.GetRedis().Get(ctx, key).Result()
		if err == redis.Nil {
			util.Rsp(c, 400, 4401)
			return
		}
		if err != nil {
			util.Rsp(c, 500, "获取验证码失败")
			return
		}
		if code != registerPayload.RegisterCaptcha {
			util.Rsp(c, 400, "验证码错误")
			return
		}
		// 校验通过后删除redis中的邮箱验证码
		if err := config.GetRedis().Del(ctx, key).Err(); err != nil {
			util.Rsp(c, 500, "删除验证码失败")
			return
		}
	}

	// 校验通过创建用户
	hashedPassword, err := util.BcryptPassword(registerPayload.Password)
	if err != nil {
		util.Rsp(c, 500, "密码加密错误")
		return
	}
	// 构造最终注册录入的用户表数据
	userRegisterData := model.User{
		Nickname: registerPayload.Nickname,
		Username: registerPayload.Username,
		Email:    registerPayload.Email,
		Password: hashedPassword,
	}
	if err := dao.UserRepo.Create(userRegisterData); err != nil {
		util.Rsp(c, 500, "用户创建失败")
		return
	}
	util.Rsp(c, 200, "用户注册成功")
}
