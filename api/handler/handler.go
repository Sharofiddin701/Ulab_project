package handler

import (
	"e-commerce/config"
	"e-commerce/pkg/logger"
	"e-commerce/service"
	"e-commerce/storage"
	"strconv"
)

type handler struct {
	cfg     *config.Config
	logger  logger.LoggerI
	storage storage.StorageI
	service service.IServiceManager
}

type Response struct {
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	Error       interface{} `json:"error"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

func NewHandler(cfg *config.Config, storage storage.StorageI, logger logger.LoggerI, service service.IServiceManager) *handler {
	return &handler{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
		service: service,
	}
}

func (h *handler) getOffsetQuery(offset string) (int, error) {

	if len(offset) <= 0 {
		return 0, nil
	}

	return strconv.Atoi(offset)
}

func (h *handler) getLimitQuery(limit string) (int, error) {

	if len(limit) <= 0 {
		return 10, nil
	}

	return strconv.Atoi(limit)
}
