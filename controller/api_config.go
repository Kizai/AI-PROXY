package controller

import (
	"encoding/json"
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
		fmt.Printf("JSON绑定失败: %v\n", err)
		fmt.Printf("请求体: %s\n", c.Request.Body)
		util.ErrorResponse(c, http.StatusBadRequest, "参数格式错误: "+err.Error())
		return
	}

	fmt.Printf("接收到的API配置: %+v\n", config)

	if err := service.CreateAPIConfig(&config); err != nil {
		fmt.Printf("创建API配置失败: %v\n", err)
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
	var req APITestRequest
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

	// 使用配置的BaseURL进行简单测试
	url := apiConfig.BaseURL
	// 尝试常见的API端点路径
	commonPaths := []string{
		"/v1/chat/completions", // OpenAI格式
		"/v1/completions",      // OpenAI旧格式
		"/chat/completions",    // 简化格式
		"/completions",         // 最简化格式
		"/",                    // 根路径
	}

	// 构建请求头
	headers := map[string]string{}
	if apiConfig.Headers != "" {
		_ = json.Unmarshal([]byte(apiConfig.Headers), &headers)
	}

	// 添加认证信息
	if apiConfig.AuthType == "bearer" && apiConfig.AuthValue != "" {
		headers["Authorization"] = "Bearer " + apiConfig.AuthValue
	}
	if apiConfig.AuthType == "api_key" && apiConfig.AuthValue != "" {
		headers["X-API-Key"] = apiConfig.AuthValue
	}
	if apiConfig.AuthType == "basic" && apiConfig.AuthValue != "" {
		headers["Authorization"] = "Basic " + apiConfig.AuthValue
	}

	// 尝试多个常见路径进行测试
	testBody := `{"test": "api_config_test", "timestamp": "` + time.Now().Format("2006-01-02 15:04:05") + `"}`

	var resp *http.Response
	var duration int64
	var testedURL string
	var foundValidEndpoint bool

	timeout := apiConfig.Timeout
	if timeout <= 0 {
		timeout = 10 // 默认10秒超时
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	// 尝试每个路径
	for _, path := range commonPaths {
		testURL := strings.TrimRight(url, "/") + path
		fmt.Printf("尝试测试路径: %s\n", testURL)

		httpReq, reqErr := http.NewRequest("POST", testURL, strings.NewReader(testBody))
		if reqErr != nil {
			continue
		}

		// 设置请求头
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
		httpReq.Header.Set("Content-Type", "application/json")

		start := time.Now()
		currentResp, currentErr := client.Do(httpReq)
		currentDuration := time.Since(start).Milliseconds()

		fmt.Printf("测试路径 %s - 状态码: %d, 错误: %v\n", testURL,
			func() int {
				if currentResp != nil {
					return currentResp.StatusCode
				}
				return 0
			}(), currentErr)

		if currentErr == nil && currentResp.StatusCode >= 200 && currentResp.StatusCode < 400 {
			// 找到有效的端点
			resp = currentResp
			duration = currentDuration
			testedURL = testURL
			foundValidEndpoint = true
			fmt.Printf("找到可用的API端点: %s\n", testURL)
			break
		} else {
			// 记录最后一个错误，但继续尝试
			if currentResp != nil {
				resp = currentResp
				duration = currentDuration
				testedURL = testURL
			}
			if currentResp != nil {
				currentResp.Body.Close()
			}
		}
	}

	// 记录请求日志
	headersJSON, _ := json.Marshal(headers)
	logErr := service.SaveRequestLog(&model.RequestLog{
		APIName:        req.Name,
		RequestMethod:  "POST",
		RequestPath:    testedURL,
		RequestHeaders: string(headersJSON),
		RequestBody:    testBody,
		ResponseStatus: func() int {
			if resp != nil {
				return resp.StatusCode
			}
			return 0 // 网络错误
		}(),
		ResponseTime: int(duration),
		ErrorMessage: func() string {
			if !foundValidEndpoint && resp != nil {
				return fmt.Sprintf("状态码: %d", resp.StatusCode)
			}
			return ""
		}(),
		UserIP:    c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		CreatedAt: time.Now(),
	})

	if logErr != nil {
		fmt.Printf("记录请求日志失败: %v\n", logErr)
	}

	// 如果没有找到有效端点，返回失败结果
	if !foundValidEndpoint {
		// 更新统计数据 - 失败情况
		updateErr := service.UpdateDailyStatistics(req.Name, time.Now(), 0, duration)
		if updateErr != nil {
			fmt.Printf("更新统计数据失败: %v\n", updateErr)
		}

		// 更新API测试状态为fail
		_ = service.UpdateAPITestStatus(req.Name, "fail", time.Now().UnixMilli())

		// 判断具体的错误类型
		var errorMsg string
		if resp != nil {
			switch resp.StatusCode {
			case 404:
				errorMsg = "API地址可访问，但所有测试路径都返回404（请检查API端点路径是否正确）"
			case 401:
				errorMsg = "API地址可访问，但认证失败（请检查API Key或Token是否正确）"
			case 403:
				errorMsg = "API地址可访问，但权限不足（请检查API Key权限）"
			case 405:
				errorMsg = "API地址可访问，但不支持POST方法（请检查API端点是否正确）"
			default:
				errorMsg = fmt.Sprintf("API响应异常，状态码: %d", resp.StatusCode)
			}
		} else {
			errorMsg = "API地址无法访问，请检查网络连接和URL配置"
		}

		util.SuccessResponse(c, gin.H{
			"success": false,
			"status": func() int {
				if resp != nil {
					return resp.StatusCode
				}
				return 0
			}(),
			"response_time": duration,
			"error":         errorMsg,
			"message":       "API配置测试失败",
		})
		return
	}
	defer resp.Body.Close()

	// 更新统计数据 - 成功情况
	updateErr := service.UpdateDailyStatistics(req.Name, time.Now(), resp.StatusCode, duration)
	if updateErr != nil {
		fmt.Printf("更新统计数据失败: %v\n", updateErr)
	}

	// 更新API测试状态为success
	_ = service.UpdateAPITestStatus(req.Name, "success", time.Now().UnixMilli())

	// 成功情况
	message := fmt.Sprintf("API配置测试成功，找到可用端点: %s", testedURL)

	util.SuccessResponse(c, gin.H{
		"success":       true,
		"status":        resp.StatusCode,
		"response_time": duration,
		"error":         "",
		"message":       message,
	})
}
