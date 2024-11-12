// 负责整体的项目初始化, 在此检查数据库、Redis连通性、日志记录器等
package config

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	globalConfig *Config
	globalDB     *gorm.DB
	globalRedis  *redis.Client
	globalLogger *zap.SugaredLogger
)

// 进行服务运行必要的初始化, 全局公开日志记录器、配置载入、数据库连接、Redis连接
func InitServer() error {
	// 全局公开zap日志记录器
	globalLogger = InitLogger()
	// 加载配置文件
	if err := InitConfig(); err != nil {
		globalLogger.Error("XXX ", err)
		return err
	}
	// 初始化数据库连接
	if err := InitDB(); err != nil {
		globalLogger.Error("XXX ", err)
		return err
	}
	// 初始化Redis连接
	if err := InitRedis(); err != nil {
		globalLogger.Error("XXX ", err)
		return err
	}
	return nil
}

// 获取配置信息
func GetConfig() *Config {
	return globalConfig
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return globalDB
}

// 获取Redis实例
func GetRedis() *redis.Client {
	return globalRedis
}

// 获取日志记录器实例
func GetLogger() *zap.SugaredLogger {
	return globalLogger
}
