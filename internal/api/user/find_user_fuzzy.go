// 根据用户昵称(nickname)模糊查询用户, 使用封装好的UserProfileResponse函数, 统一调用内部FindByFieldExact函数执行数据库查询
package user

import (
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"
	"net/url"

	"github.com/gin-gonic/gin"
)

// 查询数据库中用户总数+Redis中的在线人数接口
func FinUser(c *gin.Context) {
	// 注意需要前端进行url编码后才能请求, 这里是后端所以只需处理url路径的中文解码
	paramNickname, err := url.QueryUnescape(c.Param("param"))
	if err != nil {
		util.Rsp(c, 400, "参数解码失败, "+err.Error())
		return
	}
	if paramNickname == "" {
		util.Rsp(c, 400, "查询参数不能为空")
		return
	}
	// 调用封装的通用函数防止返回敏感信息, 传入查询方式(模糊fuzzy或精确exact)并指定要查询的数据库字段和查询参数
	service.UserProfileResponse(c, "fuzzy", "nickname", paramNickname)
}
