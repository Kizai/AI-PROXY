package service

import (
	"errors"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/repository"
)

// 获取所有API配置
func GetAllAPIConfigs() ([]model.APIConfig, error) {
	return repository.GetALLAPIConfig()
}

// 获取单个API配置
func GetAPIConfigByName(name string) (*model.APIConfig, error) {
	if name == "" {
		return nil, errors.New("API名称不能为空")
	}
	return repository.GetAPIConfigByName(name)
}

// 创建API配置
func CreateAPIConfig(config *model.APIConfig) error {
	if config.Name == "" || config.BaseURL == "" {
		return errors.New("API名称和基础url不能为空")
	}
	return repository.CreateAPIConfig(config)
}

// 更新API配置
func UpdateAPIConfig(name string, config *model.APIConfig) error {
	if name == "" {
		return errors.New("API名称不能为空")
	}
	return repository.UpdateAPIConfig(name, config)
}

// 删除api配置
func DeleteAPIConfig(name string) error {
	if name == "" {
		return errors.New("API名称不能为空")
	}
	return repository.DeleteAPIConfig(name)
}

// 更新API测试状态
func UpdateAPITestStatus(name string, status string, testTime int64) error {
	return repository.UpdateAPITestStatus(name, status, time.UnixMilli(testTime))
}
