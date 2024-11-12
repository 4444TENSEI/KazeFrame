package monitor

import (
	"KazeFrame/internal/config"
	"KazeFrame/internal/model"
	"KazeFrame/pkg/util"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 清理Redis缓存接口-根据键名
func ClearCacheByKey(c *gin.Context) {
	var clearPayload model.ClearCachePayload
	if err := c.ShouldBindJSON(&clearPayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	ctx := context.Background()
	var deletedKeys []string
	var notFoundKeys []string
	var failToDeleteKeys []string
	for _, key := range clearPayload.CacheKey {
		pattern := fmt.Sprintf("%s*", key)
		iter := config.GetRedis().Scan(ctx, 0, pattern, 0).Iterator()
		var keyFound bool
		for iter.Next(ctx) {
			keyFound = true
			deletedKey := iter.Val()
			if _, err := config.GetRedis().Del(ctx, deletedKey).Result(); err != nil {
				failToDeleteKeys = append(failToDeleteKeys, deletedKey)
			} else {
				deletedKeys = append(deletedKeys, deletedKey)
			}
		}
		if err := iter.Err(); err != nil {
			util.Rsp(c, 500, fmt.Sprintf("清除缓存键%s失败", key))
			return
		}
		if !keyFound {
			notFoundKeys = append(notFoundKeys, key)
		}
	}
	if len(deletedKeys) == 0 {
		util.Rsp(c, 404, "未找到要删除的缓存键名")
		return
	}
	c.JSON(200, gin.H{
		"deleted_key": deletedKeys,
		"not_key":     notFoundKeys,
		"fail_key":    failToDeleteKeys,
	})
}

// 清理Redis全部缓存接口
func ClearAllCache(c *gin.Context) {
	ctx := context.Background()
	iter := config.GetRedis().Scan(ctx, 0, "", 0).Iterator()
	for iter.Next(ctx) {
		config.GetRedis().Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		util.Rsp(c, 500, "清除所有缓存失败")
		return
	}
	util.Rsp(c, 200, "所有缓存已清除")
}
