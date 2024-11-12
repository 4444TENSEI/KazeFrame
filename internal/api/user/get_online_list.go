package user

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/pkg/util"
	"context"

	"github.com/gin-gonic/gin"
)

// 获取所有在线用户ID列表接口
func GetOnlineUserList(c *gin.Context) {
	ctx := context.Background()
	onlineUserKeys, err := config.GetRedis().Keys(ctx, cache.UserOnlineKey+"*").Result()
	if err != nil {
		util.Rsp(c, 500, "获取在线用户ID失败"+err.Error())
		return
	}
	onlineUsers := make([]string, 0)
	for _, key := range onlineUserKeys {
		userID := key[len(cache.UserOnlineKey):]
		onlineUsers = append(onlineUsers, userID)
	}
	if len(onlineUsers) == 0 {
		util.Rsp(c, 200, "当前无人在线~")
		return
	}
	c.JSON(200, onlineUsers)
}
