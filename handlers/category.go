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
// @Summary		Add Category
// @Description	Admin can add new categories for products
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			category	body	domain.Category	true	"category"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/add-category [POST]
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

// @Summary		Delete Category
// @Description	Admin can delete a category
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			id	query	string	true	"id"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/delete-category     [DELETE]
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

// @Summary		Update Category
// @Description	Admin can update name of a category into new name
// @Tags			Admin Category Management
// @Accept			json
// @Produce		    json
// @Param			set_new_name	body	models.SetNewName	true	"set new name"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/update-category     [PUT]
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
