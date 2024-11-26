package main

import (
	"fmt"
	"net/http"
	"os"

	"KazeFrame/internal/config"
	"KazeFrame/internal/dao"
	"KazeFrame/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 载入配置、测试数据库、Redis连接
	if err := config.InitServer(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// 载入zap日志工具、配置信息、数据库实例
	zapLog := config.GetLogger()
	cfg := config.GetConfig()
	db := config.GetDB()
	// 初始化/internal/dao包中的仓库实例
	dao.InitRepo(config.GetDB())
	// 插入基础数据到数据库表
	config.Seed(db)
	// 通过配置文件debug变量选择Gin运行模式
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 启动HTTP服务
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.RunServer(),
	}
	// 记录启动日志、打印到控制台
	zapLog.Info(">>> 服务启动于端口:", cfg.Server.Port)
	fmt.Printf("\n\n>>> 服务启动成功, 可查看: http://localhost:%s/ui/welcome.html\n\n", cfg.Server.Port)
	// 开始监听请求
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zapLog.Fatalf("XXX 服务端启动失败: %v\n", err)
	}
}
