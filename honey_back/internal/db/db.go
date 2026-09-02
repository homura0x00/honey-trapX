package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"honey_back/internal/config"
	"log"
)

type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewData(cfg *config.Config) (*Data, func(), error) {
	db, err := InitMySQL(cfg.Mysql)
	if err != nil {
		return nil, nil, fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := InitRedis(cfg.Redis)
	if err != nil {
		return nil, nil, fmt.Errorf("init redis: %w", err)
	}

	data := &Data{
		DB:  db,
		RDB: rdb,
	}

	cleanup := func() {
		log.Println("Starting to close database connections...")

		// 1. Cleaning MySQL
		if sqlDB, err := db.DB(); err == nil {
			if closeErr := sqlDB.Close(); closeErr != nil {
				log.Println("Failed to close database connections: ", closeErr)
			} else {
				log.Println("Closing database connections...")
			}
		} else {
			log.Println("Failed to close database connections: ", err)
		}

		// 2. Cleaning Redis
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis: %v\n", err)
		} else {
			log.Println("Redis connection closed successfully")
		}
	}

	return data, cleanup, nil
}
