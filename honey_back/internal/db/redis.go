package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"honey_back/internal/config"
	"strconv"
)

func InitRedis(cfg config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.IP + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return rdb, nil
}
