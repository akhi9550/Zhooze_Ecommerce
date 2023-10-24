package handlers

import (
	"Zhooze/domain"
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)
//admin
func AddCategory(c *gin.Context) {
	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	cate, err := usecase.AddCategory(category)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not add the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully added Category", cate, nil)
	c.JSON(http.StatusOK, success)

}
func DeleteCategory(c *gin.Context) {
	id := c.Query("id")
	err := usecase.DeleteCategory(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not delete the specified category", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully deleted the category", nil, nil)
	c.JSON(http.StatusOK, success)

}
func UpdateCategory(c *gin.Context) {
	var categoryUpdate models.SetNewName
	if err := c.ShouldBindJSON(&categoryUpdate); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	ok, err := usecase.UpdateCategory(categoryUpdate.Current, categoryUpdate.New)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "could not update the product", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully updated Category", ok, nil)
	c.JSON(http.StatusOK, success)

}
