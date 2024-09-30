package service

import (
	"e-commerce/pkg/logger"
	"e-commerce/storage"
)

type IServiceManager interface {
	Auth() authService
	AuthAdmin() authadminService
}

type Service struct {
	auth      authService
	authAdmin authadminService
	logger    logger.LoggerI
}

func New(storage storage.StorageI, log logger.LoggerI, redis storage.RedisI) Service {
	return Service{
		auth:      NewAuthService(storage, log, redis),
		authAdmin: NewAuthAdminService(storage, log, redis),
		logger:    log,
	}
}

func (s Service) Auth() authService {
	return s.auth
}

func (s Service) AuthAdmin() authadminService {
	return s.authAdmin
}


