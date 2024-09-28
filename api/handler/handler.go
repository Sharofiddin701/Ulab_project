package handler

import (
	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/pkg/logger"
	"e-commerce/service"
	"e-commerce/storage"
	"strconv"

	"github.com/gin-gonic/gin"
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

func handleResponseLog(c *gin.Context, log logger.LoggerI, msg string, statusCode int, data interface{}) {
	resp := models.Response{}

	if statusCode >= 100 && statusCode <= 199 {
		resp.Description = config.ERR_INFORMATION
	} else if statusCode >= 200 && statusCode <= 299 {
		resp.Description = config.SUCCESS
		log.Info("REQUEST SUCCEEDED", logger.Any("msg: ", msg), logger.Int("status: ", statusCode))

	} else if statusCode >= 300 && statusCode <= 399 {
		resp.Description = config.ERR_REDIRECTION
	} else if statusCode >= 400 && statusCode <= 499 {
		resp.Description = config.ERR_BADREQUEST
		log.Error("!!!!!!!! BAD REQUEST !!!!!!!!", logger.Any("error: ", msg), logger.Int("status: ", statusCode))
	} else {
		resp.Description = config.ERR_INTERNAL_SERVER
		log.Error("!!!!!!!! ERR_INTERNAL_SERVER !!!!!!!!", logger.Any("error: ", msg), logger.Int("status: ", statusCode))
	}

	resp.StatusCode = statusCode
	resp.Data = data

	c.JSON(resp.StatusCode, resp)
}
