package service

import (
	"context"
	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/pkg/jwt"
	"e-commerce/pkg/logger"
	"e-commerce/storage"
	"fmt"
)

type authadminService struct {
	storage storage.StorageI
	log     logger.LoggerI
	redis   storage.RedisI
}

func NewAuthAdminService(storage storage.StorageI, log logger.LoggerI, redis storage.RedisI) authadminService {
	return authadminService{
		storage: storage,
		log:     log,
		redis:   redis,
	}
}

func (a authadminService) AdminLogin(ctx context.Context, loginRequest models.AdminLoginRequest) (models.AdminLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
	admin, err := a.storage.Admin().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("error while getting user credentials by login", logger.Error(err))
		return models.AdminLoginResponse{}, err
	}

	// Note: Password comparison is commented out. You may want to uncomment and implement this.
	// if err = password.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
	//  a.log.Error("error while comparing password", logger.Error(err))
	//  return models.UserLoginResponse{}, err
	// }

	m := make(map[interface{}]interface{})
	m["user_id"] = admin.Id
	m["user_role"] = config.ADMIN_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for user login", logger.Error(err))
		return models.AdminLoginResponse{}, err
	}

	return models.AdminLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           admin.Id,
		PhoneNumber:  admin.Phone_number,
	}, nil
}
