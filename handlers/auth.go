package handlers

import (
	"net/http"

	"github.com/Neurocoda/Antimg/config"
	"github.com/Neurocoda/Antimg/middleware"
	"github.com/Neurocoda/Antimg/models"
	"github.com/Neurocoda/Antimg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// API登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if !models.ValidatePassword(user, req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateWebToken(user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "生成token失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Web登录页面
func (h *AuthHandler) LoginPage(c *gin.Context) {
	// 检查是否已登录，如果已登录直接跳转到管理页面
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		// 验证token是否有效
		claims := &middleware.Claims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err == nil && parsedToken.Valid && claims.Role == "admin" {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	c.HTML(http.StatusOK, "base.html", gin.H{
		"title": "登录 - Antimg",
		"page":  "login",
	})
}

// Web登录处理
func (h *AuthHandler) WebLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"page":  "login",
			"error": "用户名和密码不能为空",
			"title": "登录 - Antimg",
		})
		return
	}

	user, err := models.GetUserByUsername(username)
	if err != nil || !models.ValidatePassword(user, password) {
		c.HTML(http.StatusUnauthorized, "base.html", gin.H{
			"page":  "login",
			"error": "用户名或密码错误",
			"title": "登录 - Antimg",
		})
		return
	}

	token, err := middleware.GenerateWebToken(user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "base.html", gin.H{
			"page":  "login",
			"error": "登录失败",
			"title": "登录 - Antimg",
		})
		return
	}

	// 设置cookie - 7天有效期，确保会话维持
	// MaxAge: 7天 (秒), Path: "/", Domain: "", Secure: false (HTTP), HttpOnly: true
	c.SetCookie("auth_token", token, 7*24*3600, "/", "", false, true)

	// 管理员登录后直接进入图像处理工作台
	if user.Role == "admin" {
		// 构建基础URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host

		// 直接渲染工作台页面，避免重定向导致的cookie延迟问题
		c.HTML(http.StatusOK, "base.html", gin.H{
			"title":      "Antimg",
			"username":   user.Username,
			"token":      user.APIToken,
			"baseURL":    baseURL,
			"page":       "process",
			"isLoggedIn": true,
		})
	} else {
		c.HTML(http.StatusUnauthorized, "base.html", gin.H{
			"page":  "login",
			"error": "仅限管理员登录",
			"title": "登录 - Antimg",
		})
		return
	}
}

// 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

// 重置API Token
func (h *AuthHandler) ResetAPIToken(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "未找到用户信息")
		return
	}

	newToken, err := models.ResetAPIToken(username)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "重置API Token失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"api_token": newToken,
		"message":   "API Token已重置",
	})
}
