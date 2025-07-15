package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供有效的token"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "Ghj20000107...1qaz@WSX" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token无效"})
			c.Abort()
			return
		}
		c.Next()

	}
}
