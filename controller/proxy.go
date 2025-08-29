package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"AI-PROXY/service"
	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

// ForwardRequest 代理转发请求
func ForwardRequest(c *gin.Context) {
	fmt.Printf("=== 代理转发被调用 ===\n")
	fmt.Printf("完整请求路径: %s\n", c.Request.URL.Path)
	fmt.Printf("请求方法: %s\n", c.Request.Method)

	// 获取 API 名称和路径
	apiName := c.Param("apiName")
	path := c.Param("path")
	
	// 重要：保留原始查询参数
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	fmt.Printf("解析的API名称: %s\n", apiName)
	fmt.Printf("解析的路径: %s\n", path)
	fmt.Printf("原始查询参数: %s\n", c.Request.URL.RawQuery)

	// 移除requestLog相关的定义、赋值、所有service.SaveRequestLog调用及相关逻辑

	// 获取 API 配置
	apiConfig, err := service.GetAPIConfigByName(apiName)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, "API配置不存在: "+apiName)
		return
	}
	// 新增：未启用的API禁止访问
	if !apiConfig.Active {
		util.ErrorResponse(c, http.StatusForbidden, "该API已被禁用")
		return
	}

	// 控制台调试输出
	fmt.Printf("代理请求 - API名称: %s\n", apiName)
	fmt.Printf("代理请求 - 原始路径: %s\n", path)
	fmt.Printf("代理请求 - 方法: %s\n", c.Request.Method)
	fmt.Printf("代理请求 - 完整URL: %s\n", apiConfig.BaseURL+path)

	// 构建目标 URL（只做base_url和path原样拼接，不做任何补全）
	targetURL := strings.TrimRight(apiConfig.BaseURL, "/") + path
	// 自动补全协议，优先https
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	// 特殊处理Gemini API的认证方式
	if apiName == "gemini" {
		fmt.Printf("🔍 Gemini特殊处理 - 原始targetURL: %s\n", targetURL)
		
		// 检查URL中是否已经包含key参数
		if !strings.Contains(targetURL, "key=") {
			fmt.Printf("🔍 URL中未包含key参数，尝试从Authorization头提取\n")
			// 从Authorization头中提取API Key并添加到URL查询参数
			authHeader := c.GetHeader("Authorization")
			fmt.Printf("🔍 Authorization头: %s\n", authHeader)
			
			if authHeader != "" {
				// 支持 "Bearer API_KEY" 或 "API_KEY" 格式
				apiKey := strings.TrimPrefix(authHeader, "Bearer ")
				apiKey = strings.TrimSpace(apiKey)
				fmt.Printf("🔍 提取的API Key: %s\n", apiKey)
				
				// 添加key参数到URL
				separator := "?"
				if strings.Contains(targetURL, "?") {
					separator = "&"
				}
				targetURL = targetURL + separator + "key=" + apiKey
				fmt.Printf("🔍 添加key后的targetURL: %s\n", targetURL)
			} else {
				fmt.Printf("❌ 未找到Authorization头\n")
			}
		} else {
			fmt.Printf("✅ URL中已包含key参数，直接使用\n")
		}
	}

	// 添加调试日志
	fmt.Printf("代理请求 - API名称: %s\n", apiName)
	fmt.Printf("代理请求 - 原始路径: %s\n", path)
	fmt.Printf("代理请求 - 目标URL: %s\n", targetURL)
	fmt.Printf("代理请求 - 方法: %s\n", c.Request.Method)

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "读取请求体失败")
		return
	}

	// 构建请求头
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			// 对于Gemini API，跳过Authorization头，因为我们已经将其转换为URL参数
			if apiName == "gemini" && strings.ToLower(key) == "authorization" {
				continue
			}
			headers[key] = values[0]
		}
	}

	// 创建 HTTP 客户端（使用默认超时）
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
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
		util.ErrorResponse(c, http.StatusBadGateway, "请求失败: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "读取响应体失败")
		return
	}

	// 设置响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 返回响应
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
