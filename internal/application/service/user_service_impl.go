package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Vanak/Backend/bootstrap"
	userdto "github.com/Vanak/Backend/internal/application/dto/user"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/entities"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/logging/logtypes"
	"github.com/Vanak/Backend/internal/domain/ports"
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
}

func NewUserService(
	cfg *bootstrap.Config,
	otpService usecase.OtpService,
	userCacheRepository redis.UserCacheRepository,
	unitOfWork ports.UnitOfWork,
	jwtService usecase.JwtService,
) *UserService {
	logger := logging.NewLogger(cfg)

	return &UserService{
		logger:              logger,
		cfg:                 cfg,
		otpService:          otpService,
		userCacheRepository: userCacheRepository,
		unitOfWork:          unitOfWork,
		jwtService:          jwtService,
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
		walletRepo := rf.WalletRepository()

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

			userWallet := &entities.Wallet{
				Balance:  0,
				IsActive: false,
				UserID:   user.ID,
			}

			err = walletRepo.Create(userWallet)
			if err != nil {
				return err
			}

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

	if _, err := us.unitOfWork.Factory().UserRepository().FindUserById(info.UserID); err != nil {
		unauthorizedError := exception.NewUnauthorizedError(fmt.Sprintf("%s.%s", us.cfg.Constant.ErrorField.User, us.cfg.Constant.ErrorTag.Unauthorized))
		panic(unauthorizedError)
	}

	updateData := make(map[string]interface{})

	minInfo := 0
	if info.FirstName != nil {
		firtname := strings.TrimSpace(*info.FirstName)
		if firtname == "" {

			updateData[us.cfg.Constant.ProfileFilleds.Firstname] = ""
		} else {

			updateData[us.cfg.Constant.ProfileFilleds.Firstname] = &firtname
			minInfo++
		}

	}
	if info.LastName != nil {
		lastname := strings.TrimSpace(*info.LastName)
		if lastname == "" {
			updateData[us.cfg.Constant.ProfileFilleds.Lastname] = ""
		} else {
			updateData[us.cfg.Constant.ProfileFilleds.Lastname] = &lastname
			minInfo++
		}
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
				isDuplicated := us.unitOfWork.Factory().UserRepository().CheckNationalCodeExists(ctx, *info.NationalCode, info.UserID)
				if isDuplicated {
					alreadyExistsErrors.Add(us.cfg.Constant.ProfileFilleds.NationalCode, us.cfg.Constant.ErrorTag.AlreadyExist, us.cfg.Constant.ErrorField.NationalID)
				}
			}

			updateData[us.cfg.Constant.ProfileFilleds.NationalCode] = *info.NationalCode
			minInfo++
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
				isDuplicated := us.unitOfWork.Factory().UserRepository().CheckShabaExists(ctx, shaba, info.UserID)
				if isDuplicated {
					alreadyExistsErrors.Add(us.cfg.Constant.ProfileFilleds.ShabaNumber, us.cfg.Constant.ErrorTag.AlreadyExist, us.cfg.Constant.ErrorField.IBAN)
				}
			}
			updateData[us.cfg.Constant.ProfileFilleds.ShabaNumber] = shaba
			minInfo++
		}
	}

	if validationErrs.HasErrors() {
		panic(validationErrs)
	}

	if alreadyExistsErrors.HasErrors() {
		panic(alreadyExistsErrors)
	}

	if minInfo == 4 {
		updateData[us.cfg.Constant.ProfileFilleds.MinInfoCompleted] = true
	} else {
		updateData[us.cfg.Constant.ProfileFilleds.MinInfoCompleted] = false
	}

	err := errors.New("")
	if len(updateData) != 0 {
		err = us.unitOfWork.WithTransaction(func(rf ports.RepositoryFactory) error {
			userRepo := rf.UserRepository()
			walletRepo := rf.WalletRepository()
			err := userRepo.UpdateProfile(ctx, info.UserID, updateData)

			if err != nil {
				return err
			}

			if minInfo == 4 {
				err := walletRepo.ActiveWallet(info.UserID, true)
				if err != nil {
					return err
				}
			} else {
				err := walletRepo.ActiveWallet(info.UserID, false)
				if err != nil {
					return err
				}
			}

			return nil
		})

	}

	return err
}

func (us *UserService) GetInformation(userId uint) (*userdto.ProfileDTO, error) {
	user, err := us.unitOfWork.Factory().UserRepository().FindUserById(userId)
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
