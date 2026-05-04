package service

import (
	"context"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	"github.com/Vanak/Backend/internal/domain/repository/redis"
)

type UserService struct {
	logger              logging.Logger
	otpService          usecase.OtpService
	cfg                 *bootstrap.Config
	userCacheRepository redis.UserCacheRepository
}

func NewUserService(
	cfg *bootstrap.Config,
	otpService usecase.OtpService,
	userCacheRepository redis.UserCacheRepository,
) *UserService {
	logger := logging.NewLogger(cfg)

	return &UserService{
		logger:              logger,
		cfg:                 cfg,
		otpService:          otpService,
		userCacheRepository: userCacheRepository,
	}
}

func (userService *UserService) SendOtp(phone string) error {
	otp, err := userService.otpService.GenerateOtp()
	if err != nil {
		return err
	}

	err = userService.userCacheRepository.Set(context.Background(), phone, otp, time.Duration(userService.cfg.Env.OTP.ExpiryMinute)*time.Minute)

	if err != nil {
		return err
	}

	userService.logger.Info(logtypes.Auth, logtypes.Otp, "otp generated", map[logtypes.ExtraKey]interface{}{
		"phone": phone,
		"otp":   otp,
	})

	return nil
}
