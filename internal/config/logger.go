// 配置zap日志工具配置+lumberjack日志切割
package config

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 配置使用zap+lumberjack包进行文件日志记录+日志切割
func InitLogger() *zap.SugaredLogger {
	logWriter := &lumberjack.Logger{
		Filename:   "./运行日志.log",
		MaxSize:    10,   // 日志文件空间限制(MB)
		MaxBackups: 3,    // 保留的最大备份文件数量
		MaxAge:     7,    // 保留日期, 过期后自动删除
		Compress:   true, // 旧的日志文件在滚动时是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05 MST"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(time.FixedZone("UTC+8", 8*3600)).Format("2006-01-02 15:04:05 MST"))
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(logWriter),
		zap.InfoLevel,
	)
	logger := zap.New(core)
	return logger.Sugar()
}
