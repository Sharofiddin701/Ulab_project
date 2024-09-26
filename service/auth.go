package service

import (
	"context"
	nmagap "e-commerce"
	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/pkg"
	"e-commerce/pkg/jwt"
	"e-commerce/pkg/logger"
	"e-commerce/storage"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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

	// Note: Password comparison is commented out. You may want to uncomment and implement this.
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

	msg := fmt.Sprintf("iBron ilovasi ro`yxatdan o`tish uchun tasdiqlash kodi: %v", otpCode)

	err := a.redis.SetX(ctx, loginRequest.MobilePhone, otpCode, time.Minute*5)
	if err != nil {
		a.log.Error("error while setting smsCode to redis user register", logger.Error(err))
		return err
	}
	err = nmagap.SendSms(loginRequest.MobilePhone, msg)
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

func (a authService) UserLoginByPhone(ctx context.Context, req models.UserLoginByPhoneRequest) error {
	otpCode := pkg.GenerateOTP()
	msg := fmt.Sprintf("iBron ilovasi ro'yxatdan o'tish uchun tasdiqlash kodi: %v", otpCode)

	err := a.redis.SetX(ctx, req.PhoneNumber, otpCode, time.Minute*5)
	if err != nil {
		a.log.Error("error while setting smsCode to redis for user login", logger.Error(err))
		return err
	}

	err = nmagap.SendSms(req.PhoneNumber, msg)
	if err != nil {
		a.log.Error("error while sending sms code for user login", logger.Error(err))
		return err
	}
	return nil
}

func (a authService) UserLoginByPhoneConfirm(ctx context.Context, req models.UserLoginByPhoneConfirmRequest) (models.UserLoginResponse, error) {
	resp := models.UserLoginResponse{}

	storedOTP, err := a.redis.Get(ctx, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			a.log.Error("OTP code not found or expired", logger.Error(err))
			return resp, errors.New("OTP kod topilmadi yoki muddati tugagan")
		}
		a.log.Error("error while getting OTP code from redis", logger.Error(err))
		return resp, errors.New("Tizim xatosi yuz berdi")
	}

	if req.OtpCode != storedOTP {
		a.log.Error("incorrect OTP code", logger.Error(errors.New("OTP code mismatch")))
		return resp, errors.New("Noto'g'ri OTP kod")
	}

	err = a.redis.Del(ctx, req.PhoneNumber)
	if err != nil {
		a.log.Error("error while deleting OTP from redis", logger.Error(err))
		return resp, err
	}
	user, err := a.storage.Customer().GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		a.log.Error("error while getting user by phone number", logger.Error(err))
		return resp, err
	}

	resp.PhoneNumber = req.PhoneNumber
	resp.ID = user.Id

	return resp, nil
}
