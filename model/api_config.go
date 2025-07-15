package model

import (
	"time"

	"gorm.io/gorm"
)

// api配置结构体
type APIConfig struct {
	ID             uint           `json:"id" gorm:"primaryKey"`                     //主键
	Name           string         `json:"name" gorm:"uniqueIndex;size:50"`          //API名称，唯一不可重复
	BaseURL        string         `json:"base_url" gorm:"column:base_url;size:255"` //APIurl
	Description    string         `json:"description" gorm:"size:255"`
	Active         bool           `json:"active" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	LastTestStatus string         `json:"last_test_status" gorm:"column:last_test_status;size:10;default:'never'"` // 最近一次测试状态 success/fail/never
	LastTestTime   *time.Time     `json:"last_test_time" gorm:"column:last_test_time"`                             // 最近一次测试时间
}

func (APIConfig) TableName() string {
	return "api_configs"
}
