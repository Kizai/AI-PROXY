package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 创建=http客户端结构体
type HTTPClient struct {
	client *http.Client
}

// 创建http客户端
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// DoRequest 执行HTTP请求
func (c *HTTPClient) DoRequest(method, url string, headers map[string]string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 执行请求
	startTime := time.Now()
	resp, err := c.client.Do(req)
	requestTime := time.Since(startTime).Milliseconds()

	if err != nil {
		return nil, fmt.Errorf("请求失败 (%dms): %w", requestTime, err)
	}

	return resp, nil
}
