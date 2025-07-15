package service

import (
	"fmt"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/repository"
)

func UpdateDailyStatistics(apiName string, date time.Time, statusCode int, responseTime int64) error {
	return repository.UpdateDailyStatistics(apiName, date, statusCode, responseTime)
}

func GetStatistics(query *model.StatisticsQuery) (*model.StatisticsResponse, error) {
	return repository.GetStatistics(query)
}

func GetRealTimeStats() ([]model.RealTimeStats, error) {
	return repository.GetRealTimeStats()
}

// GetActiveAPICount 获取活跃API配置数量
func GetActiveAPICount() (int, error) {
	return repository.GetActiveAPICount()
}

// DebugStatistics 调试统计数据
func DebugStatistics() {
	fmt.Println("=== service.DebugStatistics 被调用 ===")
	repository.DebugStatistics()
	fmt.Println("=== service.DebugStatistics 调用完成 ===")
}
