package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Category godoc
// @ID create_categorys
// @Router /e_commerce/api/v1/category [POST]
// @Summary Create Category
// @Description Create Category
// @Tags Category
// @Accept json
// @Category json
// @Param Category body models.CategoryCreate true "CreateCategoryRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateCategory(c *gin.Context) {

	var (
		categoryCreate models.CategoryCreate
	)

	err := c.ShouldBindJSON(&categoryCreate)
	if err != nil {
		h.logger.Error(err.Error() + " : " + "error Category Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please,Enter Valid Data!")
		return
	}

	resp, err := h.storage.Category().Create(c.Request.Context(), &categoryCreate)
	if err != nil {
		h.logger.Error(err.Error() + ":" + "Error Category Create")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Creating Category Successfully!!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Category godoc
// @ID get_by_id_category
// @Router /e_commerce/api/v1/category/{id} [GET]
// @Summary Get By ID Category
// @Description Get By ID Category
// @Tags Category
// @Accept json
// @Category json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdCategory(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Category().GetByID(c.Request.Context(), &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Category.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Category Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Category godoc
// @ID get_list_category
// @Router /e_commerce/api/v1/category [GET]
// @Summary Get List Category
// @Description Get List Category
// @Tags Category
// @Accept json
// @Category json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCategory(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCategoryINVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCategory INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Category().GetList(c.Request.Context(), &models.CategoryGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Category.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListCategory Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Category godoc
// @ID delete_category
// @Router /e_commerce/api/v1/category/{id} [DELETE]
// @Summary Delete Category
// @Description Delete Category
// @Tags Category
// @Accept json
// @Category json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteCategory(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Category().Delete(c.Request.Context(), &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Category.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Category Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}

// Update Category godoc
// @ID update_category
// @Router /e_commerce/api/v1/category/{id} [PUT]
// @Summary Update Category
// @Description Update Category
// @Tags Category
// @Accept json
// @Category json
// @Param id path string true "id"
// @Param Category body models.CategoryUpdate true "UpdateCategoryRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateCategory(c *gin.Context) {
	var (
		id             = c.Param("id")
		categoryUpdate models.CategoryUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&categoryUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Category Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	categoryUpdate.Id = id
	rowsAffected, err := h.storage.Category().Update(c.Request.Context(), &categoryUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Category.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Category.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Category().GetByID(c.Request.Context(), &models.CategoryPrimaryKey{Id: categoryUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Category.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Category Successfully!")
	c.JSON(http.StatusAccepted, resp)
}
