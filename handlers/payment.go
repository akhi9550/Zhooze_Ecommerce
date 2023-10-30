package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MakePaymentRazorPay(c *gin.Context) {
	orderID := c.Query("id")
	userID := c.Query("user_id")
	user_Id, _ := strconv.Atoi(userID)

	orderDetail, razorID, err := usecase.MakePaymentRazorPay(orderID, user_Id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "could not generate order details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"final_price": orderDetail.FinalPrice * 100,
		"razor_id":    razorID,
		"user_id":     userID,
		"order_id":    orderDetail.OrderId,
		"user_name":   orderDetail.Firstname,
		"total":       int(orderDetail.FinalPrice),
	})
}
func VerifyPayment(c *gin.Context) {
	orderID := c.Query("id")
	paymentID := c.Query("payment_id")
	fmt.Println("😁😁😁😁", orderID, paymentID)

	err := usecase.SavePaymentDetails(paymentID, orderID)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "could not update payment details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully updated payment details", nil, nil)
	c.JSON(http.StatusOK, success)
}
