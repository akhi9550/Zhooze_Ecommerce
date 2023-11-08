package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
	"fmt"
	"image"
	"log"
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
// @Param			image_id	query		string	true	"image-id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/image-crop    [POST]
func CropImage(c *gin.Context) {
	imageId := c.Query("image_id")
	imageID, err := strconv.Atoi(imageId)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error in string conversion", nil, err)
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	imageUrl, err := usecase.GetImageURL(imageID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "error in cropping", nil, err)
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	fmt.Println("ddddddddddd", imageUrl)
	resp, err := http.Get(imageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image"})
		return
	}
	defer resp.Body.Close()

	inputImage, str, err := image.Decode(resp.Body)
	fmt.Println(str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to decode image"})
		return
	}

	cropRect := image.Rect(100, 100, 400, 400)
	croppedImage := imaging.Crop(inputImage, cropRect)
	img := image.Image(croppedImage)
	err = imaging.Save(img, "./uploads/", imaging.JPEGQuality(80))
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	c.JSON(http.StatusOK, response.ClientResponse(http.StatusOK, "Image cropped and saved successfully", nil, nil))
}
