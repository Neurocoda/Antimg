package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Neurocoda/Antimg/config"
	"github.com/Neurocoda/Antimg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 生成Web Session Token（7天有效期）
func GenerateWebToken(user *models.User) (string, error) {
	claims := Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天有效期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// 生成API Token（永不过期）
func GenerateAPIToken(user *models.User) (string, error) {
	claims := Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// API Token 永不过期，不设置 ExpiresAt
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		// 尝试JWT Token验证（包括API Token和Web Token）
		// 先用MapClaims解析，因为API Token使用了MapClaims格式
		mapClaims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, mapClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 从MapClaims中提取用户信息
		username, ok := mapClaims["username"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		role, ok := mapClaims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token格式错误"})
			c.Abort()
			return
		}

		// 验证用户是否存在
		user, userErr := models.GetUserByUsername(username)
		if userErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 检查是否为当前有效的API Token
		if user.APIToken == tokenString {
			c.Set("username", username)
			c.Set("role", role)
			c.Set("auth_type", "api_token")
			c.Next()
			return
		}

		// 检查是否有过期时间（Web Token特征）
		if exp, hasExp := mapClaims["exp"]; hasExp {
			// 有过期时间，说明是Web Token，检查是否过期
			if expTime, ok := exp.(float64); ok {
				if time.Unix(int64(expTime), 0).Before(time.Now()) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "token已过期"})
					c.Abort()
					return
				}
				c.Set("username", username)
				c.Set("role", role)
				c.Set("auth_type", "web_token")
				c.Next()
				return
			}
		}

		// 没有过期时间但也不是当前API Token，说明是旧的API Token，拒绝访问
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API Token已失效，请使用最新的Token"})
		c.Abort()
		return
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func WebAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查cookie中的认证token
		token, err := c.Cookie("auth_token")
		if err != nil {
			// Cookie不存在，清除可能存在的无效cookie并重定向到登录页
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 验证JWT token
		claims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			// Token无效，清除cookie并重定向到登录页
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Token有效，设置用户信息到上下文
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

