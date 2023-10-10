package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var adminDetails models.AdminLogin
	if err := c.ShouldBindJSON(&adminDetails); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	admin, err := usecase.LoginHandler(adminDetails)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Cannot authenticate user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Admin authenticated successfully", admin, nil)
	c.JSON(http.StatusOK, success)
}
func DashBoard(c *gin.Context) {
	adminDashboard, err := usecase.DashBoard()
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Dashboard could not be displayed", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Admin dashboard displayed", adminDashboard, nil)
	c.JSON(http.StatusOK, success)
}
