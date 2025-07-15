package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("=== PANIC 捕获 ===")
				fmt.Printf("panic: %+v\n", err)
				fmt.Println(string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":  "服务器内部错误",
					"detail": err,
					"stack":  string(debug.Stack()),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
