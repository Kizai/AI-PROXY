package model

import (
	"time"

	"gorm.io/gorm"
)

// 每日统计数据模型
type DailyStatistics struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Date            time.Time      `json:"date" gorm:"uniqueIndex:idx_date_api;not null"`
	APIName         string         `json:"api_name" gorm:"uniqueIndex:idx_date_api;size:50;not null"`
	TotalRequests   int            `json:"total_requests" gorm:"default:0"`   //总请求数
	SuccessRequests int            `json:"success_requests" gorm:"default:0"` //成功请求数
	ErrorCount      int            `json:"error_count" gorm:"default:0"`      //错误请求数
	AvgLatencyMs    float64        `json:"avg_latency_ms" gorm:"default:0.0"` //平均延迟(毫秒)
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DailyStatistics) TableName() string {
	return "daily_statistics"
}

// StatisticsQuery 统计查询参数
type StatisticsQuery struct {
	APIName   string    `json:"api_name" form:"api_name"`
	StartDate time.Time `json:"start_date" form:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date"`
	Page      int       `json:"page" form:"page"`
	PageSize  int       `json:"page_size" form:"page_size"`
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	Summary    StatisticsSummary `json:"summary"`
	DailyStats []DailyStatistics `json:"daily_stats"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	Pages      int               `json:"pages"`
}

// StatisticsSummary 统计摘要
type StatisticsSummary struct {
	TotalRequests int64   `json:"total_requests"`
	SuccessCount  int64   `json:"success_count"`
	ErrorCount    int64   `json:"error_count"`
	SuccessRate   float64 `json:"success_rate"`
	ErrorRate     float64 `json:"error_rate"`
	AvgRespTime   float64 `json:"avg_resp_time"`
}

// RealTimeStats 实时统计数据
type RealTimeStats struct {
	APIName       string  `json:"api_name"`
	TotalRequests int64   `json:"total_requests"`
	SuccessCount  int64   `json:"success_count"`
	ErrorCount    int64   `json:"error_count"`
	SuccessRate   float64 `json:"success_rate"`
	ErrorRate     float64 `json:"error_rate"`
	AvgRespTime   float64 `json:"avg_resp_time"`
}
