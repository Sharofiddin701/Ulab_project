package service

import (
	"context"
	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/pkg"
	"e-commerce/pkg/jwt"
	"e-commerce/pkg/logger"
	"errors"
	"fmt"

	// "food/pkg/password"

	"e-commerce/storage"
	"time"
)

type authService struct {
	storage storage.StorageI
	log     logger.LoggerI
	redis   storage.RedisI
}

func NewAuthService(storage storage.StorageI, log logger.LoggerI, redis storage.RedisI) authService {
	return authService{
		storage: storage,
		log:     log,
		redis:   redis,
	}
}

func (a authService) UserLogin(ctx context.Context, loginRequest models.UserLoginRequest) (models.UserLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
	user, err := a.storage.Customer().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("error while getting user credentials by login", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	// if err = password.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
	//  a.log.Error("error while comparing password", logger.Error(err))
	//  return models.UserLoginResponse{}, err
	// }

	m := make(map[interface{}]interface{})

	m["user_id"] = user.Id
	m["user_role"] = config.CUSTOMER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for user login", logger.Error(err))
		return models.UserLoginResponse{}, err
	}

	return models.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a authService) UserRegister(ctx context.Context, loginRequest models.UserRegisterRequest) error {
	fmt.Println(" loginRequest.Login: ", loginRequest.MobilePhone)

	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("iBron ilovasi ro‘yxatdan o‘tish uchun tasdiqlash kodi: %v", otpCode)

	err := a.redis.SetX(ctx, loginRequest.MobilePhone, otpCode, time.Minute*5)
	if err != nil {
		a.log.Error("error while setting smsCode to redis user register", logger.Error(err))
		return err
	}

	err = pkg.SendSms(loginRequest.MobilePhone, msg)
	if err != nil {
		a.log.Error("error while sending sms code to user register", logger.Error(err))
		return err
	}
	return nil
}

func (a authService) UserRegisterConfirm(ctx context.Context, req models.UserRegisterConfRequest) (models.UserLoginResponse, error) {
	resp := models.UserLoginResponse{}

	otp, err := a.redis.Get(ctx, req.MobilePhone)
	if err != nil {
		a.log.Error("error while getting sms code for user register confirm", logger.Error(err))
		return resp, err
	}
	if req.Otp != otp {
		a.log.Error("incorrect sms code for user register confirm", logger.Error(err))
		return resp, errors.New("incorrect otp code")
	}
	req.Customer.Phone_number = req.MobilePhone
	id, err := a.storage.Customer().Create(ctx, req.Customer)
	if err != nil {
		a.log.Error("error while creating user", logger.Error(err))
		return resp, err
	}
	var m = make(map[interface{}]interface{})

	m["user_id"] = id
	m["user_role"] = config.CUSTOMER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for user register confirm", logger.Error(err))
		return resp, err
	}
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken

	return resp, nil
}
