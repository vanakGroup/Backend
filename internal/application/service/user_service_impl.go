package service

import (
	"context"

	"github.com/Vanak/Backend/bootstrap"
	userdto "github.com/Vanak/Backend/internal/application/dto/user"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	"github.com/Vanak/Backend/internal/domain/repository/postgres"
	"github.com/Vanak/Backend/internal/domain/repository/redis"
)

type UserService struct {
	logger              logging.Logger
	otpService          usecase.OtpService
	cfg                 *bootstrap.Config
	userCacheRepository redis.UserCacheRepository
	userRepository      postgres.UserRepository
	jwtService          usecase.JwtService
}

func NewUserService(
	cfg *bootstrap.Config,
	otpService usecase.OtpService,
	userCacheRepository redis.UserCacheRepository,
	userRepository postgres.UserRepository,
	jwtService usecase.JwtService,
) *UserService {
	logger := logging.NewLogger(cfg)

	return &UserService{
		logger:              logger,
		cfg:                 cfg,
		otpService:          otpService,
		userCacheRepository: userCacheRepository,
		userRepository:      userRepository,
		jwtService:          jwtService,
	}
}

func (userService *UserService) SendOtp(phone string) error {
	otp, err := userService.otpService.GenerateOtp()
	if err != nil {
		return err
	}

	err = userService.otpService.SetOTp(context.Background(), phone, otp)
	// Set(context.Background(), phone, otp, time.Duration(userService.cfg.Env.OTP.ExpiryMinute)*time.Minute, false)

	if err != nil {
		return err
	}

	userService.logger.Info(logtypes.Auth, logtypes.Otp, "otp generated and send", map[logtypes.ExtraKey]interface{}{
		"phone": phone,
		"otp":   otp,
	})

	return nil
}

func (userService *UserService) VerifyOTP(verifyInfo userdto.VerifyPhoneOtp) error {
	err := userService.otpService.ValidateOtp(context.Background(), verifyInfo.Phone, verifyInfo.Otp)

	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) Login(phone string) (userdto.UserLoginResponse, error) {
	user, err := userService.userRepository.FindUserByPhone(phone)

	if err != nil {
		return userdto.UserLoginResponse{}, err
	}

	if user == nil {
		newUser := &entities.User{
			Phone: phone,
		}
		err := userService.userRepository.Create(newUser)
		if err != nil {
			return userdto.UserLoginResponse{}, err
		}
		user = newUser
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.UserLoginResponse{}, err
	}

	return userdto.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Phone:        phone,
	}, nil
}
