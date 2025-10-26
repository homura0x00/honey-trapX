package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port string
	}
	Database struct {
		Mysql MySQLConfig
		Redis RedisConfig
	}
	Jwt struct {
		Secret string `yaml:"secret"`
	}
	Session struct {
		SessionCookieKey   string `yaml:"sessionCookieKey"`
		SessionRedisPrefix string `yaml:"sessionRedisPrefix"`
		SessionExpire      int    `yaml:"sessionExpire"`
	}
}

var AppConfig *Config // 全局指针

func InitConfig() {
	viper.SetConfigName("settings") // 文件名称
	viper.SetConfigType("yaml")     // 文件类型
	viper.AddConfigPath(".")        // 程序入口的当前地址，例如/project/cmd/v1/main.go 就是当前路径

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	initMysql()
}
