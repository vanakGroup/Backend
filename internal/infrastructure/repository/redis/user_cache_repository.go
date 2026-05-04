package redis

import (
	"context"
	"encoding/json"
	"time"

	userdto "github.com/Vanak/Backend/internal/application/dto/user"
	"github.com/Vanak/Backend/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	rdb database.Cache
}

func NewUserCacheRepository(rdb database.Cache) *UserCacheRepository {
	return &UserCacheRepository{
		rdb: rdb,
	}
}

func (userCache *UserCacheRepository) Set(ctx context.Context, key string, otp string, duration time.Duration, used bool) error {
	otpData := userdto.OTPData{
		OTP:  otp,
		Used: used,
	}

	v, err := json.Marshal(otpData)

	err = userCache.rdb.GetRedis().Set(ctx, key, string(v), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (userCache *UserCacheRepository) Get(ctx context.Context, key string) (*userdto.OTPData, error) {
	v, err := userCache.rdb.GetRedis().Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
	}

	var otpData userdto.OTPData

	err = json.Unmarshal([]byte(v), &otpData)

	if err != nil {
		return nil, err
	}

	return &otpData, nil
}
