package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 系统配置结构
type Config struct {
	Server   ServerConfig         `json:"server"`
	Database DatabaseConfig       `json:"database"`
	Log      LogConfig            `json:"log"`
	APIs     map[string]APIConfig `json:"apis"`
	Auth     AuthConfig           `json:"auth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// APIConfig API配置
type APIConfig struct {
	BaseURL     string            `json:"base_url"`
	Headers     map[string]string `json:"headers"`
	AuthType    string            `json:"auth_type"`
	AuthValue   string            `json:"auth_value"`
	Timeout     int               `json:"timeout"`
	RateLimit   int               `json:"rate_limit"`
	Description string            `json:"description"`
}

// AuthConfig 管理员认证配置
type AuthConfig struct {
	Token string `json:"token"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// 全局配置变量，供全项目访问
var GlobalConfig *Config
