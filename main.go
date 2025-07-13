package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test_podman/config"
	"test_podman/server"
)

func main() {
	// 加载配置
	config := config.LoadConfig()
	
	// 创建服务器实例
	server := server.New(config)
	
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
