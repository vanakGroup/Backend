package userdto

type OTPData struct {
	OTP  string `json:"otp"`
	Used bool   `json:"used"`
}

type UserLoginResponse struct {
	AccessToken  string
	RefreshToken string
	Phone        string
}
