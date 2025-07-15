package controller

import (
	"fmt"
	"net/http"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/service"
	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

func GetStatistics(c *gin.Context) {
	var query model.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequestResponse(c, "参数格式错误"+err.Error())
		return
	}
	stats, err := service.GetStatistics(&query)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, stats)
}

// GetStatisticsSummary 获取统计摘要（仪表板用）
func GetStatisticsSummary(c *gin.Context) {
	// 获取所有统计数据，不限制时间范围
	query := &model.StatisticsQuery{
		StartDate: time.Time{}, // 空时间，表示不限制开始时间
		EndDate:   time.Time{}, // 空时间，表示不限制结束时间
	}

	stats, err := service.GetStatistics(query)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取活跃API数量
	activeAPIs, err := service.GetActiveAPICount()
	if err != nil {
		activeAPIs = 0
	}

	summary := gin.H{
		"total_requests":    stats.Summary.TotalRequests,
		"success_requests":  stats.Summary.SuccessCount,
		"error_requests":    stats.Summary.ErrorCount,
		"success_rate":      stats.Summary.SuccessRate,
		"avg_response_time": stats.Summary.AvgRespTime,
		"active_apis":       activeAPIs,
	}

	util.SuccessResponse(c, summary)
}

func GetRealTimeStats(c *gin.Context) {
	stats, err := service.GetRealTimeStats()
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, stats)
}

// GetAPIStatsTable 获取API统计表格数据
func GetAPIStatsTable(c *gin.Context) {
	// 获取所有API配置
	apis, err := service.GetAllAPIConfigs()
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取今日统计数据 - 使用当前日期，不限制时间范围
	query := &model.StatisticsQuery{
		StartDate: time.Time{}, // 空时间，表示不限制开始时间
		EndDate:   time.Time{}, // 空时间，表示不限制结束时间
	}

	stats, err := service.GetStatistics(query)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构建API统计表格数据
	var apiStats []gin.H
	for _, api := range apis {
		// 查找该API的统计数据
		var apiStat gin.H
		for _, stat := range stats.DailyStats {
			if stat.APIName == api.Name {
				successRate := float64(0)
				if stat.TotalRequests > 0 {
					successRate = float64(stat.SuccessRequests) / float64(stat.TotalRequests) * 100
				}
				apiStat = gin.H{
					"name":              api.Name,
					"total_requests":    stat.TotalRequests,
					"success_requests":  stat.SuccessRequests,
					"error_requests":    stat.ErrorCount,
					"success_rate":      successRate,
					"avg_response_time": stat.AvgLatencyMs,
					"active":            api.Active,
				}
				break
			}
		}

		// 如果没有找到统计数据，使用默认值
		if apiStat == nil {
			apiStat = gin.H{
				"name":              api.Name,
				"total_requests":    0,
				"success_requests":  0,
				"error_requests":    0,
				"success_rate":      0,
				"avg_response_time": 0,
				"active":            api.Active,
			}
		}

		apiStats = append(apiStats, apiStat)
	}

	util.SuccessResponse(c, apiStats)
}

// DebugStatistics 调试统计数据接口
func DebugStatistics(c *gin.Context) {
	fmt.Println("=== DebugStatistics 接口被调用 ===")
	service.DebugStatistics()
	fmt.Println("=== DebugStatistics 接口调用完成 ===")
	util.SuccessResponse(c, gin.H{"message": "调试信息已输出到控制台"})
}

// TestStatistics 测试统计更新接口
func TestStatistics(c *gin.Context) {
	apiName := c.Query("api_name")
	if apiName == "" {
		apiName = "test-api"
	}

	// 手动触发统计更新
	err := service.UpdateDailyStatistics(apiName, time.Now(), 200, 150)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "统计更新失败: "+err.Error())
		return
	}

	util.SuccessResponse(c, gin.H{
		"message":       "测试统计更新成功",
		"api_name":      apiName,
		"status_code":   200,
		"response_time": 150,
	})
}
