package monitor

import (
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/pkg/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 根据请求体中的时间字段, 批量删除日志表数据接口
func ClearRequestLogBytime(c *gin.Context) {
	var logDeletePayload model.DeletePayload
	if err := c.ShouldBindJSON(&logDeletePayload); err != nil {
		util.Rsp(c, 400, "请求参数错误: "+err.Error())
		return
	}
	beforeTimeStr := logDeletePayload.Value[0].(string)
	beforeTime, err := time.Parse(time.RFC3339, beforeTimeStr)
	if err != nil {
		util.Rsp(c, 400, "时间格式错误, 请参考: 2077-01-02T15:04:05Z, "+err.Error())
		return
	}
	deletedCount, err := dao.RequestLogRepo.DeleteByTime(logDeletePayload.Field, beforeTime)
	if err != nil {
		util.Rsp(c, 500, "删除日志过程出错: "+err.Error())
		return
	}
	if deletedCount == 0 {
		util.Rsp(c, 404, "没找到你想删除的数据")
		return
	}
	c.JSON(200, gin.H{"deleted_before": beforeTimeStr, "deleted_count": deletedCount})
}

// 清空日志表数据
func ClearAllRequestLog(c *gin.Context) {
	// 这里ClearAllData()传一个true代表硬删除(无法恢复的永久删除), false代表软删除(实际保留数据)
	// 但如果数据库自动迁移的结构体没有添加gorm.model或是DeletedAt字段, 那么这里无论如何都是硬删除
	deletedCount, err := dao.RequestLogRepo.ClearAllData(false)
	if err != nil {
		util.Rsp(c, 500, "删除日志过程出错: "+err.Error())
		return
	}
	if deletedCount == 0 {
		util.Rsp(c, 404, "日志数据表中没有任何数据")
		return
	}
	deletedCountStr := strconv.FormatInt(deletedCount, 10)
	util.Rsp(c, 200, "日志数据已全部清空，共删除["+deletedCountStr+"]条记录")
}

// 批量删除日志表数据接口-自定义条件, 例如删除id为1,2,3或者name为张三,李四的日志
// 暂时不用, 先留个示例
func ClearRequestLog(c *gin.Context) {
	var logDeletePayload model.DeletePayload
	if err := c.ShouldBindJSON(&logDeletePayload); err != nil {
		util.Rsp(c, 400, "请求参数错误: "+err.Error())
		return
	}
	response, err := dao.RequestLogRepo.QuickHardDelete(logDeletePayload.Field, logDeletePayload.Value)
	if err != nil {
		util.Rsp(c, 500, "删除操作失败: "+err.Error())
		return
	}
	if response.OkCount == 0 {
		util.Rsp(c, 404, "没有找到你想删除的捏")
		return
	}
	c.JSON(200, response)
}
