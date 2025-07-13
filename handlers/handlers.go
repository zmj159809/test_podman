package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"test_podman/types"
	"test_podman/version"
)

// Handlers 包含所有HTTP处理器及其依赖
type Handlers struct {
	startTime time.Time
	logger    *log.Logger
}

// New 创建新的处理器实例
func New(startTime time.Time, logger *log.Logger) *Handlers {
	return &Handlers{
		startTime: startTime,
		logger:    logger,
	}
}

// Hello 处理主要端点请求
func (h *Handlers) Hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ver, buildTime := version.GetVersion()
	response := map[string]interface{}{
		"message":    "Hello World111!",
		"timestamp":  time.Now(),
		"version":    ver,
		"build_time": buildTime,
		"uptime":     time.Since(h.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// RootHandler 智能根路径处理器
func (h *Handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	// 只处理精确的根路径
	if r.URL.Path == "/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		ver, buildTime := version.GetVersion()
		response := map[string]interface{}{
			"service":    "Test Podman Web Server",
			"version":    ver,
			"build_time": buildTime,
			"endpoints": map[string]string{
				"hello":   "/hello",
				"health":  "/health",
				"ready":   "/ready",
				"metrics": "/metrics",
			},
			"uptime": time.Since(h.startTime).String(),
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.logger.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	}
	
	// 其他路径返回 404
	http.NotFound(w, r)
}

// Health 处理健康检查请求
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ver, _ := version.GetVersion()
	health := types.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Version:   ver,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(health); err != nil {
		h.logger.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Ready 处理就绪检查请求
func (h *Handlers) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 在这里添加任何就绪检查（数据库连接等）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Printf("Error encoding ready response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Metrics 处理基础指标请求
func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ver, buildTime := version.GetVersion()
	metrics := map[string]interface{}{
		"uptime_seconds": time.Since(h.startTime).Seconds(),
		"timestamp":      time.Now().Unix(),
		"version":        ver,
		"build_time":     buildTime,
		"go_version":     runtime.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		h.logger.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
