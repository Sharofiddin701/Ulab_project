package models

type User struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}
