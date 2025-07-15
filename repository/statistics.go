package repository

import (
	"fmt"
	"time"

	"AI-PROXY/model"

	"gorm.io/gorm"
)

//更新每日统计数据

func UpdateDailyStatistics(apiName string, date time.Time, statusCode int, responseTime int64) error {
	// 检查数据库连接
	if db == nil {
		fmt.Printf("错误：数据库连接为空\n")
		return fmt.Errorf("数据库连接为空")
	}

	var stat model.DailyStatistics
	// 统一使用日期零点，确保每天只有一条记录
	date = date.Truncate(24 * time.Hour)
	dateStr := date.Format("2006-01-02")

	fmt.Printf("开始更新统计数据 - API: %s, 日期: %s, 状态码: %d, 响应时间: %dms\n", apiName, dateStr, statusCode, responseTime)

	// 使用事务确保数据一致性
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找或创建记录
	result := tx.Where("api_name = ? AND DATE(date) = ?", apiName, dateStr).First(&stat)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fmt.Printf("查询统计数据失败: %v\n", result.Error)
		tx.Rollback()
		return result.Error
	}

	// 如果记录不存在，创建新记录
	if result.RowsAffected == 0 {
		fmt.Printf("创建新的统计记录\n")
		stat = model.DailyStatistics{
			APIName:         apiName,
			Date:            date, // 使用零点时间
			TotalRequests:   1,
			SuccessRequests: 0,
			ErrorCount:      0,
			AvgLatencyMs:    float64(responseTime),
		}
		if statusCode >= 200 && statusCode < 400 {
			stat.SuccessRequests = 1
		} else {
			stat.ErrorCount = 1
		}
		err := tx.Create(&stat).Error
		if err != nil {
			fmt.Printf("创建统计记录失败: %v\n", err)
			tx.Rollback()
			return err
		}
		fmt.Printf("创建统计记录成功\n")
	} else {
		// 如果记录存在，更新统计数据
		fmt.Printf("更新现有统计记录 - 原总数: %d\n", stat.TotalRequests)
		stat.TotalRequests++
		if statusCode >= 200 && statusCode < 400 {
			stat.SuccessRequests++
		} else {
			stat.ErrorCount++
		}
		// 平均响应时间更新
		stat.AvgLatencyMs = (stat.AvgLatencyMs*float64(stat.TotalRequests-1) + float64(responseTime)) / float64(stat.TotalRequests)
		err := tx.Save(&stat).Error
		if err != nil {
			fmt.Printf("更新统计记录失败: %v\n", err)
			tx.Rollback()
			return err
		}
		fmt.Printf("更新统计记录成功 - 新总数: %d\n", stat.TotalRequests)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("提交事务失败: %v\n", err)
		return err
	}

	return nil
}

