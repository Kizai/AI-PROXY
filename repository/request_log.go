package repository

import (
	"fmt"
	"time"

	"AI-PROXY/model"
)

// 保存请求日志
func SaveRequestLog(log *model.RequestLog) error {
	return db.Create(log).Error
}

// 查询请求日志
func GetRequestLogs(query *model.RequestLogQuery) ([]model.RequestLog, int64, error) {
	var logs []model.RequestLog
	var total int64

	dbQuery := db.Model(&model.RequestLog{})
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name=?", query.APIName)
	}
	if query.RequestPath != "" {
		dbQuery = dbQuery.Where("request_path LIKE ?", "%"+query.RequestPath+"%")
	}
	if query.RequestMethod != "" {
		dbQuery = dbQuery.Where("request_method=?", query.RequestMethod)
	}
	if query.StatusCode != 0 {
		dbQuery = dbQuery.Where("status_code=?", query.StatusCode)
	}
	if query.HasError {
		dbQuery = dbQuery.Where("error_message != ''")
	} else {
		dbQuery = dbQuery.Where("error_message = ''")
	}

	// 注意：旧方法不再使用时间筛选，新方法使用GetRequestLogsWithTime
	// 这里保留兼容性，但不进行时间筛选

	//统计总数
	dbQuery.Count(&total)

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := 10 // 固定页面大小
	offset := (page - 1) * pageSize

	result := dbQuery.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs)
	return logs, total, result.Error
}

// GetRequestLogsWithTime 查询请求日志（带时间参数）
func GetRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) ([]model.RequestLog, int64, error) {
	var logs []model.RequestLog
	var total int64

	dbQuery := db.Model(&model.RequestLog{})
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name=?", query.APIName)
	}
	if query.RequestPath != "" {
		dbQuery = dbQuery.Where("request_path LIKE ?", "%"+query.RequestPath+"%")
	}
	if query.RequestMethod != "" {
		dbQuery = dbQuery.Where("request_method=?", query.RequestMethod)
	}
	if query.StatusCode != 0 {
		dbQuery = dbQuery.Where("status_code=?", query.StatusCode)
	}
	if query.HasError {
		dbQuery = dbQuery.Where("error_message != ''")
	} else {
		dbQuery = dbQuery.Where("error_message = ''")
	}

	// 使用解析后的时间参数
	if !startTime.IsZero() {
		dbQuery = dbQuery.Where("created_at>=?", startTime)
	}
	if !endTime.IsZero() {
		dbQuery = dbQuery.Where("created_at<=?", endTime)
	}

	//统计总数
	dbQuery.Count(&total)

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := 10 // 固定页面大小
	offset := (page - 1) * pageSize

	result := dbQuery.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs)
	return logs, total, result.Error
}

func DeleteRequestLogs(query *model.RequestLogQuery) (int64, error) {
	fmt.Printf("=== 开始删除请求日志 ===\n")

	// 先统计要删除的记录数
	var count int64
	dbQuery := db.Model(&model.RequestLog{})

	// 应用查询条件
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name=?", query.APIName)
	}
	if query.RequestPath != "" {
		dbQuery = dbQuery.Where("request_path LIKE ?", "%"+query.RequestPath+"%")
	}
	if query.RequestMethod != "" {
		dbQuery = dbQuery.Where("request_method=?", query.RequestMethod)
	}
	if query.StatusCode != 0 {
		dbQuery = dbQuery.Where("response_status=?", query.StatusCode)
	}
	if query.HasError {
		dbQuery = dbQuery.Where("error_message != ''")
	} else {
		dbQuery = dbQuery.Where("error_message = ''")
	}

	// 统计要删除的记录数
	dbQuery.Count(&count)
	fmt.Printf("将要删除 %d 条记录\n", count)

	// 如果记录数太多，分批删除
	if count > 10000 {
		fmt.Printf("记录数较多，分批删除...\n")
		var deletedCount int64
		batchSize := 1000

		for {
			var batchCount int64
			result := dbQuery.Limit(batchSize).Delete(&model.RequestLog{})
			if result.Error != nil {
				fmt.Printf("分批删除失败: %v\n", result.Error)
				return deletedCount, result.Error
			}
			batchCount = result.RowsAffected
			deletedCount += batchCount
			fmt.Printf("已删除 %d 条记录，总计 %d 条\n", batchCount, deletedCount)

			if batchCount < int64(batchSize) {
				break
			}
		}

		fmt.Printf("=== 删除完成，共删除 %d 条记录 ===\n", deletedCount)
		return deletedCount, nil
	} else {
		// 记录数不多，直接删除
		result := dbQuery.Delete(&model.RequestLog{})
		fmt.Printf("=== 删除完成，共删除 %d 条记录 ===\n", result.RowsAffected)
		return result.RowsAffected, result.Error
	}
}

