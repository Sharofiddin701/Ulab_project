package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Order godoc
// @ID create_order
// @Router /e_commerce/api/v1/order [POST]
// @Summary Create Order
// @Description Create Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Order body models.OrderCreate true "CreateOrderRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateOrder(c *gin.Context) {
	var orderCreate models.OrderCreate

	if err := c.ShouldBindJSON(&orderCreate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	resp, err := h.storage.Order().Create(c.Request.Context(), &orderCreate)
	if err != nil {
		h.logger.Error("error in Order.Create: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Order Created Successfully!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Order godoc
// @ID get_by_id_order
// @Router /e_commerce/api/v1/order/{id} [GET]
// @Summary Get By ID Order
// @Description Get By ID Order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("invalid UUID format!")
		c.JSON(http.StatusBadRequest, "Invalid ID!")
		return
	}

	order, err := h.storage.Order().GetByID(c.Request.Context(), &models.OrderPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error("error in Order.GetByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Order Retrieved Successfully!")
	c.JSON(http.StatusOK, order)
}

// GetList Orders godoc
// @ID get_list_orders
// @Router /e_commerce/api/v1/order [GET]
// @Summary Get List Orders
// @Description Get List Orders
// @Tags Order
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param customer_id query string false "customer_id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListOrders(c *gin.Context) {
	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error("Invalid offset: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error("Invalid limit: " + err.Error())
		c.JSON(http.StatusBadRequest, "Invalid limit")
		return
	}

	customerID := c.Query("customer_id")

	resp, err := h.storage.Order().GetList(c.Request.Context(), &models.OrderGetListRequest{
		Offset:     offset,
		Limit:      limit,
		CustomerId: customerID,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error("error in Order.GetList: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Orders Retrieved Successfully!")
	c.JSON(http.StatusOK, resp)
}

// Update Order godoc
// @ID update_order
// @Router /e_commerce/api/v1/order/{id} [PUT]
// @Summary Update Order
// @Description Update Order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param Order body models.OrderUpdate true "UpdateOrderRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, "Invalid ID!")
		return
	}

	var orderUpdate models.OrderUpdate
	if err := c.ShouldBindJSON(&orderUpdate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	orderUpdate.Id = id

	rowsAffected, err := h.storage.Order().Update(c.Request.Context(), &orderUpdate)
	if err != nil {
		h.logger.Error("error in Order.Update: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("No rows affected in Order.Update")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	h.logger.Info("Order Updated Successfully!")
	c.JSON(http.StatusAccepted, orderUpdate)
}

// Delete Order godoc
// @ID delete_order
// @Router /e_commerce/api/v1/order/{id} [DELETE]
// @Summary Delete Order
// @Description Delete Order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, "Invalid ID!")
		return
	}

	if err := h.storage.Order().Delete(c.Request.Context(), &models.OrderPrimaryKey{Id: id}); err != nil {
		h.logger.Error("error in Order.Delete: " + err.Error())
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Order Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
