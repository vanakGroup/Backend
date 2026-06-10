package usecase

import (
	"context"

	userdto "github.com/Vanak/Backend/internal/application/dto/user"
)

type UserService interface {
	SendOtp(phone string) error
	VerifyOTP(verifyInfo userdto.VerifyPhoneOtp) error
	Login(phone string) (userdto.UserLoginResponse, error)
}

type UserProfileUseCase interface {
	UpdateInformation(ctx context.Context, info userdto.UserFilled) error
	GetInformation(userId uint) (*userdto.ProfileDTO,error)
}
