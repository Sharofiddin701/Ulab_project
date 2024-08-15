package handler

import (
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create Banner godoc
// @ID create_banner
// @Router /e_commerce/api/v1/banner [POST]
// @Summary Create Banner
// @Description Create Banner
// @Tags Banner
// @Accept json
// @Produce json
// @Param Banner body models.BannerCreate true "CreateBannerRequest"
// @Success 200 {object} models.Banner "Success Request"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateBanner(c *gin.Context) {
	var bannerCreate models.BannerCreate

	err := c.ShouldBindJSON(&bannerCreate)
	if err != nil {
		h.logger.Error("Error while binding JSON: " + err.Error())
		c.JSON(http.StatusBadRequest, Response{Data: "Invalid request payload"})
		return
	}

	resp, err := h.storage.Banner().Create(c.Request.Context(), &bannerCreate)
	if err != nil {
		h.logger.Error("Error while creating banner: " + err.Error())
		c.JSON(http.StatusInternalServerError, Response{Data: "Failed to create banner"})
		return
	}

	h.logger.Info("Banner created successfully")
	c.JSON(http.StatusOK, resp)
}

// GetByID Banner godoc
// @ID get_by_id_banner
// @Router /e_commerce/api/v1/banner/{id} [GET]
// @Summary Get By ID Banner
// @Description Get By ID Banner
// @Tags Banner
// @Accept json
// @Banner json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdBanner(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	request, err := h.storage.Banner().GetByID(c.Request.Context(), &models.BannerPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Banner.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetByID Banner Response!")
	c.JSON(http.StatusOK, request)
}

// GetList Banner godoc
// @ID get_list_banner
// @Router /e_commerce/api/v1/banner [GET]
// @Summary Get List Banner
// @Description Get List Banner
// @Tags Banner
// @Accept json
// @Banner json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListBanner(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListBannerINVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListBanner INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Banner().GetList(c.Request.Context(), &models.BannerGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Banner.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListBanner Response!")
	c.JSON(http.StatusOK, resp)
}

// Delete Banner godoc
// @ID delete_banner
// @Router /e_commerce/api/v1/banner/{id} [DELETE]
// @Summary Delete Banner
// @Description Delete Banner
// @Tags Banner
// @Accept json
// @Banner json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteBanner(c *gin.Context) {
	var id = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.logger.Error("is not valid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id!")
		return
	}

	err := h.storage.Banner().Delete(c.Request.Context(), &models.BannerPrimaryKey{Id: id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Banner.Delete!")
		c.JSON(http.StatusInternalServerError, "Unable to delete data, please try again later!")
		return
	}

	h.logger.Info("Banner Deleted Successfully!")
	c.JSON(http.StatusNoContent, nil)
}

// Update Banner godoc
// @ID update_banner
// @Router /e_commerce/api/v1/banner/{id} [PUT]
// @Summary Update Banner
// @Description Update Banner
// @Tags Banner
// @Accept json
// @Banner json
// @Param id path string true "id"
// @Param Banner body models.BannerUpdate true "UpdateBannerRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateBanner(c *gin.Context) {
	var (
		id           = c.Param("id")
		bannerUpdate models.BannerUpdate
	)

	if !helper.IsValidUUID(id) {
		h.logger.Error("is invalid uuid!")
		c.JSON(http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&bannerUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "error Banner Should Bind Json!")
		c.JSON(http.StatusBadRequest, "Please, Enter Valid Data!")
		return
	}

	bannerUpdate.Id = id
	rowsAffected, err := h.storage.Banner().Update(c.Request.Context(), &bannerUpdate)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Banner.Update!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	if rowsAffected <= 0 {
		h.logger.Error("storage.Banner.Update!")
		c.JSON(http.StatusBadRequest, "Unable to update data. Please try again later!")
		return
	}

	resp, err := h.storage.Banner().GetByID(c.Request.Context(), &models.BannerPrimaryKey{Id: bannerUpdate.Id})
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "storage.Banner.GetByID!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("Update Banner Successfully!")
	c.JSON(http.StatusAccepted, resp)
}
