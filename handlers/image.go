package handlers

// import (
// 	"Zhooze/usecase"
// 	"Zhooze/utils/response"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // @Summary Get Products Details to users
// // @Description Retrieve product images
// // @Tags User Product
// // @Accept json
// // @Produce json
// // @Param product_id query string true "product_id"
// // @Success 200 {object} response.Response{}
// // @Failure 500 {object} response.Response{}
// // @Router /user/products/image  [GET]
// func ShowImages(c *gin.Context) {
// 	product_id := c.Query("product_id")
// 	productID, err := strconv.Atoi(product_id)
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusInternalServerError, "error in string conversion", nil, err)
// 		c.JSON(http.StatusInternalServerError, errs)
// 		return
// 	}
// 	image, err := usecase.ShowImages(productID)
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadGateway, "could not retrive images", nil, err.Error())
// 		c.JSON(http.StatusBadGateway, errs)
// 		return
// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully retrive images", image, nil)
// 	c.JSON(http.StatusOK, success)
// }
