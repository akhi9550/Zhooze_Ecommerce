package handlers

import (
	"Zhooze/domain"
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowAllProducts(c *gin.Context) {
	pageString := c.Param("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	products, err := usecase.ShowAllProducts(page, count)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all products", products, nil)
	c.JSON(http.StatusOK, success)
}
func FilterCategory(c *gin.Context) {
	var data map[string]int
	if err := c.ShouldBindJSON(&data); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	productCategory, err := usecase.FilterCategory(data)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve products by category", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully filtered the category", productCategory, nil)
	c.JSON(http.StatusOK, success)
}
func AllProducts(c *gin.Context) {
	products, err := usecase.SeeAllProducts()
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all products", products, nil)
	c.JSON(http.StatusOK, success)

}

func AddProducts(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	products, err := usecase.AddProducts(product)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not add the product", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully added products", products, nil)
	c.JSON(http.StatusOK, success)

}
func DeleteProducts(c *gin.Context) {
	id := c.Query("id")
	err := usecase.DeleteProducts(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "could not delete the specified products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, success)
}

func UpdateProduct(c *gin.Context) {

	var p models.ProductUpdate

	if err := c.BindJSON(&p); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	a, err := usecase.UpdateProduct(p.ProductId, p.Stock)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not update the product quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated the product quantity", a, nil)
	c.JSON(http.StatusOK, successRes)

}
