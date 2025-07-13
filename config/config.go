package config

import (
	"os"
	"strconv"

	"test_podman/types"
)

// LoadConfig 从环境变量加载配置，使用默认值
func LoadConfig() *types.Config {
	config := &types.Config{
		Port:         5667,
		ReadTimeout:  15,
		WriteTimeout: 15,
		IdleTimeout:  60,
		LogLevel:     "info",
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if readTimeout := os.Getenv("READ_TIMEOUT"); readTimeout != "" {
		if rt, err := strconv.Atoi(readTimeout); err == nil {
			config.ReadTimeout = rt
		}
	}

	if writeTimeout := os.Getenv("WRITE_TIMEOUT"); writeTimeout != "" {
		if wt, err := strconv.Atoi(writeTimeout); err == nil {
			config.WriteTimeout = wt
		}
	}

	if idleTimeout := os.Getenv("IDLE_TIMEOUT"); idleTimeout != "" {
		if it, err := strconv.Atoi(idleTimeout); err == nil {
			config.IdleTimeout = it
		}
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config
}
