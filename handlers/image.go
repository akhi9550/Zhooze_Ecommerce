package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
	"image"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

// @Summary		Croping Added Images
// @Description	Croping Image
// @Tags			Image
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			id	query		string	true	"image-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/image-crop    [POST]
func CropImage(c *gin.Context) {
	imageId := c.Query("image_id")
	imageID, err := strconv.Atoi(imageId)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error in string conversion", nil, err)
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	imageUrl, err := usecase.CropImage(imageID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error in cropping", nil, err)
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	inputImage, err := imaging.Open(imageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	cropRect := image.Rect(100, 100, 400, 400)
	croppedImage := imaging.Crop(inputImage, cropRect)
	err = imaging.Save(croppedImage, imageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "Image cropped and saved successfully", nil, nil))
}
