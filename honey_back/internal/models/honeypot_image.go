package models

import "time"

// HoneypotImage 蜜罐镜像模板（schema 由 gorm AutoMigrate 依据本模型生成）
type HoneypotImage struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	Type        string    `gorm:"size:50;not null"`
	Name        string    `gorm:"size:100;not null;uniqueIndex"`
	DockerImage string    `gorm:"size:255;not null"`
	DefaultPort int       `gorm:"not null"`
	Description string    `gorm:"size:512;not null;default:''"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (HoneypotImage) TableName() string { return "honeypot_images" }
