package monitor

import (
	"KazeFrame/internal/dao"
	"KazeFrame/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 查询日志数据+日志数据量接口
func GetRequestLog(c *gin.Context) {
	pageStr, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	allUserData, dataCount, err := dao.RequestLogRepo.FindTableData(pageStr, pageSizeStr)
	if err != nil {
		util.Rsp(c, 500, "数据查询失败, "+err.Error())
		return
	}
	c.JSON(200, gin.H{
		"count": dataCount,
		"data":  allUserData,
	})
}
