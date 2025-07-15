package repository

import (
	"time"

	"AI-PROXY/model"

	"gorm.io/gorm"
)

var db *gorm.DB

// 初始化数据库连接，根据model目录下的数据结构自动创建api配置表
func InitDB(database *gorm.DB) {
	db = database

	// 只进行自动迁移，不删除现有表，保留历史数据
	database.AutoMigrate(&model.APIConfig{})
}

// 查询所有api配置
func GetALLAPIConfig() ([]model.APIConfig, error) {
	var configs []model.APIConfig
	result := db.Find(&configs)
	return configs, result.Error
}

// 查询单个api配置
func GetAPIConfigByName(name string) (*model.APIConfig, error) {
	var config model.APIConfig
	result := db.Where("name=?", name).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// 创建API配置
func CreateAPIConfig(config *model.APIConfig) error {
	return db.Create(config).Error
}

// 更新配置
func UpdateAPIConfig(name string, config *model.APIConfig) error {
	return db.Model(&model.APIConfig{}).Where("name= ?", name).Updates(config).Error
}

// 删除API配置
func DeleteAPIConfig(name string) error {
	return db.Where("name=?", name).Delete(&model.APIConfig{}).Error
}

// GetActiveAPICount 获取活跃API配置数量
func GetActiveAPICount() (int, error) {
	var count int64
	result := db.Model(&model.APIConfig{}).Where("active = ?", true).Count(&count)
	return int(count), result.Error
}

// 更新API测试状态
func UpdateAPITestStatus(name string, status string, testTime time.Time) error {
	return db.Model(&model.APIConfig{}).Where("name = ?", name).Updates(map[string]interface{}{
		"last_test_status": status,
		"last_test_time":   testTime,
	}).Error
}
