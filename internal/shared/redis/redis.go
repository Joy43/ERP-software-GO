package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/vidatechbd/assmi-super-shop-erp-backend/internal/config"
)

func NewClient(cfg config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0, 
	})

	//-------- Test connection------------
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return client, nil
}
