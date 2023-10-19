package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
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
