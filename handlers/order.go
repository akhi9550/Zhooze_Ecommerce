package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ApproveOrder(c *gin.Context) {
	orderId := c.Param("order_id")
	err := usecase.ApproveOrder(orderId)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't approve the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Approved Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}
func CancelOrderFromAdmin(c *gin.Context) {
	order_id := c.Param("order_id")
	err := usecase.CancelOrderFromAdmin(order_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't cancel the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
	}
	success := response.ClientResponse(http.StatusOK, "Order Cancel Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}
func GetAllOrderDetailsForAdmin(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	allOrderDetails, err := usecase.GetAllOrderDetailsForAdmin(page)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not retrieve order details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Details Retrieved successfully", allOrderDetails, nil)
	c.JSON(http.StatusOK, success)
}
func RefundUser(c *gin.Context) {
	OrderID := c.Param("order_id")
	err := usecase.RefundUser(OrderID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Refund was not possible", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Refunded the user", nil, nil)
	c.JSON(http.StatusOK, success)
}
func GetOrderDetails(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page count not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	id, _ := c.Get("user_id")
	UserID := id.(int)
	OrderDetails, err := usecase.GetOrderDetails(UserID, page, pageSize)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}
	fmt.Println("full order details is ", OrderDetails)

	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", OrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}
func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	id, _ := c.Get("user_id")
	userID := id.(int)
	err := usecase.CancelOrders(orderID, userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not cancel the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Cancel Successfull", nil, nil)
	c.JSON(http.StatusOK, success)
}
