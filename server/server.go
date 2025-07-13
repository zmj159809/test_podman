package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"test_podman/handlers"
	"test_podman/middleware"
	"test_podman/types"
)

// Server 封装HTTP服务器及其附加功能
type Server struct {
	config    *types.Config
	server    *http.Server
	startTime time.Time
	logger    *log.Logger
	handlers  *handlers.Handlers
}

// New 创建新的服务器实例
func New(config *types.Config) *Server {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lshortfile)
	startTime := time.Now()

	return &Server{
		config:    config,
		startTime: startTime,
		logger:    logger,
		handlers:  handlers.New(startTime, logger),
	}
}

// setupRoutes 配置HTTP路由
func (s *Server) setupRoutes() {
	mux := http.NewServeMux()

	// 获取日志中间件
	logMiddleware := middleware.LoggingMiddleware(s.logger)

	// 主要处理器，带日志中间件
	mux.HandleFunc("/hello", logMiddleware(s.handlers.Hello))

	// 健康检查端点
	mux.HandleFunc("/health", logMiddleware(s.handlers.Health))

	// 就绪检查端点
	mux.HandleFunc("/ready", logMiddleware(s.handlers.Ready))

	// 指标端点
	mux.HandleFunc("/metrics", logMiddleware(s.handlers.Metrics))

	// 404 处理器
	mux.HandleFunc("/", logMiddleware(s.handlers.RootHandler))

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(s.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.IdleTimeout) * time.Second,
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
