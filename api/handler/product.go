package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Create Product godoc
// @ID create_product
// @Router /e_commerce/api/v1/product [POST]
// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Accept json
// @Product json
// @Param Product body models.ProductCreate true "CreateProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateProduct(c *gin.Context) {
	var newProduct models.ProductCreate

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid input data",
		})
		return
	}

	resp, err := h.storage.Product().Create(c.Request.Context(), &newProduct)
	if err != nil {
		h.logger.Error("error in Product.Create: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	h.logger.Info("Product Created Successfully!")
	c.JSON(http.StatusCreated, Response{
		Data: resp,
	})
}

// GetByID Product godoc
// @ID get_by_id_product
// @Router /e_commerce/api/v1/product/{id} [GET]
// @Summary Get By ID Product
// @Description Get By ID Product
// @Tags Product
// @Accept json
// @Product json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	product, err := h.storage.Product().GetByID(c.Request.Context(), &models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error("error in Product.GetByID: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	if product == nil {
		h.logger.Warn("Product not found with ID: " + id)
		c.JSON(http.StatusNotFound, Response{
			Error: "Product not found",
		})
		return
	}

	h.logger.Info("Product Retrieved Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: product,
	})
}

// GetList Product godoc
// @ID get_list_product
// @Router /e_commerce/api/v1/product [GET]
// @Summary Get List Product
// @Description Get List Product
// @Tags Product
// @Accept json
// @Product json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param favorite query string false "favorite"
// @Success 200 {object} Response{data=models.ProductGetListResponse} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListProduct(c *gin.Context) {
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

	favoriteStr := c.Query("favorite")
	var favorite *bool
	if favoriteStr != "" {
		fav, err := strconv.ParseBool(favoriteStr)
		if err != nil {
			h.logger.Error("Invalid favorite value: " + err.Error())
			c.JSON(http.StatusBadRequest, Response{
				Error: "Invalid favorite value",
			})
			return
		}
		favorite = &fav
	}

	resp, err := h.storage.Product().GetList(c.Request.Context(), &models.ProductGetListRequest{
		Offset:   offset,
		Limit:    limit,
		Favorite: favorite,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error("error in Product.GetList: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	h.logger.Info("Products Retrieved Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: resp,
	})
}

// Update Product godoc
// @ID update_product
// @Router /e_commerce/api/v1/product/{id} [PUT]
// @Summary Update Product
// @Description Update Product
// @Tags Product
// @Accept json
// @Product json
// @Param id path string true "id"
// @Param Product body models.ProductUpdate true "UpdateProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	var productUpdate models.ProductUpdate
	if err := c.ShouldBindJSON(&productUpdate); err != nil {
		h.logger.Error("error in ShouldBindJSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid input data",
		})
		return
	}

	productUpdate.Id = id

	rowsAffected, err := h.storage.Product().Update(c.Request.Context(), &productUpdate)
	if err != nil {
		h.logger.Error("error in Product.Update: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Server error",
		})
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("No rows affected in Product.Update")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Unable to update data. Please try again later!",
		})
		return
	}

	h.logger.Info("Product Updated Successfully!")
	c.JSON(http.StatusOK, Response{
		Data: "Product updated successfully",
	})
}

// Delete Product godoc
// @ID delete_product
// @Router /e_commerce/api/v1/product/{id} [DELETE]
// @Summary Delete Product
// @Description Delete Product
// @Tags Product
// @Accept json
// @Product json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("Invalid UUID format!")
		c.JSON(http.StatusBadRequest, Response{
			Error: "Invalid ID format",
		})
		return
	}

	if err := h.storage.Product().Delete(c.Request.Context(), &models.ProductPrimaryKey{Id: id}); err != nil {
		h.logger.Error("error in Product.Delete: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{
			Error: "Unable to delete data. Please try again later!",
		})
		return
	}

	h.logger.Info("NewProduct Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
