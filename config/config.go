package config

import (
	"os"
)

type Config struct {
	Port          string
	JWTSecret     string
	AdminUsername string
	AdminPassword string
}

var AppConfig *Config

func Init() {
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" || jwtSecret == "your-secret-key-change-in-production" {
		// Production environment must set secure JWT secret
		panic("Please set secure JWT_SECRET environment variable! Minimum 32 characters required")
	}
	if len(jwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters for security")
	}

	AppConfig = &Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     jwtSecret,
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "password"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}