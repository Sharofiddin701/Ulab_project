package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Product godoc
// @ID create_order_product
// @Router /e_commerce/api/v1/order_product [POST]
// @Summary Create OrderProduct
// @Description Create OrderProduct
// @Tags OrderProduct
// @Accept json
// @Produce json
// @Param OrderProduct body models.OrderProductCreate true "CreateOrderProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateOrderProduct(c *gin.Context) {
	var orderproductCreate models.OrderProductCreate

	if err := c.ShouldBindJSON(&orderproductCreate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid input data",
		})
		return
	}

	resp, err := h.storage.OrderProduct().Create(c.Request.Context(), &orderproductCreate)
	if err != nil {
		h.logger.Error("error in OrderProduct.Create: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	h.logger.Info("OrderProduct Created Successfully!")
	c.JSON(http.StatusCreated, Response{
		Data: resp,
	})
}

// GetByID OrderProduct godoc
// @ID get_by_id_order_product
// @Router /e_commerce/api/v1/order_product/{id} [GET]
// @Summary Get By ID OrderProduct
// @Description Get By ID OrderProduct
// @Tags OrderProduct
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdOrderProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	product, err := h.storage.OrderProduct().GetByID(c.Request.Context(), &models.OrderProductPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error("error in OrderProduct.GetByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	if product == nil {
		h.logger.Warn("OrderProduct not found with ID: " + id)
		c.JSON(http.StatusNotFound, Response{
			Error: "Product not found",
		})
		return
	}

	h.logger.Info("OrderProduct Retrieved Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: product,
	})
}

// GetList OrderProducts godoc
// @ID get_list_order_products
// @Router /e_commerce/api/v1/order_product [GET]
// @Summary Get List OrderProducts
// @Description Get List OrderProducts
// @Tags OrderProduct
// @Accept json
// @Produce json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param category_id query string false "category_id"
// @Param brand_id query string false "brand_id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListOrderProducts(c *gin.Context) {
	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error("Invalid offset: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid offset",
		})
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error("Invalid limit: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid limit",
		})
		return
	}

	orderID := c.Query("order_id")
	productID := c.Query("product_id")

	resp, err := h.storage.OrderProduct().GetList(c.Request.Context(), &models.OrderProductGetListRequest{
		Offset:    offset,
		Limit:     limit,
		OrderId:   orderID,
		ProductId: productID,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error("error in OrderProduct.GetList: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	h.logger.Info("OrderProducts Retrieved Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: resp,
	})
}

// Update OrderProduct godoc
// @ID update_order_product
// @Router /e_commerce/api/v1/order_product/{id} [PUT]
// @Summary Update OrderProduct
// @Description Update OrderProduct
// @Tags OrderProduct
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param OrderProduct body models.OrderProductUpdate true "UpdateOrderProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateOrderProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	var orderproductUpdate models.OrderProductUpdate
	if err := c.ShouldBindJSON(&orderproductUpdate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid input data",
		})
		return
	}

	orderproductUpdate.Id = id

	rowsAffected, err := h.storage.OrderProduct().Update(c.Request.Context(), &orderproductUpdate)
	if err != nil {
		h.logger.Error("error in OrderProduct.Update: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("No rows affected in OrderProduct.Update")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Unable to update data. Please try again later!",
		})
		return
	}

	h.logger.Info("OrderProduct Updated Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: "OrderProduct updated successfully",
	})
}

// Delete OrderProduct godoc
// @ID delete_order_product
// @Router /e_commerce/api/v1/order_product/{id} [DELETE]
// @Summary Delete OrderProduct
// @Description Delete OrderProduct
// @Tags OrderProduct
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteOrderProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	if err := h.storage.OrderProduct().Delete(c.Request.Context(), &models.OrderProductPrimaryKey{Id: id}); err != nil {
		h.logger.Error("error in OrderProduct.Delete: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Unable to delete data. Please try again later!",
		})
		return
	}

	h.logger.Info("OrderProduct Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
