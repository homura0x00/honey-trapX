package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load 从 yaml 读取配置并补齐默认值。
// path 为空时使用 ./settings.yaml（相对运行目录）。
func Load(path string) (*Config, error) {
	if path == "" {
		path = "settings.yaml"
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("could not load config file %s: %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	// 默认值
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 3000
	}
	if cfg.Server.Timeout == 0 {
		cfg.Server.Timeout = 30
	}
	if cfg.Mysql.Host == "" {
		cfg.Mysql.Host = "127.0.0.1"
	}
	if cfg.Mysql.Port == 0 {
		cfg.Mysql.Port = 3306
	}
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	if cfg.Session.CookieName == "" {
		cfg.Session.CookieName = "honey_session"
	}
	if cfg.Session.RedisPrefix == "" {
		cfg.Session.RedisPrefix = "session:"
	}
	if cfg.Session.ExpireHours == 0 {
		cfg.Session.ExpireHours = 24
	}
	if cfg.Admin.Account == "" {
		cfg.Admin.Account = "admin"
	}
	if cfg.Admin.Password == "" {
		cfg.Admin.Password = "admin123456"
	}
	if cfg.Admin.UserName == "" {
		cfg.Admin.UserName = "管理员"
	}
	return &cfg, nil
}
