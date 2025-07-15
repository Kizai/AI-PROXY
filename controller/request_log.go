package controller

import (
	"net/http"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/service"
	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

// parseTimeFlexible 灵活解析时间字符串，支持多种格式
func parseTimeFlexible(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// 支持带毫秒和不带毫秒的ISO格式
	layouts := []string{
		time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05.000Z", // "2006-01-02T15:04:05.000Z"
		"2006-01-02T15:04:05Z",     // "2006-01-02T15:04:05Z"
		"2006-01-02 15:04:05",      // 常见格式
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}

func GetRequestLogs(c *gin.Context) {
	var query model.RequestLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequestResponse(c, "参数格式错误"+err.Error())
		return
	}

	// 解析时间参数
	startTime, _ := parseTimeFlexible(query.StartTime)
	endTime, _ := parseTimeFlexible(query.EndTime)

	logs, total, err := service.GetRequestLogsWithTime(&query, startTime, endTime)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, gin.H{
		"logs":  logs,
		"total": total,
	})
}

func DeleteRequestLogs(c *gin.Context) {
	var query model.RequestLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequestResponse(c, "参数格式错误"+err.Error())
		return
	}

	// 解析时间参数
	startTime, _ := parseTimeFlexible(query.StartTime)
	endTime, _ := parseTimeFlexible(query.EndTime)

	count, err := service.DeleteRequestLogsWithTime(&query, startTime, endTime)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, gin.H{
		"message": "日志删除成功",
		"count":   count,
	})
}

// ExportRequestLogs 导出请求日志
func ExportRequestLogs(c *gin.Context) {
	var query model.RequestLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequestResponse(c, "参数格式错误"+err.Error())
		return
	}

	// 解析时间参数
	startTime, _ := parseTimeFlexible(query.StartTime)
	endTime, _ := parseTimeFlexible(query.EndTime)

	csvData, err := service.ExportRequestLogsWithTime(&query, startTime, endTime)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=request-logs.csv")
	c.Data(http.StatusOK, "text/csv", csvData)
}
