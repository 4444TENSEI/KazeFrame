// 配置文件加载路径、结构在此定义
// 项目打包后，配置文件放在static/server/config/config.yaml方便管理
// 因为还需要存放其他静态资源, 例如邮件模板、后期个人开发的前端页面等

package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 配置文件结构
type Config struct {
	Server struct {
		StaticPath string `mapstructure:"static_path"`
		Debug      bool   `mapstructure:"debug"`
		Port       string `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		DSN string `mapstructure:"dsn"`
	}
	Token struct {
		JwtKey     string `mapstructure:"jwt_key"`
		AccessExp  int    `mapstructure:"access_exp"`
		RefreshExp int    `mapstructure:"refresh_exp"`
	} `mapstructure:"token"`
	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"database"`
	} `mapstructure:"redis"`
	Email struct {
		Enable         bool   `mapstructure:"enable"`
		SenderName     string `mapstructure:"sender_name"`
		SenderEmail    string `mapstructure:"sender_email"`
		SenderPassword string `mapstructure:"sender_password"`
		SmtpServer     string `mapstructure:"smtp_server"`
		SmtpPort       int    `mapstructure:"smtp_port"`
	} `mapstructure:"email"`
	CORS []string `mapstructure:"cors"`
}

// 载入配置文件进行初始化
func InitConfig() error {
	var cfg Config
	globalConfig = &cfg
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./static/server/config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("加载配置文件失败: %w", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	// 热重载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Printf("配置文件重载失败: %s", err)
			return
		}
		globalConfig = &cfg
	})
	return nil
}
