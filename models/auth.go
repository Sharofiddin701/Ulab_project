package models

import "time"

type AuthRequest struct {
	Phone_number string `json:"phone_number" binding:"required,phone_number"`
	Password     string `json:"password" binding:"required,min=6,max=16"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type VerifyRequest struct {
	Phone_number string `json:"phone_number" binding:"required,phone_number"`
	Code         string `json:"code" binding:"required,min=6"`
}

type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	ID           int64     `json:"id"`
	Phone_number string    `json:"phone_number"`
	CreatedAt    time.Time `json:"created_at"`
}
