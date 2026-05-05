package usecase

import userdto "github.com/Vanak/Backend/internal/application/dto/user"

type UserService interface {
	SendOtp(phone string) error
	VerifyOTP(verifyInfo userdto.VerifyPhoneOtp) error
}
