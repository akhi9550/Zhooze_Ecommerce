package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	fmt.Println("qqqqqqqqqqq")
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
func GetUsers(c *gin.Context) {
	users, err := usecase.ShowAllUsers()
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all Users", users, nil)
	c.JSON(http.StatusOK, success)

}

func BlockUser(c *gin.Context) {
	id := c.Query("id")
	err := usecase.BlockedUser(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)

}
func UnBlockUser(c *gin.Context) {
	id := c.Query("id")
	err := usecase.UnBlockedUser(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)

}
