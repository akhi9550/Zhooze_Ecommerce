package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Get Products Details to users
// @Description Retrieve all product Details with pagination to users
// @Tags User Product
// @Accept json
// @Produce json
// @Param page query string false "Page number"
// @Param count query string false "Page Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/products     [GET]
func ShowAllProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
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

// @Summary Show Products of specified category
// @Description Show all the Products belonging to a specified category
// @Tags User Product
// @Accept json
// @Produce json
// @Param data body map[string]int true "Category IDs and quantities"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/products/filter [POST]
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


// @Summary Get Products Details
// @Description Retrieve all product Details
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query string false "Page number"
// @Param count query string false "Page Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/products   [GET]
func ShowAllProductsFromAdmin(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	countStr := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(countStr)
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

// @Summary Get Products Details
// @Tags Admin Product Management
// @Description Search for a product
// @Accept json
// @Produce json
// @Security Bearer
// @Param prefix body models.SearchItems  true "Product details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/products/search  [GET]
func  SearchProducts(c *gin.Context) {

	var prefix models.SearchItems

	if err := c.ShouldBindJSON(&prefix); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields are provided in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	productDetails, err := usecase.SearchProductsOnPrefix(prefix.ProductName)

	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not retrive products by prefix search", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrived all details", productDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Add Products
// @Description Add product from admin side
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body models.Product true "Product details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/products [POST]
func AddProducts(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := validator.New().Struct(product)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	if product.Stock < 1 {
		errs := response.ClientResponse(http.StatusBadRequest, "Invalid Stock", nil, nil)
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
// @Router /admin/products    [DELETE]
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
// @Router /admin/products    [PUT]
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

// @Summary Add Product Image
// @Description Add product Product from admin side
// @Tags Admin Product Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param product_id query int  true "Product_id"
// @Param files formData file true "Image file to upload" collectionFormat "multi"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/products/upload-image 	[POST]
func UploadImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Retrieving images from form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No files provided", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	for _, file := range files {
		err := usecase.UpdateProductImage(id, file)
		if err != nil {
			log.Println(err)
			errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change one or more images", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorRes)
			return
		}
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully changed images", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
