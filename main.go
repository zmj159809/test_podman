package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// 版本信息，由构建时通过 ldflags 注入
var (
	Version   = "dev"      // 版本号，默认为 dev
	BuildTime = "unknown"  // 构建时间
)

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

// Server 封装HTTP服务器及其附加功能
type Server struct {
	config    *Config
	server    *http.Server
	startTime time.Time
	logger    *log.Logger
}

// NewServer 创建新的服务器实例
func NewServer(config *Config) *Server {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)
	
	return &Server{
		config:    config,
		startTime: time.Now(),
		logger:    logger,
	}
}

// loadConfig 从环境变量加载配置，使用默认值
func loadConfig() *Config {
	config := &Config{
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

// setupRoutes 配置HTTP路由
func (s *Server) setupRoutes() {
	mux := http.NewServeMux()

	// 主要处理器，带日志中间件
	mux.HandleFunc("/", s.loggingMiddleware(s.helloHandler))
	
	// 健康检查端点
	mux.HandleFunc("/health", s.loggingMiddleware(s.healthHandler))
	
	// 就绪检查端点
	mux.HandleFunc("/ready", s.loggingMiddleware(s.readyHandler))
	
	// 指标端点
	mux.HandleFunc("/metrics", s.loggingMiddleware(s.metricsHandler))

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.IdleTimeout) * time.Second,
	}
}

// loggingMiddleware 添加请求日志记录
func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 创建自定义ResponseWriter来捕获状态码
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next(wrapper, r)
		
		duration := time.Since(start)
		s.logger.Printf("%s %s %d %v %s", 
			r.Method, r.URL.Path, wrapper.statusCode, duration, r.RemoteAddr)
	}
}

// responseWriter 包装http.ResponseWriter以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// helloHandler 处理主要端点请求
func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message":   "Hello World111!",
		"timestamp": time.Now(),
		"version":   Version,
		"build_time": BuildTime,
		"uptime":    time.Since(s.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// healthHandler 处理健康检查请求
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(s.startTime).String(),
		Version:   Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(health); err != nil {
		s.logger.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// readyHandler 处理就绪检查请求
func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 在这里添加任何就绪检查（数据库连接等）
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"status": "ready",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Printf("Error encoding ready response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// metricsHandler 处理基础指标请求
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := map[string]interface{}{
		"uptime_seconds": time.Since(s.startTime).Seconds(),
		"timestamp":      time.Now().Unix(),
		"version":        Version,
		"build_time":     BuildTime,
		"go_version":     fmt.Sprintf("%s", os.Getenv("GO_VERSION")),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		s.logger.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	s.setupRoutes()
	
	s.logger.Printf("Server starting on port %d", s.config.Port)
	s.logger.Printf("Configuration: %+v", s.config)
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed to start: %w", err)
	}
	
	return nil
}

// Stop 优雅地停止HTTP服务器
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Println("Shutting down server...")
	
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	
	s.logger.Println("Server stopped gracefully")
	return nil
}

func main() {
	// 加载配置
	config := loadConfig()
	
	// 创建服务器实例
	server := NewServer(config)
	
	// 设置优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	// 在协程中启动服务器
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	
	// 等待关闭信号
	<-quit
	
	// 创建带超时的关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 尝试优雅关闭
	if err := server.Stop(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	
	log.Println("Application exited successfully")
}
