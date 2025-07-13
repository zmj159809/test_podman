package types

import "time"

// Config 应用程序配置结构
type Config struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
	LogLevel     string `json:"log_level"`
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
}
