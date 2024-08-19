package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Brand godoc
// @ID create_brends
// @Router /e_commerce/api/v1/brand [POST]
// @Summary Create Brand
// @Description Create Brand
// @Tags Brand
// @Accept json
// @Brand json
// @Param Brand body models.BrandCreate true "CreateBrandRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateBrand(c *gin.Context) {

	var (
		brandCreate models.BrandCreate
	)

	err := c.ShouldBindJSON(&brandCreate)
	if err != nil {
		h.logger.Error(err.Error() + " : " + "error Brand Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please,Enter Valid Data!")
		return
	}

	resp, err := h.storage.Brand().Create(c.Request.Context(), &brandCreate)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "Error Brand Create")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Creating Brand Successfully!!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Brand godoc
// @ID get_by_id_brand
// @Router /e_commerce/api/v1/brand/{id} [GET]
// @Summary Get By ID Brand
// @Description Get By ID Brand
// @Tags Brand
// @Accept json
// @Brand json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdBrand(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Brand().GetByID(c.Request.Context(), &models.BrandPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Brand.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Brand Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Brand godoc
// @ID get_list_brand
// @Router /e_commerce/api/v1/brand [GET]
// @Summary Get List Brand
// @Description Get List Brand
// @Tags Brand
// @Accept json
// @Brand json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListBrand(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListBrandINVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListBrand INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Brand().GetList(c.Request.Context(), &models.BrandGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Brand.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListBrand Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Brand godoc
// @ID delete_brand
// @Router /e_commerce/api/v1/brand/{id} [DELETE]
// @Summary Delete Brand
// @Description Delete Brand
// @Tags Brand
// @Accept json
// @Brand json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteBrand(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Brand().Delete(c.Request.Context(), &models.BrandPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Brand.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Brand Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}

// Update Brand godoc
// @ID update_brand
// @Router /e_commerce/api/v1/brand/{id} [PUT]
// @Summary Update Brand
// @Description Update Brand
// @Tags Brand
// @Accept json
// @Brand json
// @Param id path string true "id"
// @Param Brand body models.BrandUpdate true "UpdateBrandRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateBrand(c *gin.Context) {
	var (
		id          = c.Param("id")
		brandUpdate models.BrandUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&brandUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Brand Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	brandUpdate.Id = id
	rowsAffected, err := h.storage.Brand().Update(c.Request.Context(), &brandUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Brand.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Brand.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Brand().GetByID(c.Request.Context(), &models.BrandPrimaryKey{Id: brandUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Brand.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Brand Successfully!")
	c.JSON(http.StatusAccepted, resp)
}
