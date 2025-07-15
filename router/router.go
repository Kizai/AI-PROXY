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

	// 首页（普通用户）
	r.GET("/", func(c *gin.Context) {
		c.File("./web/home.html")
	})

	// 管理后台页面入口
	r.GET("/admin", func(c *gin.Context) {
		c.File("./web/admin.html")
	})
	r.GET("/admin/", func(c *gin.Context) {
		c.File("./web/admin.html")
	})

	// 静态资源托管（支持 / 和 /admin 下的资源访问）
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.Static("/assets", "./web/assets")
	r.Static("/pages", "./web/pages")
	r.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	r.Static("/admin/css", "./web/css")
	r.Static("/admin/js", "./web/js")
	r.Static("/admin/assets", "./web/assets")
	r.Static("/admin/pages", "./web/pages")

	// 管理后台接口路由
	admin := r.Group("/admin")
	admin.Use(middleware.AdminAuth())
	admin.GET("/api-config", controller.GetAllAPIConfigs)
	admin.GET("/api-config/:name", controller.GetAPIConfig)
	admin.POST("/api-config", controller.CreateAPIConfig)
	admin.PUT("/api-config/:name", controller.UpdateAPIConfig)
	admin.DELETE("/api-config/:name", controller.DeleteAPIConfig)
	admin.POST("/api-config/test", controller.TestAPIConfig)

	// 代理转发路由（必须放在最后）
	fmt.Printf("注册代理转发路由: /:apiName/*path\n")
	r.Any("/:apiName/*path", func(c *gin.Context) {
		fmt.Printf("匹配到代理转发路由: %s\n", c.Request.URL.Path)
		controller.ForwardRequest(c)
	})

	// SPA兜底，支持前端路由刷新
	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("匹配到NoRoute: %s\n", c.Request.URL.Path)
		c.File("./web/admin.html")
	})

	fmt.Printf("=== 路由配置完成 ===\n")

	return r
}
