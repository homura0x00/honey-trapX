package honey

import (
	"context"

	"gorm.io/gorm"

	"honey_back/internal/models"
)

// defaultImages 内置镜像模板 seed（公开镜像，先验证编排；真实蜜罐镜像留给 capture 切片）。
var defaultImages = []models.HoneypotImage{
	{Type: "ssh", Name: "SSH 蜜罐", DockerImage: "linuxserver/openssh-server:latest", DefaultPort: 2222, Description: "SSH 服务（示例模板）"},
	{Type: "wordpress", Name: "WordPress 蜜罐", DockerImage: "wordpress:latest", DefaultPort: 80, Description: "WordPress 站点（示例模板）"},
	{Type: "redis", Name: "Redis 蜜罐", DockerImage: "redis:7", DefaultPort: 6379, Description: "Redis 服务（示例模板）"},
}

// EnsureDefaultImages 表为空时写入内置镜像模板（幂等，启动时调用）。
func EnsureDefaultImages(ctx context.Context, db *gorm.DB) error {
	var cnt int64
	if err := db.WithContext(ctx).Model(&models.HoneypotImage{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	return db.WithContext(ctx).Create(&defaultImages).Error
}
