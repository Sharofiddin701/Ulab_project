package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Location godoc
// @ID create_location
// @Router /e_commerce/api/v1/location [POST]
// @Summary Create Location
// @Description Create Location
// @Tags Location
// @Accept json
// @Location json
// @Param Location body models.LocationCreate true "CreateLocationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateLocation(c *gin.Context) {
	var (
		locationCreate models.LocationCreate
	)

	err := c.ShouldBindJSON(&locationCreate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Location Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	resp, err := h.storage.Location().Create(c.Request.Context(), &locationCreate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Location.Create")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Create Location Successfully!!")
	c.JSON(http.StatusCreated, resp)
}

// GetByID Location godoc
// @ID get_by_id_location
// @Router /e_commerce/api/v1/location/{id} [GET]
// @Summary Get By ID Location
// @Description Get By ID Location
// @Tags Location
// @Accept json
// @Location json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdLocation(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Location().GetByID(c.Request.Context(), &models.LacationPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Location.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Location Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Location godoc
// @ID get_list_location
// @Router /e_commerce/api/v1/location [GET]
// @Summary Get List Location
// @Description Get List Location
// @Tags Location
// @Accept json
// @Location json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListLocation(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListLocation INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListLocation INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Location().GetList(c.Request.Context(), &models.LocationGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Location.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListLocation Response!")
	c.JSON(http.StatusOK, resp)
}

// Update Location godoc
// @ID update_location
// @Router /e_commerce/api/v1/location/{id} [PUT]
// @Summary Update Location
// @Description Update Location
// @Tags Location
// @Accept json
// @Location json
// @Param id path string true "id"
// @Param Location body models.LocationUpdate true "UpdateLocationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateLocation(c *gin.Context) {
	var (
		id             = c.Param("id")
		locationUpdate models.LocationUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&locationUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Location Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	locationUpdate.Id = id
	rowsAffected, err := h.storage.Location().Update(c.Request.Context(), &locationUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Location.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Location.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Location().GetByID(c.Request.Context(), &models.LacationPrimaryKey{Id: locationUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Location.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Location Successfully!")
	c.JSON(http.StatusAccepted, resp)
}

// Delete Location godoc
// @ID delete_location
// @Router /e_commerce/api/v1/location/{id} [DELETE]
// @Summary Delete Location
// @Description Delete Location
// @Tags Location
// @Accept json
// @Location json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteLocation(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Location().Delete(c.Request.Context(), &models.LacationPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Location.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Location Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
