package models

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type Register struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	Exist       bool   `json:"exist"`
}

type UserPhoneNumber struct {
	PhoneNumber string `json:"phone_number"`
}
