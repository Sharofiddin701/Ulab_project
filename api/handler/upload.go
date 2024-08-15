package handler

import (
	"e-commerce/pkg/helper"

	"github.com/gin-gonic/gin"
)

// upload Multiple Files godoc
// @ID upload_multiple_files
// @Router /e_commerce/api/v1/upload-files [POST]
// @Summary Upload Multiple Files
// @Description Upload Multiple Files
// @Tags Upload File
// @Accept multipart/form-data
// @Procedure json
// @Param file formData []file true "File to upload"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "File error")
		return
	}

	resp, err := helper.UploadFiles(form)
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "Upload error")
		return
	}

	c.JSON(200, resp)
}

// delete file godoc
// @ID delete_file
// @Router /e_commerce/api/v1/delete-file [DELETE]
// @Summary Delete File
// @Description Delete File
// @Tags Upload File
// @Accept multipart/form-data
// @Procedure json
// @Param id query string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteFile(c *gin.Context) {

	err := helper.DeleteFile(c.Query("id"))
	if err != nil {
		h.logger.Error(err.Error() + "  :  " + "Upload error")
		c.JSON(500, err.Error())
		return
	}
	c.JSON(204, "success")
}
