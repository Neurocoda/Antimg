package routes

import (
	"net/http"
	"time"

	"github.com/Neurocoda/Antimg/config"
	"github.com/Neurocoda/Antimg/handlers"
	"github.com/Neurocoda/Antimg/middleware"
	"github.com/Neurocoda/Antimg/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	// Set reasonable memory limit: 100MB to prevent OOM attacks
	r.MaxMultipartMemory = 100 << 20 // 100MB

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

	// 公开路由 - 直接显示工作台界面
	r.GET("/", func(c *gin.Context) {
		// 检查是否已登录
		var username string
		var token string
		var isLoggedIn bool

		cookieToken, err := c.Cookie("auth_token")
		if err == nil && cookieToken != "" {
			// 验证token是否有效
			claims := &middleware.Claims{}
			parsedToken, err := jwt.ParseWithClaims(cookieToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.AppConfig.JWTSecret), nil
			})
			if err == nil && parsedToken.Valid && claims.Role == "admin" {
				username = claims.Username
				isLoggedIn = true
				// 获取用户的API Token
				if user, userErr := models.GetUserByUsername(username); userErr == nil {
					token = user.APIToken
				}
			}
		}

		// 构建基础URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host

		// 渲染工作台页面
		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":      "Antimg",
			"username":   username,
			"token":      token,
			"baseURL":    baseURL,
			"page":       "process",
			"isLoggedIn": isLoggedIn,
		})
	})

	// 认证相关路由
	r.GET("/login", authHandler.LoginPage)
	r.POST("/login", authHandler.WebLogin)
	r.GET("/logout", authHandler.Logout)

	// API路由组
	api := r.Group("/api")
	// 添加速率限制：每分钟最多30个请求
	api.Use(middleware.RateLimitMiddleware(30, time.Minute))
	{
		// 公开API
		api.POST("/login", authHandler.Login)

		// 需要认证的API
		apiAuth := api.Group("/")
		apiAuth.Use(middleware.AuthMiddleware())
		// 图片处理API添加更严格的速率限制：每分钟最多10个请求
		apiAuth.Use(middleware.RateLimitMiddleware(10, time.Minute))
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
