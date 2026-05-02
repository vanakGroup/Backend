package userdto

type OTPData struct {
	OTP  string `json:"otp"`
	Used bool   `json:"used"`
}
