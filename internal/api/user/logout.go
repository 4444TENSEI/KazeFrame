package user

import (
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 退出登录接口
func UserLogout(c *gin.Context) {
	// 从Cookie取ID, 设置redis用户状态为离线
	ckUid, err := c.Cookie("ck_uid")
	if err != nil {
		util.Rsp(c, 401, "未找到用户ID")
		return
	}
	if err := service.SetUserOffline(ckUid); err != nil {
		util.Rsp(c, 500, "设置Redis离线状态失败")
		return
	}
	// 遍历浏览器产生的站内Cookie数据并清空
	for _, cookie := range c.Request.Cookies() {
		c.SetCookie(cookie.Name, "", -1, "/", "", false, true)
	}
	util.Rsp(c, 200, "退出登录成功, 所有本地数据已经清空")
}
