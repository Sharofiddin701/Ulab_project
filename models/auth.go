package models

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
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

type UserRegisterConfRequest struct {
	MobilePhone string          `json:"phone_number"`
	Otp         string          `json:"sms_code"`
	Customer    *CustomerCreate `json:"customer"`
}
