package userdto

type OTPData struct {
	OTP  string `json:"otp"`
	Used bool   `json:"used"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Phone        string `json:"phone"`
}
