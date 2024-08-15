package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Customer godoc
// @ID create_customers
// @Router /e_commerce/api/v1/customer [POST]
// @Summary Create Customer
// @Description Create Customer
// @Tags Customer
// @Accept json
// @Customery json
// @Param Customer body models.CustomerCreate true "CreateCustomerRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateCustomer(c *gin.Context) {

	var (
		customerCreate models.CustomerCreate
	)

	err := c.ShouldBindJSON(&customerCreate)
	if err != nil {
		h.logger.Error(err.Error() + " : " + "error Customer Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please,Enter Valid Data!")
		return
	}

	resp, err := h.storage.Customer().Create(c.Request.Context(), &customerCreate)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "Error Customer Create")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Creating Customer Successfully!!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Customer godoc
// @ID get_by_id_customer
// @Router /e_commerce/api/v1/customer/{id} [GET]
// @Summary Get By ID Customer
// @Description Get By ID Customer
// @Tags Customer
// @Accept json
// @Customer json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdCustomer(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Customer().GetByID(c.Request.Context(), &models.CustomerPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Customer.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Customer Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Customer godoc
// @ID get_list_customer
// @Router /e_commerce/api/v1/customer [GET]
// @Summary Get List Customer
// @Description Get List Customer
// @Tags Customer
// @Accept json
// @Customer json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCustomer(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCustomer INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCustomer INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Customer().GetList(c.Request.Context(), &models.CustomerGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Customer.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListCustomer Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Customer godoc
// @ID delete_customer
// @Router /e_commerce/api/v1/customer/{id} [DELETE]
// @Summary Delete Customer
// @Description Delete Customer
// @Tags Customer
// @Accept json
// @Customer json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteCustomer(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Customer().Delete(c.Request.Context(), &models.CustomerPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Customer.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Customer Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}

// Update Customer godoc
// @ID update_customer
// @Router /e_commerce/api/v1/customer/{id} [PUT]
// @Summary Update Customer
// @Description Update Customer
// @Tags Customer
// @Accept json
// @Customer json
// @Param id path string true "id"
// @Param Customer body models.CustomerUpdate true "UpdateCustomerRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateCustomer(c *gin.Context) {
	var (
		id             = c.Param("id")
		customerUpdate models.CustomerUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&customerUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Customer Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	customerUpdate.Id = id
	rowsAffected, err := h.storage.Customer().Update(c.Request.Context(), &customerUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Customer.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Customer.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Customer().GetByID(c.Request.Context(), &models.CustomerPrimaryKey{Id: customerUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Customer.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Customer Successfully!")
	c.JSON(http.StatusAccepted, resp)
}
