package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(cfg *bootstrap.Config) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Env.Redis.Host, cfg.Env.Redis.Port),
		Password:        cfg.Env.Redis.Password,
		DB:              0,
		DialTimeout:     cfg.Env.Redis.DialTimeout * time.Second,
		ReadTimeout:     cfg.Env.Redis.ReadTimeout * time.Second,
		WriteTimeout:    cfg.Env.Redis.WriteTimeout * time.Second,
		PoolSize:        cfg.Env.Redis.PoolSize,
		PoolTimeout:     cfg.Env.Redis.PoolTimeout,
		ConnMaxIdleTime: cfg.Env.Redis.ConnMaxIdleTime * time.Millisecond,
		ConnMaxLifetime: cfg.Env.Redis.ConnMaxLifetime * time.Millisecond,
	})
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	if redisClient != nil {
		redisClient.Close()
	}
}

func Set[T any](ctx context.Context, c *redis.Client, key string, value T, duration time.Duration) error {
	v, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.Set(ctx, key, v, duration).Err()
}

func Get[T any](ctx context.Context, c *redis.Client, key string) (T, error) {
	var dest = *new(T)
	v, err := c.Get(ctx, key).Result()

	if err != nil {
		return dest, err
	}

	err = json.Unmarshal([]byte(v), dest)

	if err != nil {
		return dest, err
	}

	return dest, nil
}
