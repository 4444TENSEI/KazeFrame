// 临时状态操作
package cache

import (
	"KazeFrame/internal/config"
	"context"
	"fmt"
	"time"
)

// 通过刷新访问令牌时, 顺带记录用户在线状态到Redis, 方便进行用户在线数量统计
func SetUserOnline(id string) error {
	onlineKey := fmt.Sprintf("%s%s", UserOnlineKey, id)
	ctx := context.Background()
	if err := config.GetRedis().Set(ctx, onlineKey, "online", 0).Err(); err != nil {
		return err
	}
	timestamp := time.Now().Unix()
	if err := config.GetRedis().HSet(ctx, UserLastSeenKey, id, timestamp).Err(); err != nil {
		return err
	}
	return nil
}
