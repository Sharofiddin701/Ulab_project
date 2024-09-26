package service

import (
	"e-commerce/pkg/logger"
	"e-commerce/storage"
)

type IServiceManager interface {
	Auth() authService
}

type Service struct {
	auth   authService
	logger logger.LoggerI
}

func New(storage storage.StorageI, log logger.LoggerI, redis storage.RedisI) Service {
	return Service{
		auth:   NewAuthService(storage, log, redis),
		logger: log,
	}
}

func (s Service) Auth() authService {
	return s.auth
}
