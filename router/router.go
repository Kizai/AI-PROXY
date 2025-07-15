package router

import (
	"fmt"

	"AI-PROXY/controller"
	"AI-PROXY/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	fmt.Printf("=== 路由配置开始 ===\n")

	// 1. 首页
	r.GET("/", func(c *gin.Context) {
		fmt.Printf("匹配到首页路由: %s\n", c.Request.URL.Path)
		c.File("./web/index.html")
	})

	// 2. 代理转发路由（高优先级）
	fmt.Printf("注册代理转发路由: /:apiName/*path\n")
	r.Any("/:apiName/*path", func(c *gin.Context) {
		fmt.Printf("匹配到代理转发路由: %s\n", c.Request.URL.Path)
		controller.ForwardRequest(c)
	})

	// 3. 管理后台路由
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth())

	//API配置管理
	admin.GET("/api-config", controller.GetAllAPIConfigs)
	admin.GET("/api-config/:name", controller.GetAPIConfig)
	admin.POST("/api-config", controller.CreateAPIConfig)
	admin.PUT("/api-config/:name", controller.UpdateAPIConfig)
	admin.DELETE("/api-config/:name", controller.DeleteAPIConfig)
	admin.POST("/api-config/test", controller.TestAPIConfig)

	//日志管理
	admin.GET("/logs", controller.GetRequestLogs)
	admin.DELETE("/logs", controller.DeleteRequestLogs)
	admin.GET("/logs/export", controller.ExportRequestLogs)
	admin.POST("/logs/clear", controller.DeleteRequestLogs) // 清空日志使用POST方法

	//统计数据
	admin.GET("/stats", controller.GetStatisticsSummary) // 修改为GetStatisticsSummary
	admin.GET("/stats/realtime", controller.GetRealTimeStats)
	admin.GET("/stats/api-table", controller.GetAPIStatsTable)
	admin.GET("/stats/debug", controller.DebugStatistics) // 调试接口
	admin.GET("/stats/test", controller.TestStatistics)   // 测试统计更新接口

	// 4. 静态资源托管（低优先级）
	r.Static("/js", "./web/js")
	r.Static("/css", "./web/css")
	r.Static("/assets", "./web/assets")
	r.Static("/pages", "./web/pages") // 添加pages目录的静态文件服务
	r.StaticFile("/index.html", "./web/index.html")

	// 5. SPA兜底，支持前端路由刷新
	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("匹配到NoRoute: %s\n", c.Request.URL.Path)
		c.File("./web/index.html")
	})

	fmt.Printf("=== 路由配置完成 ===\n")

	return r
}
