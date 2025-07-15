package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger 初始化日志
func InitLogger() error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel("info")
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志输出到控制台
	Logger.SetOutput(os.Stdout)

	return nil
}

// 记录日志信息
func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

// 记录错误日志
func LogError(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Error(message)
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
