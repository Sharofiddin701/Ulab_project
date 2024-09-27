package models

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	PhoneNumber  string `json:"phone_number"`
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type UserRegisterRequest struct {
	MobilePhone string `json:"mobile_phone"`
}

type UserLoginByPhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type UserLoginByPhoneConfirmRequest struct {
	PhoneNumber string `json:"phone_number"`
	OtpCode     string `json:"otp_code"`
}

type UserRegisterConfRequest struct {
	Customer *CustomerCreate `json:"customer"`
}
