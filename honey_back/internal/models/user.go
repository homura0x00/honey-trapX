package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// User 用户表（schema 由 gorm AutoMigrate 依据本模型生成）
type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	UserAccount  string         `gorm:"size:64;not null;uniqueIndex"`
	UserPassword string         `gorm:"size:512;not null"`
	UserName     string         `gorm:"size:64;not null;default:''"`
	UserRole     string         `gorm:"size:16;not null;default:'user'"`
	LastLoginAt  *time.Time     `gorm:"column:last_login_at"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (User) TableName() string { return "users" }