// 根据查询条件统计数据和摘要
func GetStatistics(query *model.StatisticsQuery) (*model.StatisticsResponse, error) {
	var stats []model.DailyStatistics
	var total int64

	//查询条件拼接
	dbQuery := db.Model(&model.DailyStatistics{})
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name = ?", query.APIName)
	}
	if !query.StartDate.IsZero() {
		// 使用日期字符串格式进行查询
		startDateStr := query.StartDate.Format("2006-01-02")
		dbQuery = dbQuery.Where("DATE(date) >= ?", startDateStr)
	}
	if !query.EndDate.IsZero() {
		// 使用日期字符串格式进行查询
		endDateStr := query.EndDate.Format("2006-01-02")
		dbQuery = dbQuery.Where("DATE(date) <= ?", endDateStr)
	}

	//统计总数
	dbQuery.Count(&total)

	// 如果没有时间限制，查询所有数据（不分页）
	if query.StartDate.IsZero() && query.EndDate.IsZero() {
		fmt.Printf("查询所有统计数据（无时间限制）\n")
		result := dbQuery.Order("date ASC").Find(&stats)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		//分页处理，保持默认值，防止出错
		page := query.Page
		if page < 1 {
			page = 1
		}
		pageSize := query.PageSize
		if pageSize < 1 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize

		//按照结果升序，结果放在stat切片里
		result := dbQuery.Order("date ASC").Limit(pageSize).Offset(offset).Find(&stats)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	fmt.Printf("查询到 %d 条统计记录\n", len(stats))

	// 统计摘要
	summary := model.StatisticsSummary{}
	for _, stat := range stats {
		summary.TotalRequests += int64(stat.TotalRequests)
		summary.SuccessCount += int64(stat.SuccessRequests)
		summary.ErrorCount += int64(stat.ErrorCount)
		summary.AvgRespTime += stat.AvgLatencyMs
	}
	if len(stats) > 0 {
		summary.AvgRespTime /= float64(len(stats))
		if summary.TotalRequests > 0 {
			summary.SuccessRate = float64(summary.SuccessCount) / float64(summary.TotalRequests) * 100
			summary.ErrorRate = float64(summary.ErrorCount) / float64(summary.TotalRequests) * 100
		} else {
			summary.SuccessRate = 0
			summary.ErrorRate = 0
		}
	} else {
		summary.AvgRespTime = 0
		summary.SuccessRate = 0
		summary.ErrorRate = 0
	}

	fmt.Printf("统计摘要 - 总请求: %d, 成功: %d, 错误: %d, 成功率: %.2f%%\n",
		summary.TotalRequests, summary.SuccessCount, summary.ErrorCount, summary.SuccessRate)

	pages := 0
	if query.PageSize > 0 {
		pages = int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	}
	return &model.StatisticsResponse{
		Summary:    summary,
		DailyStats: stats,
		Total:      total,
		Page:       query.Page,
		Size:       query.PageSize,
		Pages:      pages,
	}, nil
}

func GetRealTimeStats() ([]model.RealTimeStats, error) {
	var stats []model.RealTimeStats
	today := time.Now().Format("2006-01-02")
	var dailyStats []model.DailyStatistics
	result := db.Where("date = ?", today).Find(&dailyStats)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, stat := range dailyStats {
		successRate := float64(0)
		errorRate := float64(0)
		if stat.TotalRequests > 0 {
			successRate = float64(stat.SuccessRequests) / float64(stat.TotalRequests) * 100
			errorRate = float64(stat.ErrorCount) / float64(stat.TotalRequests) * 100
		}
		stats = append(stats, model.RealTimeStats{
			APIName:       stat.APIName,
			TotalRequests: int64(stat.TotalRequests),
			SuccessCount:  int64(stat.SuccessRequests),
			ErrorCount:    int64(stat.ErrorCount),
			SuccessRate:   successRate,
			ErrorRate:     errorRate,
			AvgRespTime:   stat.AvgLatencyMs,
		})
	}
	return stats, nil
}

// DebugStatistics 调试函数：打印所有统计数据
func DebugStatistics() {
	fmt.Println("=== 调试：所有统计数据 ===")

	// 检查数据库连接
	if db == nil {
		fmt.Println("错误：数据库连接为空")
		return
	}

	// 测试数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("获取数据库连接失败: %v\n", err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("数据库连接测试失败: %v\n", err)
		return
	}
	fmt.Println("数据库连接正常")

	var stats []model.DailyStatistics
	result := db.Find(&stats)
	if result.Error != nil {
		fmt.Printf("查询统计数据失败: %v\n", result.Error)
		return
	}

	fmt.Printf("总共找到 %d 条统计记录\n", len(stats))
	for _, stat := range stats {
		fmt.Printf("API: %s, 日期: %s, 总请求: %d, 成功: %d, 错误: %d, 平均延迟: %.1fs\n",
			stat.APIName,
			stat.Date.Format("2006-01-02 15:04:05"),
			stat.TotalRequests,
			stat.SuccessRequests,
			stat.ErrorCount,
			stat.AvgLatencyMs/1000.0)
	}
	fmt.Println("=== 调试结束 ===")
}
