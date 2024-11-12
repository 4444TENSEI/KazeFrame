// 在这里定义好的结构体用于Gorm自动建表，使用方式是：
// 按照Gorm规范建立好新的结构体，将结构体名称首字母设置为大写
// 全局搜索"AutoMigrate"，在内部增加你的结构名称，即可在项目运行时自动创建表
package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户表user
type User struct {
	ID            uint   `gorm:"primarykey;" json:"uid" comment:"自增用户UID"`
	Username      string `gorm:"type:varchar(255);not null;uniqueIndex;" json:"username" comment:"用户账号"`
	Email         string `gorm:"type:varchar(255);not null;uniqueIndex;" json:"email" comment:"邮箱"`
	Password      string `gorm:"type:varchar(255);not null;" json:"password" comment:"密码"`
	Nickname      string `gorm:"type:varchar(255);default:'默认用户名';index;" json:"nickname" comment:"昵称"`
	Signature     string `gorm:"type:varchar(255);default:'还没有个性签名~';" json:"signature" comment:"个性签名"`
	Gender        string `gorm:"type:varchar(255);default:null;" json:"gender" comment:"性别"`
	AvatarUrl     string `gorm:"type:varchar(255);default:null;" json:"avatar_url" comment:"头像图片在线地址"`
	BackgroundUrl string `gorm:"type:varchar(255);default:null;" json:"background_url" comment:"背景图片在线地址"`
	RoleLevel     int    `gorm:"size:1;default:2;" json:"role_level" comment:"注册默认为用户角色, 1:游客, 2:用户, 3:管理员"`
	gorm.Model
}

// 请求日志表request_log
// 这里的UID和RoleLevel设置为string类型是为了方便从Cookie中直接读取无需转换
type RequestLog struct {
	ID            uint      `gorm:"primarykey;" json:"log_id" comment:"日志序号和用户无关"`
	RequestTime   time.Time `gorm:"type:datetime;default:null;" json:"request_time" comment:"操作时间"`
	RequestIP     string    `gorm:"type:varchar(255);default:null;" json:"request_ip" comment:"操作IP"`
	RequestRoute  string    `gorm:"type:varchar(255);" json:"request_route" comment:"请求接口"`
	RequestMethod string    `gorm:"type:varchar(255);" json:"request_method" comment:"请求方法"`
	ResponeCode   int       `json:"respone_code" comment:"响应头状态码"`
	UID           string    `gorm:"index;" json:"uid" comment:"用户ID"`
	Username      string    `gorm:"type:varchar(255);" json:"nickname" comment:"用户账号"`
	RoleLevel     string    `json:"role_level" comment:"权限等级"`
}
