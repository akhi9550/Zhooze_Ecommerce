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

// @Summary Get Products Details to users
// @Description Retrieve all product Details with pagination to users
// @Tags User Product
// @Accept json
// @Produce json
// @Param page query string true "Page number"
// @Param count query string true "Page Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /page   [GET]
func ShowAllProducts(c *gin.Context) {
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Page count not in right format ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	products, err := usecase.ShowAllProducts(page, count)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't retrieve products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all products", products, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Get Products Details to users category based
// @Description Retrieve all product Details with pagination to users
// @Tags User Product
// @Accept json
// @Produce json
// @Param data body map[string]int true "Category IDs and quantities"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /filter   [POST]
func FilterCategory(c *gin.Context) {
	var data map[string]int
	if err := c.ShouldBindJSON(&data); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	productCategory, err := usecase.FilterCategory(data)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't retrieve products by category", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully filtered the category", productCategory, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Get Products Details to users
// @Description Retrieve all product Details with pagination to users
// @Tags User Product
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /products   [GET]
func AllProducts(c *gin.Context) {
	products, err := usecase.SeeAllProducts()
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't retrieve products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all products", products, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary Get Products Details to users
// @Description Retrieve all product Details
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query string true "Page number"
// @Param count query string true "Page Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /product   [GET]
func ShowAllProductsFromAdmin(c *gin.Context) {
	pageString := c.Query("page")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Page count not in right format ", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	products, err := usecase.ShowAllProductsFromAdmin(page, count)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't retrieve products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all products", products, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Add Products
// @Description Add product from admin side
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body domain.Product true "Product details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /product [POST]
func AddProducts(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
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

// @Summary Delete product
// @Description Delete a product from the admin side
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string true "product id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /product    [DELETE]
func DeleteProducts(c *gin.Context) {
	id := c.Query("id")
	err := usecase.DeleteProducts(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not delete the specified products", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Update Products quantity
// @Description Update quantity of already existing product
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param productUpdate body models.ProductUpdate true "Product details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /product [PATCH]
func UpdateProduct(c *gin.Context) {
	var p models.ProductUpdate
	if err := c.BindJSON(&p); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	a, err := usecase.UpdateProduct(p.ProductId, p.Stock)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the product quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated the product quantity", a, nil)
	c.JSON(http.StatusOK, successRes)

}
