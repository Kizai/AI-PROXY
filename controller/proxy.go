package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"AI-PROXY/model"
	"AI-PROXY/service"
	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

// ForwardRequest 代理转发请求
func ForwardRequest(c *gin.Context) {
	fmt.Printf("=== 代理转发被调用 ===\n")
	fmt.Printf("完整请求路径: %s\n", c.Request.URL.Path)
	fmt.Printf("请求方法: %s\n", c.Request.Method)

	startTime := time.Now()

	// 获取 API 名称和路径
	apiName := c.Param("apiName")
	path := c.Param("path")

	fmt.Printf("解析的API名称: %s\n", apiName)
	fmt.Printf("解析的路径: %s\n", path)

	// 记录请求日志
	requestLog := &model.RequestLog{
		APIName:       apiName,
		RequestPath:   path,
		RequestMethod: c.Request.Method,
		UserIP:        c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		CreatedAt:     time.Now(),
	}

	// 获取 API 配置
	apiConfig, err := service.GetAPIConfigByName(apiName)
	if err != nil {
		requestLog.ResponseStatus = http.StatusNotFound
		requestLog.ErrorMessage = "API配置不存在: " + err.Error()
		service.SaveRequestLog(requestLog)
		util.ErrorResponse(c, http.StatusNotFound, "API配置不存在: "+apiName)
		return
	}

	// 控制台调试输出
	fmt.Printf("代理请求 - API名称: %s\n", apiName)
	fmt.Printf("代理请求 - 原始路径: %s\n", path)
	fmt.Printf("代理请求 - 方法: %s\n", c.Request.Method)
	fmt.Printf("代理请求 - 完整URL: %s\n", apiConfig.BaseURL+path)

	// 构建目标 URL
	targetURL := apiConfig.BaseURL
	if !strings.HasSuffix(targetURL, "/") {
		targetURL += "/"
	}
	targetURL += strings.TrimPrefix(path, "/")

	// 添加调试日志
	fmt.Printf("代理请求 - API名称: %s\n", apiName)
	fmt.Printf("代理请求 - 原始路径: %s\n", path)
	fmt.Printf("代理请求 - 目标URL: %s\n", targetURL)
	fmt.Printf("代理请求 - 方法: %s\n", c.Request.Method)

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		requestLog.ResponseStatus = http.StatusBadRequest
		requestLog.ErrorMessage = "读取请求体失败: " + err.Error()
		service.SaveRequestLog(requestLog)
		util.ErrorResponse(c, http.StatusBadRequest, "读取请求体失败")
		return
	}
	requestLog.RequestBody = string(body)

	// 构建请求头
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// 添加 API 配置中的认证信息
	if apiConfig.AuthType == "bearer" {
		headers["Authorization"] = "Bearer " + apiConfig.AuthValue
	} else if apiConfig.AuthType == "api_key" {
		headers["X-API-Key"] = apiConfig.AuthValue
	}

	// 添加 API 配置中的自定义请求头
	if apiConfig.Headers != "" {
		// 这里可以解析 Headers 字段（JSON 格式）
		// 暂时跳过，因为当前 Headers 是 string 类型
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(apiConfig.Timeout) * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		requestLog.ResponseStatus = http.StatusInternalServerError
		requestLog.ErrorMessage = "创建请求失败: " + err.Error()
		service.SaveRequestLog(requestLog)
		util.ErrorResponse(c, http.StatusInternalServerError, "创建请求失败")
		return
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 打印请求头调试信息
	fmt.Printf("代理请求 - 请求头: %+v\n", req.Header)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		requestLog.ResponseStatus = http.StatusBadGateway
		requestLog.ErrorMessage = "请求失败: " + err.Error()
		service.SaveRequestLog(requestLog)
		util.ErrorResponse(c, http.StatusBadGateway, "请求失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		requestLog.ResponseStatus = http.StatusInternalServerError
		requestLog.ErrorMessage = "读取响应体失败: " + err.Error()
		service.SaveRequestLog(requestLog)
		util.ErrorResponse(c, http.StatusInternalServerError, "读取响应体失败")
		return
	}

	// 设置响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 记录请求日志
	requestLog.ResponseStatus = resp.StatusCode
	requestLog.ResponseTime = int(time.Since(startTime).Milliseconds())
	service.SaveRequestLog(requestLog)

	// 更新统计数据
	fmt.Printf("=== 开始更新统计数据 ===\n")
	fmt.Printf("API名称: %s\n", apiName)
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应时间: %dms\n", requestLog.ResponseTime)
	fmt.Printf("当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	err = service.UpdateDailyStatistics(apiName, time.Now(), resp.StatusCode, int64(requestLog.ResponseTime))
	if err != nil {
		fmt.Printf("统计更新失败: %v\n", err)
	} else {
		fmt.Printf("统计更新成功\n")
	}
	fmt.Printf("=== 统计更新结束 ===\n")

	// 返回响应
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
