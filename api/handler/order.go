package handler

import (
	"context"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Order godoc
// @ID          create_order
// @Router      /e_commerce/api/v1/order [POST]
// @Summary     Create Order
// @Description Create Order
// @Tags        Order
// @Accept      json
// @Order       json
// @Param       Order body models.SwaggerOrderCreateRequest true "CreateOrderRequest"
// @Success     201 {object} Response{data=string} "Success Request"
// @Response    400 {object} Response{data=string} "Bad Request"
// @Failure     500 {object} Response{data=string} "Server error"
func (h *handler) CreateOrder(c *gin.Context) {
	var (
		request models.OrderCreateRequest
	)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("error reading body: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}
	h.logger.Info("Incoming JSON: " + string(body))

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error("error unmarshalling JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid JSON!"})
		return
	}

	if request.Order.CustomerId == "" {
		h.logger.Error("Customer ID is empty!")
		c.JSON(http.StatusBadRequest, Response{Data: "Customer ID is required!"})
		return
	}
	for _, item := range request.Items {
		if item.ProductId == "" {
			h.logger.Error("Product ID is empty for one of the items!")
			c.JSON(http.StatusBadRequest, Response{Data: "Product ID is required for each item!"})
			return
		}
	}

	order, err := h.storage.Order().CreateOrder(&request)
	if err != nil {
		h.logger.Error("error in Order.CreateOrder: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.logger.Info("Order Created Successfully!")
	c.JSON(http.StatusCreated, Response{Data: order})
}

// GetByID Order godoc
// @ID get_by_id_order
// @Router /e_commerce/api/v1/order/{id} [GET]
// @Summary Get By ID Order
// @Description Get By ID Order
// @Tags Order
// @Accept json
// @Order json
// @Param id path string true "id"
// @Success 200 {object} Response{data=models.Order} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid ID!"})
		return
	}

	order, err := h.storage.Order().GetOrder(id)
	if err != nil {
		h.logger.Error("error in Order.GetOrder: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	h.logger.Info("Order Retrieved Successfully!")
	c.JSON(http.StatusOK, Response{Data: order})
}

// GetList 		Order godoc
// @ID   		get_all_orders
// @Router      /e_commerce/api/v1/order [GET]
// @Summary     Get All Products
// @Description Retrieve all products
// @Tags     Order
// @Accept   json
// @Produce  json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success  200 {object} []models.OrderCreateRequest
// @Response  400 {object} Response{data=string} "Bad Request"
// @Failure  500 {object} Response{data=string} "Server error"
func (h *handler) GetAllOrders(c *gin.Context) {

	var req models.OrderGetListRequest

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetAllOrders INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, Response{Data: "INVALID OFFSET"})
		return
	}
	req.Offset = offset

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetAllOrders INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, Response{Data: "INVALID LIMIT"})
		return
	}
	req.Limit = limit

	orders, err := h.storage.Order().GetAll(context.Background(), &req)
	if err != nil {
		h.logger.Error(err.Error() + " : Error while getting all orders")
		c.JSON(http.StatusInternalServerError, Response{Data: "Error while getting all orders"})
		return
	}

	h.logger.Info("Orders retrieved successfully")
	c.JSON(http.StatusOK, orders)
}

// Update Order godoc
// @ID update_order
// @Router /e_commerce/api/v1/order/{id} [PUT]
// @Summary Update Order
// @Description Update Order
// @Tags Order
// @Accept json
// @Order json
// @Param id path string true "id"
// @Param Order body models.OrderUpdate true "UpdateOrderRequest"
// @Success 202 {object} Response{data=models.OrderUpdate} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid ID!"})
		return
	}

	var orderUpdate models.Order
	if err := c.ShouldBindJSON(&orderUpdate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Please, Enter Valid Data!"})
		return
	}

	orderUpdate.Id = id

	err := h.storage.Order().UpdateOrder(orderUpdate)
	if err != nil {
		h.logger.Error("error in Order.UpdateOrder: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Server Error!"})
		return
	}

	// if rowsAffected <= 0 {
	// 	h.logger.Error("No rows affected in Order.UpdateOrder")
	// 	c.JSON(http.StatusBadRequest, Response{Data: "Unable to update data. Please try again later!"})
	// 	return
	// }

	h.logger.Info("Order Updated Successfully!")
	c.JSON(http.StatusAccepted, Response{Data: orderUpdate})
}

// Delete Order godoc
// @ID delete_order
// @Router /e_commerce/api/v1/order/{id} [DELETE]
// @Summary Delete Order
// @Description Delete Order
// @Tags Order
// @Accept json
// @Order json
// @Param id path string true "id"
// @Success 204 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid ID!"})
		return
	}

	if err := h.storage.Order().DeleteOrder(id); err != nil {
		h.logger.Error("error in Order.DeleteOrder: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Unable to delete data, please try again later!"})
		return
	}

	h.logger.Info("Order Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
