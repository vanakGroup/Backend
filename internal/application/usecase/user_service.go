package usecase

type UserService interface {
	SendOtp(phone string) error
}
