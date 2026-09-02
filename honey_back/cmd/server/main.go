package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"honey_back/internal/auth"
	"honey_back/internal/config"
	"honey_back/internal/db"
	"honey_back/internal/docker"
	"honey_back/internal/honey"
	"honey_back/internal/server"
)

func main() {
	createDB := flag.Bool("create-db", false, "仅创建数据库后退出（部署时先执行一次）")
	createAdmin := flag.Bool("create-admin", false, "创建默认管理员后退出（部署时执行一次）")
	flag.Parse()

	cfg, err := config.Load("")
	if err != nil {
		log.Fatal(err)
	}

	// 部署指令 1：建库（连 MySQL 实例，不指定库）
	if *createDB {
		if err := db.EnsureDatabase(cfg.Mysql); err != nil {
			log.Fatalf("create-db: %v", err)
		}
		log.Printf("数据库已就绪: %s", cfg.Mysql.DBName)
		return
	}

	data, cleanup, err := db.NewData(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	ctx := context.Background()

	// gorm AutoMigrate：模型即 schema，每次启动幂等
	if err := db.Migrate(data.DB); err != nil {
		log.Fatalf("migrate (gorm AutoMigrate): %v", err)
	}

	authSvc := auth.NewService(data.DB, data.RDB, cfg.Session)

	// 部署指令 2：创建默认管理员
	if *createAdmin {
		if err := authSvc.BootstrapAdmin(ctx, cfg.Admin); err != nil {
			log.Fatalf("create-admin: %v", err)
		}
		log.Printf("默认管理员已就绪: %s（密码见 settings.yaml admin 段，请尽快修改）", cfg.Admin.Account)
		return
	}

	// 空表时写入内置镜像模板
	if err := honey.EnsureDefaultImages(ctx, data.DB); err != nil {
		log.Fatalf("seed images: %v", err)
	}

	dockerMgr, err := docker.NewManager()
	if err != nil {
		log.Fatalf("连接 Docker 失败（请确认 Docker 已启动）: %v", err)
	}
	honeySvc := honey.NewService(data.DB, dockerMgr)

	engine := server.New(cfg, authSvc, honeySvc)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	go func() {
		log.Printf("honey_back listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown:", err)
	}
	log.Println("server exited")
}
