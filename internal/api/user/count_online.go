package user

import (
	"KazeFrame/internal/cache"
	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/pkg/util"
	"context"

	"github.com/gin-gonic/gin"
)

// 查询数据库中用户总数+Redis中的在线人数接口
func GetUserOnlineCount(c *gin.Context) {
	//查询数据库用户总数
	totalCount, err := dao.UserRepo.CountTableData("", "")
	if err != nil {
		util.Rsp(c, 500, "操作失败"+err.Error())
		return
	}
	// 查询Redis在线用户人数
	ctx := context.Background()
	onlineCount, err := config.GetRedis().Keys(ctx, cache.UserOnlineKey+"*").Result()
	if err != nil {
		util.Rsp(c, 500, "操作失败"+err.Error())
		return
	}
	c.JSON(200, gin.H{"total_user_count": totalCount, "online_user_count": len(onlineCount)})
}
