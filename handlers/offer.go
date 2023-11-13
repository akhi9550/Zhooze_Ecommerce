package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary Add  Product Offer
// @Description Add a new Offer for a product by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon body models.ProductOfferReceiver true "Add new Product Offer"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/product-offer [post]
func AddProdcutOffer(c *gin.Context) {

	var productOffer models.ProductOfferReceiver

	if err := c.ShouldBindJSON(&productOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(productOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = usecase.AddProductOffer(productOffer)

	if err != nil {
		if errors.Is(err, errors.New("the offer already exists")) {
			errRes := response.ClientResponse(http.StatusForbidden, "could not add offers", nil, err.Error())
			c.JSON(http.StatusForbidden, errRes)
			return
		}
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}

// @Summary Add  Category Offer
// @Description Add a new Offer for a Category by specifying a limit
// @Tags Admin Offer Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param coupon body models.CategoryOfferReceiver true "Add new Category Offer"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/offer/category-offer [post]
func AddCategoryOffer(c *gin.Context) {

	var categoryOffer models.CategoryOfferReceiver

	if err := c.ShouldBindJSON(&categoryOffer); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "request fields in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(categoryOffer)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = usecase.AddCategoryOffer(categoryOffer)

	if err != nil {
		if errors.Is(err,errors.New("the offer already exists")) {
			errRes := response.ClientResponse(http.StatusForbidden, "could not add offers", nil, err.Error())
			c.JSON(http.StatusForbidden, errRes)
			return
		}
		errRes := response.ClientResponse(http.StatusInternalServerError, "could not add offer", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added offer", nil, nil)
	c.JSON(http.StatusCreated, successRes)

}
