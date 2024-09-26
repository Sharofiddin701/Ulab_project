package handler

import (
	"e-commerce/models"
	"fmt"

	// check "food/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogin godoc
// @Router       /e_commerce/api/v1/login [POST]
// @Summary      Customer login
// @Description  Login to Voltify
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.UserLoginRequest true "login"
// @Success      201  {object}  models.UserLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *handler) UserLogin(c *gin.Context) {
	loginReq := models.UserLoginRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		h.logger.Error(err.Error() + ":" + "error while binding body")
		c.JSON(http.StatusBadRequest, "error while binding body")
		return
	}

	fmt.Println("loginReq: ", loginReq)

	//TODO: need validate login & password

	loginResp, err := h.service.Auth().UserLogin(c.Request.Context(), loginReq)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "error while login")
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	h.logger.Info("Successfully login")
	c.JSON(http.StatusOK, loginResp)

}

// UserRegister godoc
// @Router       /e_commerce/api/v1/sendcode [POST]
// @Summary      Customer register
// @Description  Registering to Voltify
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterRequest true "register"
// @Success      201  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *handler) UserRegister(c *gin.Context) {
	loginReq := models.UserRegisterRequest{}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		h.logger.Error(err.Error() + ":" + "error while binding body")
		c.JSON(http.StatusBadRequest, "error while binding body")
		return
	}
	fmt.Println("loginReq: ", loginReq)

	// if err := check.ValidateEmailAddress(loginReq.MobilePhone); err != nil {
	//  handleResponseLog(c, h.log, "error while validating email" + loginReq.MobilePhone, http.StatusBadRequest, err.Error())
	//  return
	// }

	// if err := check.CheckEmailExists(loginReq.MobilePhone); err != nil {
	//  handleResponseLog(c, h.log, "error email does not exist" + loginReq.MobilePhone, http.StatusBadRequest, err.Error())
	// }

	err := h.service.Auth().UserRegister(c.Request.Context(), loginReq)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "error while registering")
		c.JSON(http.StatusInternalServerError, "error while registering")
		return
	}

	h.logger.Info("Successfully registered")
	c.JSON(http.StatusOK, "Successfully registered")
}

// UserRegisterConfirm godoc
// @Router       /e_commerce/api/v1/verifycode [POST]
// @Summary      Customer register
// @Description  Registering to Voltify
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        register body models.UserRegisterConfRequest true "register"
// @Success      201  {object}  models.UserLoginResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h *handler) UserRegisterConfirm(c *gin.Context) {
	req := models.UserRegisterConfRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(err.Error() + ":" + "error while binding body")
		c.JSON(http.StatusBadRequest, "error while binding body")
		return
	}
	fmt.Println("req: ", req)

	//TODO: need validate login & password

	confResp, err := h.service.Auth().UserRegisterConfirm(c.Request.Context(), req)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "error while registering")
		c.JSON(http.StatusInternalServerError, "error while registering")
		return
	}
	h.logger.Info("Successfully login")
	c.JSON(http.StatusOK, confResp)

}
