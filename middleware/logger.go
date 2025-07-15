package middleware

import (
	"fmt"
	"time"

	"AI-PROXY/util"

	"github.com/gin-gonic/gin"
)

// 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 添加明显的调试输出
		fmt.Printf("=== 收到请求 ===\n")
		fmt.Printf("路径: %s\n", c.Request.URL.Path)
		fmt.Printf("方法: %s\n", c.Request.Method)
		fmt.Printf("用户代理: %s\n", c.Request.UserAgent())
		fmt.Printf("================\n")

		c.Next()
		duration := time.Since(start).Milliseconds()
		util.LogRequest(
			"",
			c.Request.URL.Path,
			c.Request.Method,
			c.Writer.Status(),
			duration,
			"",
		)
	}
}
