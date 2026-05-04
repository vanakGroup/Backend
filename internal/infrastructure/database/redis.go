package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	GetRedis() *redis.Client
	CloseRedis()
}

type RedisDatabase struct {
	redisClient *redis.Client
}

func InitRedis(cfg *bootstrap.Config) RedisDatabase {
	addr := fmt.Sprintf("%s:%s", cfg.Env.Redis.Host, cfg.Env.Redis.Port)
	fmt.Println("🔍 Redis attempt address:", addr)
	redisClient := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%s", cfg.Env.Redis.Host, cfg.Env.Redis.Port),
		Password:        cfg.Env.Redis.Password,
		DB:              0,
		DialTimeout:     cfg.Env.Redis.DialTimeout * time.Second,
		ReadTimeout:     cfg.Env.Redis.ReadTimeout * time.Second,
		WriteTimeout:    cfg.Env.Redis.WriteTimeout * time.Second,
		PoolSize:        cfg.Env.Redis.PoolSize,
		PoolTimeout:     cfg.Env.Redis.PoolTimeout,
		ConnMaxIdleTime: cfg.Env.Redis.ConnMaxIdleTime * time.Second,
		ConnMaxLifetime: cfg.Env.Redis.ConnMaxLifetime * time.Second,
	})

	return RedisDatabase{redisClient: redisClient}
}

// func (rb *RedisDatabase) GetRedis() *redis.Client {
// 	return rb.redisClient
// }

func (rb *RedisDatabase) GetRedis() *redis.Client {
	// یک تست پینگ سریع برای اطمینان از زنده بودن اتصال
	err := rb.redisClient.Ping(context.Background()).Err()
	if err != nil {
		fmt.Printf("❌ Redis Connection Error: %v\n", err)
	} else {
		fmt.Println("✅ Redis is alive and connected!")
	}
	return rb.redisClient
}

func (rb *RedisDatabase) CloseRedis() {
	if rb.redisClient != nil {
		rb.redisClient.Close()
	}
}
