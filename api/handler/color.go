package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Color godoc
// @ID create_color
// @Router /e_commerce/api/v1/color [POST]
// @Summary Create Color
// @Description Create Color
// @Tags Color
// @Accept json
// @Color json
// @Param Color body models.ColorCreate true "CreateColorRequest"
// @Success 200 {object} models.Color "Success Request"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateColor(c *gin.Context) {
	var colorCreate models.ColorCreate

	err := c.ShouldBindJSON(&colorCreate)
	if err != nil {
		h.logger.Error("Error while binding JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid request payload"})
		return
	}

	resp, err := h.storage.Color().Create(c.Request.Context(), &colorCreate)
	if err != nil {
		h.logger.Error("Error while creating color: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Failed to create color"})
		return
	}

	h.logger.Info("Color created successfully")
	c.JSON(http.StatusOK, resp)
}

// GetList Color godoc
// @ID get_list_color
// @Router /e_commerce/api/v1/color [GET]
// @Summary Get List Color
// @Description Get List Color
// @Tags Color
// @Accept json
// @Color json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListColor(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListColor INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListAdmin INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Color().GetList(c.Request.Context(), &models.ColorGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Color.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListColor Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Color godoc
// @ID delete_color
// @Router /e_commerce/api/v1/color/{id} [DELETE]
// @Summary Delete Color
// @Description Delete Color
// @Tags Color
// @Accept json
// @Color json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteColor(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Color().Delete(c.Request.Context(), &models.ColorPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Color.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Color Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}
