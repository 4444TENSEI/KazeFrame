// 初始化数据库连接实例, 以及gorm自动迁移自动建表
package config

import (
	"KazeFrame/internal/model"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据库迁移列表, 在AutoMigrate函数内添加要自动建表的结构体, 服务初始化时会自动建表
func MigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.User{}, &model.RequestLog{}); err != nil {
		return errors.Wrap(err, "执行数据库迁移失败")
	}
	return nil
}

// 初始化数据库表和进行数据库自动迁移建表
func InitDB() error {
	var err error
	dsn := GetConfig().Database.DSN
	globalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			// 数据库表前缀
			// TablePrefix: "kz_",
		},
	})
	if err != nil {
		return errors.Wrap(err, "数据库连接失败, 请检查配置文件中的数据库连接信息或数据库运行状态")
	}
	if err := MigrateDB(globalDB); err != nil {
		return err
	}
	return nil
}
