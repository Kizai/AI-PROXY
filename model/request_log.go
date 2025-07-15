package model

import (
	"time"

	"gorm.io/gorm"
)

// 请求日志结构体
type RequestLog struct {
	ID             uint           `json:"id"grom:"primarykey"`
	APIName        string         `json:"api_name"grom:"size:50"`
	RequestPath    string         `json:"request_path"grom:"size:255"`
	RequestMethod  string         `json:"request_method"grom:"size:10"` //get.post.put,delete
	RequestHeaders string         `json:"request_headers"grom:"type:text"`
	RequestBody    string         `json:"request_body"grom:"type:text"`
	ResponseStatus int            `json:"response_status"`
	ResponseTime   int            `json:"response_time"`
	ErrorMessage   string         `json:"error_message"grom:"type:text"`
	UserIP         string         `json:"user_ip"grom:"size:50"`
	UserAgent      string         `json:"user_agent"grom:"size:255"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"grom:"index"`
}

func (RequestLog) TableName() string {
	return "request_logs"
}

// 请求日志查询结构体
type RequestLogQuery struct {
	APIName       string `form:"api_name" json:"api_name"`
	RequestPath   string `form:"request_path" json:"request_path"`
	RequestMethod string `form:"request_method" json:"request_method"`
	StatusCode    int    `form:"status_code" json:"status_code"`
	StartTime     string `form:"start_time" json:"start_time"`
	EndTime       string `form:"end_time" json:"end_time"`
	HasError      bool   `form:"has_error" json:"has_error"`
	Page          int    `form:"page" json:"page"`
}

// 请求日志响应结构体
type RequestLogResponse struct {
	Total int64        `json:"total"`
	Logs  []RequestLog `json:"logs"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Pages int          `json:"pages"`
}
