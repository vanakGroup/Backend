package usecase

import "context"

type OtpService interface {
	GenerateOtp() (string, error)
	SetOTp(ctx context.Context, mobileNumber string, otp string) error
	ValidateOtp(ctx context.Context, mobileNumber string, otp string) error
}
