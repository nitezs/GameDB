package cache

import (
	"GameDB/internal/config"
	"GameDB/internal/log"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	db *redis.Client
}

var Redis RedisCache

func InitRedis() {
	Redis = RedisCache{}
	Redis.db = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DBIndex,
	})
	err := Redis.CheckConnection()
	if err != nil {
		log.Logger.Panic("Cannot connect to redis")
	}
	log.Logger.Info("Connected to redis")
}

func (r *RedisCache) CheckConnection() error {
	ctx := context.Background()
	result, err := r.db.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if result != "PONG" {
		return fmt.Errorf("unexpected response from Redis: %s", result)
	}
	return nil
}

func (r *RedisCache) Get(key string) (string, bool) {
	ctx := context.Background()
	value, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return value, true
}

func (r *RedisCache) Add(key string, value interface{}) error {
	ctx := context.Background()
	cmd := r.db.Set(ctx, key, value, 7*24*time.Hour)
	return cmd.Err()
}

func (r *RedisCache) AddWithExpire(key string, value interface{}, expire time.Duration) error {
	ctx := context.Background()
	cmd := r.db.Set(ctx, key, value, expire)
	return cmd.Err()
}
