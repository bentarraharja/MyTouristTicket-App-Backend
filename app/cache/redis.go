package cache

import (
	"context"
	"my-tourist-ticket/app/configs"
	"time"

	"github.com/go-redis/redis/v8"
)

// NewRedis membuat instance Redis baru dengan konfigurasi tertentu
func NewRedis(cfg *configs.AppConfig) RedisInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_ADDR,
		Password: cfg.REDIS_PASSWORD,
		DB:       cfg.REDIS_DB,
	})

	return &Redis{client}
}

type RedisInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	Update(ctx context.Context, key string, val interface{}) error
	Del(ctx context.Context, key string) error
}

type Redis struct {
	client *redis.Client
}

// Get mengambil nilai dari Redis berdasarkan kunci
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set menetapkan nilai ke Redis dengan kunci tertentu dan opsi kedaluwarsa
func (r *Redis) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, val, expiration).Err()
}

// Update mengupdate nilai di Redis berdasarkan kunci
func (r *Redis) Update(ctx context.Context, key string, val interface{}) error {
	return r.client.Set(ctx, key, val, redis.KeepTTL).Err()
}

// Del menghapus nilai dari Redis berdasarkan kunci
func (r *Redis) Del(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
