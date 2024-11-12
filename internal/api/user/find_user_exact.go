// 根据ID精准查询用户, 使用封装好的UserProfileResponse函数, 统一调用内部FindByFieldExact函数执行数据库查询
package user

import (
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 通过URL中的用户ID精准获取用户资料接口
func GetUserByID(c *gin.Context) {
	paramID := c.Param("param")
	if paramID == "" {
		util.Rsp(c, 400, "用户ID不能为空")
		return
	}
	// 调用封装的通用函数防止返回敏感信息, 传入查询方式(模糊fuzzy或精确exact)并指定要查询的数据库字段和查询参数
	service.UserProfileResponse(c, "exact", "id", paramID)
}
