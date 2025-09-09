package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server struct {
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	}
	Upload struct {
		MaxFileSize  int64    `yaml:"max_file_size"`
		AllowedTypes []string `yaml:"allowed_types"`
		StoragePath  string   `yaml:"storage_path"`
		ChunkSize    int64    `yaml:"chunk_size"`
	}
	Auth struct {
		JWTSecret       string        `yaml:"jwt_secret"`
		TokenExpiration time.Duration `yaml:"token_expiration"`
	}
	RateLimit struct {
		RequestsPerMinute int `yaml:"requests_per_minute"`
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Set defaults
	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.ReadTimeout = getDurationEnv("READ_TIMEOUT", 30*time.Second)
	cfg.Server.WriteTimeout = getDurationEnv("WRITE_TIMEOUT", 30*time.Second)

	cfg.Upload.MaxFileSize = getInt64Env("MAX_FILE_SIZE", 25*1024*1024) // 25MB
	cfg.Upload.AllowedTypes = []string{"image/jpeg", "image/png", "application/pdf"}
	cfg.Upload.StoragePath = getEnv("STORAGE_PATH", "./storage")
	cfg.Upload.ChunkSize = getInt64Env("CHUNK_SIZE", 1024*1024) // 1MB

	cfg.Auth.JWTSecret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	cfg.Auth.TokenExpiration = getDurationEnv("TOKEN_EXPIRATION", 24*time.Hour)

	cfg.RateLimit.RequestsPerMinute = getIntEnv("RATE_LIMIT", 60)

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(cfg.Upload.StoragePath, 0755); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
