package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/service"
	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

// 获取所有API配置
func GetAllAPIConfigs(c *gin.Context) {
	configs, err := service.GetAllAPIConfigs()
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	util.SuccessResponse(c, configs)
}

// 获取单个API配置
func GetAPIConfig(c *gin.Context) {
	name := c.Param("name")
	config, err := service.GetAPIConfigByName(name)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	util.SuccessResponse(c, config)
}

// 创建API配置
func CreateAPIConfig(c *gin.Context) {
	var config model.APIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "参数格式错误: "+err.Error())
		return
	}
	if config.Name == "" || config.BaseURL == "" {
		util.ErrorResponse(c, http.StatusBadRequest, "API名称和基址URL为必填项")
		return
	}
	if err := service.CreateAPIConfig(&config); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	util.SuccessResponse(c, "API配置创建成功")
}

// 更新API配置
func UpdateAPIConfig(c *gin.Context) {
	name := c.Param("name")
	var config model.APIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		util.BadRequestResponse(c, "参数格式错误："+err.Error())
		return
	}
	if err := service.UpdateAPIConfig(name, &config); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	util.SuccessResponse(c, "API配置更新成功")
}

// 删API配置
func DeleteAPIConfig(c *gin.Context) {
	name := c.Param("name")
	if err := service.DeleteAPIConfig(name); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	util.SuccessResponse(c, "API配置删除成功")
}

// API测试请求体结构
type APITestRequest struct {
	Name string `json:"name"`
}

// API测试响应结构
type APITestResponse struct {
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
	ResponseTime int64  `json:"response_time"`
	Error        string `json:"error"`
	Message      string `json:"message"`
}

func TestAPIConfig(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 查找API配置
	apiConfig, err := service.GetAPIConfigByName(req.Name)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "API配置不存在: "+err.Error())
		return
	}

	// 自动补全协议
	url := apiConfig.BaseURL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	client := &http.Client{Timeout: 10 * time.Second}
	reqTest, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "创建请求失败")
		return
	}
	start := time.Now()
	resp, err := client.Do(reqTest)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		util.SuccessResponse(c, gin.H{
			"success":       false,
			"status":        0,
			"response_time": duration,
			"error":         err.Error(),
			"message":       "API地址无法访问: " + err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 400
	msg := "API地址可访问"
	if !success {
		msg = "API地址无法访问，状态码: " + fmt.Sprint(resp.StatusCode)
	}
	util.SuccessResponse(c, gin.H{
		"success":       success,
		"status":        resp.StatusCode,
		"response_time": duration,
		"error":         "",
		"message":       msg,
	})
}
