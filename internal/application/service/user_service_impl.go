package service

import (
	"context"
	"fmt"

	"github.com/Vanak/Backend/bootstrap"
	userdto "github.com/Vanak/Backend/internal/application/dto/user"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/domain/ports/repository/postgres"
	"github.com/Vanak/Backend/internal/domain/ports/repository/redis"
	"github.com/Vanak/Backend/internal/domain/utils"
)

type UserService struct {
	logger              logging.Logger
	otpService          usecase.OtpService
	cfg                 *bootstrap.Config
	userCacheRepository redis.UserCacheRepository
	unitOfWork          ports.UnitOfWork
	jwtService          usecase.JwtService
	userRepo            postgres.UserProfileRepository
}

func NewUserService(
	cfg *bootstrap.Config,
	otpService usecase.OtpService,
	userCacheRepository redis.UserCacheRepository,
	unitOfWork ports.UnitOfWork,
	jwtService usecase.JwtService,
	userRepo postgres.UserProfileRepository,
) *UserService {
	logger := logging.NewLogger(cfg)

	return &UserService{
		logger:              logger,
		cfg:                 cfg,
		otpService:          otpService,
		userCacheRepository: userCacheRepository,
		unitOfWork:          unitOfWork,
		jwtService:          jwtService,
		userRepo:            userRepo,
	}
}

func (us *UserService) SendOtp(phone string) error {
	otp, err := us.otpService.GenerateOtp()
	if err != nil {
		return err
	}

	err = us.otpService.SetOTp(context.Background(), phone, otp)

	if err != nil {
		return err
	}

	us.logger.Info(logtypes.Auth, logtypes.Otp, "otp generated and send", map[logtypes.ExtraKey]interface{}{
		"phone": phone,
		"otp":   otp,
	})

	return nil
}

func (us *UserService) VerifyOTP(verifyInfo userdto.VerifyPhoneOtp) error {
	err := us.otpService.ValidateOtp(context.Background(), verifyInfo.Phone, verifyInfo.Otp)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) Login(phone string) (userdto.UserLoginResponse, error) {
	var user *entities.User

	err := us.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
		userRepo := rf.UserRepository()

		var err error
		user, err = userRepo.FindUserByPhone(phone)

		if err != nil {
			return err
		}

		if user == nil {
			newUser := &entities.User{
				Phone: phone,
			}
			err := userRepo.Create(newUser)
			if err != nil {
				return err
			}
			user = newUser
		}

		return nil
	})

	if err != nil {
		return userdto.UserLoginResponse{}, err
	}

	accessToken, refreshToken, err := us.jwtService.GenerateToken(user.ID)

	if err != nil {
		return userdto.UserLoginResponse{}, err
	}

	return userdto.UserLoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Phone:        phone,
	}, nil
}

func (us *UserService) UpdateInformation(ctx context.Context, info userdto.UserFilled) error {
	alreadyExistsErrors := exception.NewAlreadyExistsError()
	validationErrs := exception.NewValidationErrors()

	updateData := make(map[string]interface{})

	if info.FirstName != nil {
		updateData[us.cfg.Constant.ProfileFilleds.Firstname] = *info.FirstName
	}
	if info.LastName != nil {
		updateData[us.cfg.Constant.ProfileFilleds.Lastname] = *info.LastName
	}
	if info.NationalCode != nil {
		code := *info.NationalCode

		if code == "" {
			updateData[us.cfg.Constant.ProfileFilleds.NationalCode] = nil
		} else {
			if len(code) != 10 {
				validationErrs.Add(us.cfg.Constant.ErrorField.NationalID, us.cfg.Constant.ErrorTag.NationalCodeLength)
			}

			if !utils.IsNumeric(code) {
				validationErrs.Add(us.cfg.Constant.ErrorField.NationalID, us.cfg.Constant.ErrorTag.Numeric)
			}

			if !validationErrs.HasErrors() {
				isDuplicated := us.userRepo.CheckNationalCodeExists(ctx, *info.NationalCode, info.UserID)
				if isDuplicated {
					alreadyExistsErrors.Add(us.cfg.Constant.ProfileFilleds.NationalCode, us.cfg.Constant.ErrorTag.AlreadyExist, us.cfg.Constant.ErrorField.NationalID)
				}
			}

			updateData[us.cfg.Constant.ProfileFilleds.NationalCode] = *info.NationalCode
		}
	}

	if info.ShabaNumber != nil {
		shaba := *info.ShabaNumber
		if shaba == "" {
			updateData[us.cfg.Constant.ProfileFilleds.ShabaNumber] = nil
		} else {
			if len(shaba) != 24 {
				validationErrs.Add(us.cfg.Constant.ErrorField.IBAN, us.cfg.Constant.ErrorTag.ShabaLength)
			}

			if !utils.IsNumeric(shaba) {
				validationErrs.Add(us.cfg.Constant.ErrorField.IBAN, us.cfg.Constant.ErrorTag.Numeric)
			}

			if !validationErrs.HasErrors() {
				isDuplicated := us.userRepo.CheckShabaExists(ctx, shaba, info.UserID)
				if isDuplicated {
					alreadyExistsErrors.Add(us.cfg.Constant.ProfileFilleds.ShabaNumber, us.cfg.Constant.ErrorTag.AlreadyExist, us.cfg.Constant.ErrorField.IBAN)
				}
			}
			updateData[us.cfg.Constant.ProfileFilleds.ShabaNumber] = shaba
		}
	}

	if validationErrs.HasErrors() {
		panic(validationErrs)
	}

	if alreadyExistsErrors.HasErrors() {
		panic(alreadyExistsErrors)
	}

	if _, err := us.userRepo.FindUserById(info.UserID); err != nil {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", us.cfg.Constant.ErrorField.User, us.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	if len(updateData) != 0 {
		if err := us.userRepo.UpdateProfile(ctx, info.UserID, updateData); err != nil {
			panic(err)
		}
	}

	return nil
}

func (us *UserService) GetInformation(userId uint) (*userdto.ProfileDTO, error) {
	user, err := us.userRepo.FindUserById(userId)
	if err != nil {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", us.cfg.Constant.ErrorField.User, us.cfg.Constant.ErrorTag.Unauthorized))
		return nil, unauthorizedError
	}

	userProfile := userdto.ProfileDTO{
		Phone:        user.Phone,
		FirstName:    utils.SafeString(user.FirstName),
		LastName:     utils.SafeString(user.LastName),
		NationalCode: utils.SafeString(user.NationalCode),
		ShabaNumber:  utils.SafeString(user.ShabaNumber),
	}

	return &userProfile, nil

}
