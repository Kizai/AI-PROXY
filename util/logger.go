package util

import (
	"io"
	"os"

	"AI-PROXY/config"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志，支持同时输出到控制台和文件
func InitLogger(logConfig *config.LogConfig) error {
	Logger = logrus.New()
	level, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 打开日志文件
	file, err := os.OpenFile(logConfig.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.SetOutput(os.Stdout)
		return err
	}
	// 同时输出到控制台和文件
	mw := io.MultiWriter(os.Stdout, file)
	Logger.SetOutput(mw)
	return nil
}

// 记录请求日志
func LogRequest(apiName, method, path string, statusCode int, responseTime int64, errorMessage string) {
	fields := logrus.Fields{
		"api_name":      apiName,
		"method":        method,
		"path":          path,
		"status_code":   statusCode,
		"response_time": responseTime,
		"error_message": errorMessage,
	}
	if statusCode >= 200 && statusCode < 400 {
		Logger.WithFields(fields).Info("请求成功")
	} else {
		Logger.WithFields(fields).Error("请求失败")
	}
}
