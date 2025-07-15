package service

import (
	"time"

	"AI-PROXY/model"
	"AI-PROXY/repository"
)

// SaveRequestLog 保存请求日志
func SaveRequestLog(log *model.RequestLog) error {
	return repository.SaveRequestLog(log)
}

// GetRequestLogs 获取请求日志列表
func GetRequestLogs(query *model.RequestLogQuery) ([]model.RequestLog, int64, error) {
	return repository.GetRequestLogs(query)
}

// GetRequestLogsWithTime 获取请求日志列表（带时间参数）
func GetRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) ([]model.RequestLog, int64, error) {
	return repository.GetRequestLogsWithTime(query, startTime, endTime)
}

// DeleteRequestLogs 删除请求日志
func DeleteRequestLogs(query *model.RequestLogQuery) (int64, error) {
	return repository.DeleteRequestLogs(query)
}

// DeleteRequestLogsWithTime 删除请求日志（带时间参数）
func DeleteRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) (int64, error) {
	return repository.DeleteRequestLogsWithTime(query, startTime, endTime)
}

// ExportRequestLogs 导出请求日志为CSV
func ExportRequestLogs(query *model.RequestLogQuery) ([]byte, error) {
	return repository.ExportRequestLogs(query)
}

// ExportRequestLogsWithTime 导出请求日志为CSV（带时间参数）
func ExportRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) ([]byte, error) {
	return repository.ExportRequestLogsWithTime(query, startTime, endTime)
}
