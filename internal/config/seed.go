// 录入初始数据到数据库表(表中不存在任何记录时)
package config

import (
	"KazeFrame/internal/dao"
	"KazeFrame/internal/model"
	"KazeFrame/pkg/util"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// 如果项目启动时数据表没有数据, 则插入基础数据
func Seed(db *gorm.DB) error {
	userRepo := dao.UserRepo
	var userCount int64
	if err := userRepo.DB.Model(&model.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}
	if userCount == 0 {
		users := []*model.User{
			{ID: 100001, Nickname: "管理员", Username: "admin", Email: "admin@xxx.com", Password: "123456", RoleLevel: 3},
			{ID: 100002, Nickname: "用户喵", Username: "user", Email: "user@xxx.com", Password: "123456", RoleLevel: 2},
			{ID: 100003, Nickname: "游客喵", Username: "visitor", Email: "visitor@xxx.com", Password: "123456", RoleLevel: 1},
		}
		for _, user := range users {
			hashedPassword, err := util.BcryptPassword(user.Password)
			if err != nil {
				return fmt.Errorf("密码加密失败: %w", err)
			}
			user.Password = hashedPassword
			if err := userRepo.Create(*user); err != nil {
				return fmt.Errorf("初始数据录入失败: %w", err)
			}
			log.Printf("user表初始数据录入成功:\n%+v", user)
		}
	}
	return nil
}
