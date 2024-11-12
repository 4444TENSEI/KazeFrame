// 用户业务通用操作逻辑/函数
package service

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/pkg/util"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 判断是否存在的数据
// 如果存在则返回第一条数据的指针(用于取出数据库信息, 可以用于登录比较密码)
// 如果不存在则返回个空指针提示需要注册
func IsExistingData(fieldName string, value string) *model.User {
	users, err := dao.UserRepo.FindByFieldExact(fieldName, value)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}
	if len(users) > 0 {
		return &users[0]
	}
	return nil
}

// 通过在Redis中查找用户ID的来检查用户是否在线
func IsUserOnline(id string) (bool, error) {
	ctx := context.Background()
	onlineStatus, err := config.GetRedis().Get(ctx, fmt.Sprintf("%s%s", cache.UserOnlineKey, id)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return onlineStatus == "online", nil
}

// 根据用户名删除Redis中的在线状态
func SetUserOffline(id string) error {
	onlineKey := fmt.Sprintf("%s%s", cache.UserOnlineKey, id)
	ctx := context.Background()
	if err := config.GetRedis().Del(ctx, onlineKey).Err(); err != nil {
		return err
	}
	return nil
}

// 便捷清空登陆令牌和刷新令牌, 用于登录失效时
func ClearToken(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.Abort()
}

// 统一构造用户资料响应, 防止直接返回例如密码等重要、敏感信息在响应内容中
// 需要指定查询是精准查询还是模糊查询, 以及字段名以及对应要查询的内容
func UserProfileResponse(c *gin.Context, searchType string, field, value string) {
	var userData *model.User
	// 精准查询
	if searchType == "exact" {
		users, err := dao.UserRepo.FindByFieldExact(field, value)
		if err != nil && err != gorm.ErrRecordNotFound {
			util.Rsp(c, 500, "查询用户数据时发生错误"+err.Error())
			return
		}
		userData = &users[0]
	}
	// 模糊查询
	if searchType == "fuzzy" {
		users, err := dao.UserRepo.FindByFieldFuzzy(field, value)
		if err != nil && err != gorm.ErrRecordNotFound {
			util.Rsp(c, 500, "查询用户数据时发生错误"+err.Error())
			return
		}
		userData = &users[0]
	}
	// 从Redis数据判断用户是否在线
	online, err := IsUserOnline(value)
	if err != nil {
		util.Rsp(c, 500, "检查用户在线状态时发生错误"+err.Error())
		return
	}
	response := gin.H{
		"uidaa":      userData.ID,
		"username":   userData.Username,
		"email":      userData.Email,
		"nickname":   userData.Nickname,
		"signature":  userData.Signature,
		"gender":     userData.Gender,
		"avatar_url": userData.AvatarUrl,
		"role_level": userData.RoleLevel,
		"online":     online,
		"CreatedAt":  userData.CreatedAt,
		"UpdatedAt":  userData.UpdatedAt,
	}
	c.JSON(200, response)
}
