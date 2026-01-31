package config

import (
	"os"
)

type Config struct {
	RedisAddr string
	RedisPass string
}

func LoadConfig() *Config {
	return &Config{
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASS", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
