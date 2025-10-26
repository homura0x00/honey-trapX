package main

import (
	"context"
	"honey_back/honey_server/config"
	"honey_back/honey_server/internal/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ctx = context.Background()

func main() {
	// 1. 初始化
	config.InitConfig()
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色
	//gin.DisableConsoleColor()
	// 记录到日志文件
	//f, _ := os.Create("gin_dev.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 启动 web 服务器
	r := routers.SetupRouter()
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":8123"
	}
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
