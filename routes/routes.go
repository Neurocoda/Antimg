package routes

import (
	"Antimg/config"
	"Antimg/handlers"
	"Antimg/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// 移除内存限制，允许处理大文件
	r.MaxMultipartMemory = 0

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// 加载HTML模板
	r.LoadHTMLGlob("templates/**/*.html")

	// 调试：打印模板加载情况
	log.Println("模板加载完成")

	// 创建处理器
	authHandler := handlers.NewAuthHandler()
	imageHandler := handlers.NewImageHandler()

	// 公开路由
	r.GET("/", func(c *gin.Context) {
		// 检查是否已登录，如果已登录直接跳转到管理页面
		token, err := c.Cookie("auth_token")
		if err == nil && token != "" {
			// 验证token是否有效
			claims := &middleware.Claims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})
			if err == nil && parsedToken.Valid && claims.Role == "admin" {
				c.Redirect(302, "/admin")
				return
			}
		}
		c.Redirect(302, "/login")
	})

	// 认证相关路由
	r.GET("/login", authHandler.LoginPage)
	r.POST("/login", authHandler.WebLogin)
	r.GET("/logout", authHandler.Logout)

	// API路由组
	api := r.Group("/api")
	{
		// 公开API
		api.POST("/login", authHandler.Login)

		// 需要认证的API
		apiAuth := api.Group("/")
		apiAuth.Use(middleware.AuthMiddleware())
		{
			// 图片处理API
			apiAuth.POST("/attack", imageHandler.AttackWatermark)
		}
	}

	// Web路由组 - 管理员登录后直接进入图像处理工作台
	web := r.Group("/")
	web.Use(middleware.WebAuthMiddleware())
	{
		// 图像处理工作台（管理员登录后的主界面）
		web.GET("/admin", imageHandler.ProcessPage)
		web.POST("/admin/process", imageHandler.WebProcessImage)
		web.POST("/admin/reset-api-token", authHandler.ResetAPIToken)
	}

	return r
}
