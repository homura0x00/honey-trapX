package models

import "time"

// Deployment 状态常量
const (
	DeploymentCreating = "creating"
	DeploymentRunning  = "running"
	DeploymentStopped  = "stopped"
	DeploymentError    = "error"
)

// Deployment 蜜罐部署实例（v1 硬删除；schema 由 gorm AutoMigrate 依据本模型生成）
type Deployment struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	Name          string    `gorm:"size:100;not null"`
	UserID        int64     `gorm:"not null;index"`
	ImageID       int64     `gorm:"not null"`
	ContainerID   string    `gorm:"size:100;not null;default:''"`
	Status        string    `gorm:"size:16;not null;default:'creating'"`
	HostPort      int       `gorm:"not null;uniqueIndex"`
	ContainerPort int       `gorm:"not null"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (Deployment) TableName() string { return "deployments" }