// DeleteRequestLogsWithTime 删除请求日志（带时间参数）
func DeleteRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) (int64, error) {
	dbQuery := db.Model(&model.RequestLog{})

	// 使用解析后的时间参数
	if !startTime.IsZero() {
		dbQuery = dbQuery.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		dbQuery = dbQuery.Where("created_at <= ?", endTime)
	}
	result := dbQuery.Delete(&model.RequestLog{})
	return result.RowsAffected, result.Error
}

// ExportRequestLogs 导出请求日志为CSV
func ExportRequestLogs(query *model.RequestLogQuery) ([]byte, error) {
	var logs []model.RequestLog

	dbQuery := db.Model(&model.RequestLog{})
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name=?", query.APIName)
	}
	if query.RequestPath != "" {
		dbQuery = dbQuery.Where("request_path LIKE ?", "%"+query.RequestPath+"%")
	}
	if query.RequestMethod != "" {
		dbQuery = dbQuery.Where("request_method=?", query.RequestMethod)
	}
	if query.StatusCode != 0 {
		dbQuery = dbQuery.Where("status_code=?", query.StatusCode)
	}
	if query.HasError {
		dbQuery = dbQuery.Where("error_message != ''")
	} else {
		dbQuery = dbQuery.Where("error_message = ''")
	}
	// 注意：旧方法不再使用时间筛选，新方法使用ExportRequestLogsWithTime
	// 这里保留兼容性，但不进行时间筛选

	result := dbQuery.Order("created_at DESC").Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	// 生成CSV数据
	csvData := "时间,API名称,请求路径,请求方法,状态码,响应时间(ms),错误信息,用户IP\n"
	for _, log := range logs {
		csvData += fmt.Sprintf("%s,%s,%s,%s,%d,%d,%s,%s\n",
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.APIName,
			log.RequestPath,
			log.RequestMethod,
			log.ResponseStatus,
			log.ResponseTime,
			log.ErrorMessage,
			log.UserIP,
		)
	}

	return []byte(csvData), nil
}

// ExportRequestLogsWithTime 导出请求日志为CSV（带时间参数）
func ExportRequestLogsWithTime(query *model.RequestLogQuery, startTime, endTime time.Time) ([]byte, error) {
	var logs []model.RequestLog

	dbQuery := db.Model(&model.RequestLog{})
	if query.APIName != "" {
		dbQuery = dbQuery.Where("api_name=?", query.APIName)
	}
	if query.RequestPath != "" {
		dbQuery = dbQuery.Where("request_path LIKE ?", "%"+query.RequestPath+"%")
	}
	if query.RequestMethod != "" {
		dbQuery = dbQuery.Where("request_method=?", query.RequestMethod)
	}
	if query.StatusCode != 0 {
		dbQuery = dbQuery.Where("status_code=?", query.StatusCode)
	}
	if query.HasError {
		dbQuery = dbQuery.Where("error_message != ''")
	} else {
		dbQuery = dbQuery.Where("error_message = ''")
	}

	// 使用解析后的时间参数
	if !startTime.IsZero() {
		dbQuery = dbQuery.Where("created_at>=?", startTime)
	}
	if !endTime.IsZero() {
		dbQuery = dbQuery.Where("created_at<=?", endTime)
	}

	result := dbQuery.Order("created_at DESC").Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	// 生成CSV数据
	csvData := "时间,API名称,请求路径,请求方法,状态码,响应时间(ms),错误信息,用户IP\n"
	for _, log := range logs {
		csvData += fmt.Sprintf("%s,%s,%s,%s,%d,%d,%s,%s\n",
			log.CreatedAt.Format("2006-01-02 15:04:05"),
			log.APIName,
			log.RequestPath,
			log.RequestMethod,
			log.ResponseStatus,
			log.ResponseTime,
			log.ErrorMessage,
			log.UserIP,
		)
	}

	return []byte(csvData), nil
}
