package config

import (
	"os"
	"strconv"
)

const (
	defaultDBDSN         = "root:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	defaultRedisAddr     = "localhost:6379"
	defaultRedisPassword = ""
	defaultRedisDB       = 0
	defaultJWTSecret     = "imaichika_secret_key"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func GetDBDSN() string {
	return getEnv("GO_BLOG_DB_DSN", defaultDBDSN)
}

func GetRedisAddr() string {
	return getEnv("GO_BLOG_REDIS_ADDR", defaultRedisAddr)
}

func GetRedisPassword() string {
	return getEnv("GO_BLOG_REDIS_PASSWORD", defaultRedisPassword)
}

func GetRedisDB() int {
	return getEnvInt("GO_BLOG_REDIS_DB", defaultRedisDB)
}

func GetJWTSecret() string {
	return getEnv("GO_BLOG_JWT_SECRET", defaultJWTSecret)
}
