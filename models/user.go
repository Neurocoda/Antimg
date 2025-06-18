package models

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/Neurocoda/Antimg/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	APIToken  string    `json:"api_token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 简单的内存存储（生产环境应使用数据库）
var (
	users          = make(map[string]*User)
	userID    uint = 1
	userMutex sync.RWMutex
)

func init() {
	// 创建默认管理员用户，使用配置中的密码
	createDefaultAdmin()
}

func createDefaultAdmin() {
	// 确保配置已初始化
	if config.AppConfig == nil {
		config.Init()
	}

	adminPassword := config.AppConfig.AdminPassword
	if adminPassword == "" {
		adminPassword = "admin123" // 后备密码
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)

	// 创建用户对象
	user := &User{
		ID:        userID,
		Username:  config.AppConfig.AdminUsername,
		Password:  string(hashedPassword),
		Email:     "admin@example.com",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 生成API Token（需要用户信息）
	apiToken, _ := generateJWTAPIToken(user)
	user.APIToken = apiToken

	users[config.AppConfig.AdminUsername] = user
	userID++
}

func GetUserByUsername(username string) (*User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	user, exists := users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	// 返回用户副本，避免外部修改
	userCopy := *user
	return &userCopy, nil
}

func ValidatePassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// 生成随机字符串API Token（旧方法，已弃用）
func generateRandomAPIToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 生成JWT格式的API Token（永不过期）
func generateJWTAPIToken(user *User) (string, error) {
	// 添加随机数确保每次生成的token都不同
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	// 直接在这里实现JWT生成，避免循环导入
	claims := map[string]interface{}{
		"username": user.Username,
		"role":     user.Role,
		"iat":      time.Now().Unix(),
		"jti":      hex.EncodeToString(randomBytes), // JWT ID，确保唯一性
		// 不设置exp字段，表示永不过期
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// 签名token
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 重置用户的API Token
func ResetAPIToken(username string) (string, error) {
	userMutex.Lock()
	defer userMutex.Unlock()

	user, exists := users[username]
	if !exists {
		return "", ErrUserNotFound
	}

	// 生成新的JWT API Token
	newToken, err := generateJWTAPIToken(user)
	if err != nil {
		return "", err
	}

	user.APIToken = newToken
	user.UpdatedAt = time.Now()

	return newToken, nil
}

// 通过API Token获取用户
func GetUserByAPIToken(apiToken string) (*User, error) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	for _, user := range users {
		if user.APIToken == apiToken {
			// 返回用户副本，避免外部修改
			userCopy := *user
			return &userCopy, nil
		}
	}
	return nil, ErrUserNotFound
}
