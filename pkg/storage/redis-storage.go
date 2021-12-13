package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vsychov/go-rating-stars/pkg/config"
	"time"
)

// RedisStorage redis client
type RedisStorage struct {
	client *redis.Client
}

// CreateRedis create instance of RedisStorage
func CreateRedis(config config.Config) (storage StorageInterface, err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	storage = &RedisStorage{
		client: rdb,
	}

	err = rdb.Ping(context.TODO()).Err()
	if err != nil {
		return
	}

	return
}

func (store *RedisStorage) SetNX(key string, value float64, ttl time.Duration) (err error) {
	result := store.client.SetNX(context.TODO(), key, value, ttl).Val()
	if !result {
		return fmt.Errorf("unable to acquire lock")
	}

	return
}

func (store *RedisStorage) Get(key string) (value float64, err error) {
	value, err = store.client.Get(context.TODO(), key).Float64()
	if err != nil {
		return
	}

	return
}

func (store *RedisStorage) Incr(key string, value float64) (result float64, err error) {
	result = store.client.IncrByFloat(context.TODO(), key, value).Val()

	return
}

func (store *RedisStorage) Cron() error {
	return nil
}
