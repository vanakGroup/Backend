package service

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/repository/redis"
)

type OtpService struct {
	logger              logging.Logger
	cfg                 *bootstrap.Config
	UserCacheRepository redis.UserCacheRepository
}

func NewOtpService(cfg *bootstrap.Config, UserCacheRepository redis.UserCacheRepository) *OtpService {
	logger := logging.NewLogger(cfg)

	return &OtpService{cfg: cfg, logger: logger, UserCacheRepository: UserCacheRepository}
}

func (os *OtpService) GenerateOtp() (string, error) {
	length := os.cfg.Env.OTP.Length

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length-1)), nil)

	n, err := crand.Int(crand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return "", err
	}

	num := new(big.Int).Add(n, min)

	return fmt.Sprintf("%0*d", length, num), nil
}

func (os *OtpService) SetOTp(ctx context.Context, mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", os.cfg.Constant.RedisKey.OtpKey, mobileNumber)

	res, err := os.UserCacheRepository.Get(ctx, key)

	if err == nil && res != nil {
		if !res.Used {
			return &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", os.cfg.Constant.ErrorField.OTP, os.cfg.Constant.ErrorTag.AlreadyExist)}
		} else if res.Used {
			return &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", os.cfg.Constant.ErrorField.OTP, os.cfg.Constant.ErrorTag.AlreadyUsed)}
		}
	}

	err = os.UserCacheRepository.Set(ctx, key, otp, time.Duration(os.cfg.Env.OTP.ExpiryMinute)*time.Minute)
	return nil
}

func (os *OtpService) ValidateOtp(ctx context.Context, mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", os.cfg.Constant.RedisKey.OtpKey, mobileNumber)

	res, err := os.UserCacheRepository.Get(ctx, key)

	if err != nil {
		return err
	} else if res.Used {
		return &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", os.cfg.Constant.ErrorField.OTP, os.cfg.Constant.ErrorTag.AlreadyUsed)}
	} else if !res.Used && res.OTP != otp {
		return &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", os.cfg.Constant.ErrorField.OTP, os.cfg.Constant.ErrorTag.Invalid)}
	}
	// else if !res.Used && res.OTP == otp {
	// 	res.Used = true
	// 	err = os.UserCacheRepository.rdb.GetRedis().Set(ctx, key, res, time.Duration(os.cfg.Env.OTP.ExpiryMinute)*time.Minute).Err()
	// }

	return nil
}
