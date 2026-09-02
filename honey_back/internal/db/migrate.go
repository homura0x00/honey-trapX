package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"honey_back/internal/config"
	"honey_back/internal/models"
)

// EnsureDatabase 部署时用 flag 调用的建库指令：连到 MySQL 实例（不指定库），
// 若目标库不存在则创建。只在部署阶段手动执行一次。
func EnsureDatabase(cfg config.MySQL) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Config)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open mysql (no db): %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping mysql: %w", err)
	}

	name := strings.ReplaceAll(cfg.DBName, "`", "")
	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name))
	if err != nil {
		return fmt.Errorf("create database%s: %w", cfg.DBName, err)
	}
	return nil
}

// Migrate 用 gorm AutoMigrate 建表/补列。模型（internal/models）是 schema 事实源，
// 不再执行任何 SQL 迁移文件。
func Migrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&models.User{},
		&models.HoneypotImage{},
		&models.Deployment{},
	)
}
