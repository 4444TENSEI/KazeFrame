package user

import (
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 根据请求体中的字段名的参数值, 批量删除用户接口
func DeleteUser(c *gin.Context) {
	var userDeletePayload model.DeletePayload
	if err := c.ShouldBindJSON(&userDeletePayload); err != nil {
		util.Rsp(c, 400, "请求参数错误: "+err.Error())
		return
	}
	response, err := dao.UserRepo.QuickHardDelete(userDeletePayload.Field, userDeletePayload.Value)
	if err != nil {
		util.Rsp(c, 500, "删除操作失败: "+err.Error())
		return
	}
	if response.OkCount == 0 {
		util.Rsp(c, 404, "没找到你想删除的数据捏")
		return
	}
	c.JSON(200, response)
}
