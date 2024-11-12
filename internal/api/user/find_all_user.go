package user

import (
	"KazeFrame/internal/dao"
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有用户信息、在线状态接口
func GetAllUser(c *gin.Context) {
	pageStr, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSizeStr, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	allUserData, dataCount, err := dao.UserRepo.FindTableData(pageStr, pageSizeStr)
	if err != nil {
		util.Rsp(c, 500, "数据查询失败, "+err.Error())
		return
	}
	// 获取Redis中的各个用户在线状态
	// 创建一个新的切片来存储只包含所需字段的用户数据
	usersResponse := make([]map[string]interface{}, 0, len(allUserData))
	for _, userData := range allUserData {
		// 获取Redis中的各个用户在线状态
		isOnline, _ := service.IsUserOnline(strconv.Itoa(int(userData.ID)))
		userMap := map[string]interface{}{
			"uid":        userData.ID,
			"username":   userData.Username,
			"email":      userData.Email,
			"nickname":   userData.Nickname,
			"signature":  userData.Signature,
			"gender":     userData.Gender,
			"avatar_url": userData.AvatarUrl,
			"role_level": userData.RoleLevel,
			"online":     strconv.FormatBool(isOnline),
			"CreatedAt":  userData.CreatedAt,
			"UpdatedAt":  userData.UpdatedAt,
		}
		usersResponse = append(usersResponse, userMap)
	}
	c.JSON(200, gin.H{
		"count": dataCount,
		"data":  usersResponse,
	})
}
