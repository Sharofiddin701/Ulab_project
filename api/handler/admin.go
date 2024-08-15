package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Admin godoc
// @ID create_admins
// @Router /e_commerce/api/v1/admin [POST]
// @Summary Create Admin
// @Description Create Admin
// @Tags Admin
// @Accept json
// @Adminy json
// @Param Admin body models.AdminCreate true "CreateAdminRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateAdmin(c *gin.Context) {
	var (
		adminCreate models.AdminCreate
	)

	err := c.ShouldBindJSON(&adminCreate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Admin Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	resp, err := h.storage.Admin().Create(c.Request.Context(), &adminCreate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Admin.Create")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Create Admin Successfully!!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Admin godoc
// @ID get_by_id_admin
// @Router /e_commerce/api/v1/admin/{id} [GET]
// @Summary Get By ID Admin
// @Description Get By ID Admin
// @Tags Admin
// @Accept json
// @Admin json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdAdmin(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Admin().GetByID(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Admin.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Admin Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Admin godoc
// @ID get_list_admin
// @Router /e_commerce/api/v1/admin [GET]
// @Summary Get List Admin
// @Description Get List Admin
// @Tags Admin
// @Accept json
// @Admin json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListAdmin(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCategory INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListCategory INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Admin().GetList(c.Request.Context(), &models.AdminGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Admin.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListAdmin Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Admin godoc
// @ID delete_admin
// @Router /e_commerce/api/v1/admin/{id} [DELETE]
// @Summary Delete Admin
// @Description Delete Admin
// @Tags Admin
// @Accept json
// @Admin json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteAdmin(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Admin().Delete(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Admin.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Admin Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}

// Update Admin godoc
// @ID update_admin
// @Router /e_commerce/api/v1/admin/{id} [PUT]
// @Summary Update Admin
// @Description Update Admin
// @Tags Admin
// @Accept json
// @Admin json
// @Param id path string true "id"
// @Param Admin body models.AdminUpdate true "UpdateAdminRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateAdmin(c *gin.Context) {
	var (
		id          = c.Param("id")
		adminUpdate models.AdminUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&adminUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Admin Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	adminUpdate.Id = id
	rowsAffected, err := h.storage.Admin().Update(c.Request.Context(), &adminUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Admin.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Admin.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Admin().GetByID(c.Request.Context(), &models.AdminPrimaryKey{Id: adminUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Admin.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Admin Successfully!")
	c.JSON(http.StatusAccepted, resp)
}
