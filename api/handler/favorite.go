package handler

import (
	"e-commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetList Favorite godoc
// @ID get_list_favorite
// @Router /e_commerce/api/v1/favorite [GET]
// @Summary Get List Favorite
// @Description Get List Favorite
// @Tags Favorite
// @Accept json
// @Favorite json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListFavorite(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListFavorite INVALID OFFSET!")
		c.JSON(http.StatusBadRequest, "INVALID OFFSET")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "GetListFavorite INVALID LIMIT!")
		c.JSON(http.StatusBadRequest, "INVALID LIMIT")
		return
	}

	resp, err := h.storage.Favorite().GetList(c.Request.Context(), &models.FavoriteGetListRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil && err.Error() != "no rows in result set" {
		h.logger.Error(err.Error() + "  :  " + "storage.Favorite.GetList!")
		c.JSON(http.StatusInternalServerError, "Server Error!")
		return
	}

	h.logger.Info("GetListFavorite Response!")
	c.JSON(http.StatusOK, resp)
}
