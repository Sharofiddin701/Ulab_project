package models

type AdminLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	PhoneNumber  string `json:"phone_number"`
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthInfoAdmin struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type AdminRegisterRequest struct {
	MobilePhone string `json:"mobile_phone"`
}

type AdminLoginByPhoneRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type AdminLoginByPhoneConfirmRequest struct {
	PhoneNumber string `json:"phone_number"`
	OtpCode     string `json:"otp_code"`
}

type AdminRegisterConfRequest struct {
	Customer *CustomerCreate `json:"customer"`
}
