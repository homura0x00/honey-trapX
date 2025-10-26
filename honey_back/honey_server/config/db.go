package config

import (
	"fmt"
	"honey_back/honey_server/global"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
	Pwd     string `yaml:"pwd"`
}

func initMysql() {
	host := AppConfig.Database.Mysql.Host
	port := AppConfig.Database.Mysql.Port
	user := AppConfig.Database.Mysql.User
	password := AppConfig.Database.Mysql.Password
	dbname := AppConfig.Database.Mysql.Dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[SERVER] MySQL is not connected! ", err)
	}

	global.MysqlDb = db
}

func initRedis() {
	// TODO 配置 Redis 连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     AppConfig.Database.Redis.Address,
		Password: AppConfig.Database.Redis.Pwd,
		DB:       0,
	})
	global.RedisDb = rdb
}
