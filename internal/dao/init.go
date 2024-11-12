// 数据库操作的相关结构和函数

package dao

import (
	"KazeFrame/internal/model"

	"gorm.io/gorm"
)

// 用于存储不同表的数据库操作实例的全局变量
// 如果你在intel/model/table_basic新增了数据库表，请在这里以及InitRepo()对应新增的数据库操作实例
// 即可使用通用的数据库CRUD操作例如dao.UserRepo.CountTableData("", "")
var (
	UserRepo       *customRepo[model.User]
	RequestLogRepo *customRepo[model.RequestLog]
	// 后期新增的数据库操作实例...
)

// 初始化数据库仓库，为每个表的连接不同的数据库操作实例
func InitRepo(db *gorm.DB) {
	UserRepo = NewRepo[model.User](db)
	RequestLogRepo = NewRepo[model.RequestLog](db)
	// 后期新增的数据库表的数据库操作实例...
}

// T代表具体的数据库模型, 可实现数据库表的通用操作接口结构
type customRepo[T any] struct {
	DB *gorm.DB
}

// 创建一个新的数据库操作实例, T代表具体的数据库模型
func NewRepo[T any](db *gorm.DB) *customRepo[T] {
	return &customRepo[T]{DB: db}
}
