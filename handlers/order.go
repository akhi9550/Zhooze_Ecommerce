package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Approve Order
// @Description Approve Order from admin side which is in processing state
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param    id   query   string   true    "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order/approve [GET]
func ApproveOrder(c *gin.Context) {
	orderId := c.Query("id")
	err := usecase.ApproveOrder(orderId)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't approve the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Approved Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Cancel Order Admin
// @Description Cancel Order from admin side
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param id query string true "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order/cancel   [GET]
func CancelOrderFromAdmin(c *gin.Context) {
	order_id := c.Query("id")
	err := usecase.CancelOrderFromAdmin(order_id)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't cancel the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Order Cancel Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Get All order details for admin
// @Description Get all order details to the admin side
// @Tags Admin Order Management
// @Accept   json
// @Produce  json
// @Security Bearer
// @Param page query string true "Page number"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /admin/order   [GET]
func GetAllOrderDetailsForAdmin(c *gin.Context) {
	pageStr := c.Query("page")
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

// // @Summary Order Items From Cart
// // @Description Add cart to the order using  cart id
// // @Tags  User Order Management
// // @Accept json
// // @Produce json
// // @Security Bearer
// // @Param cart_id query int true "cart id"
// // @Param address_id query int true "address id"
// // @Param payment_id query int true "payment id"
// // @Success 200 {object} response.Response{}
// // @Failure 500 {object} response.Response{}
// // @Router /user/order    [POST]
// func OrderItemsFromCart(c *gin.Context) {
// 	id, _ := c.Get("user_id")
// 	cart_id := c.Query("cart_id")
// 	cartID, err := strconv.Atoi(cart_id)
// 	if err != nil {
// 		errRes := response.ClientResponse(http.StatusBadRequest, "error in cart id", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errRes)
// 		return
// 	}
// 	address_id := c.Query("address_id")
// 	addressID, err := strconv.Atoi(address_id)
// 	if err != nil {
// 		errRes := response.ClientResponse(http.StatusBadRequest, "error in address id", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errRes)
// 		return
// 	}
// 	payment_id := c.Query("payment_id")
// 	paymentID, err := strconv.Atoi(payment_id)
// 	if err != nil {
// 		errRes := response.ClientResponse(http.StatusBadRequest, "error in address id", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errRes)
// 		return
// 	}
// 	Order, err := usecase.OrderItemsFromCart(id.(int), cartID, addressID, paymentID)
// 	if err != nil {
// 		errRes := response.ClientResponse(http.StatusBadRequest, "error in ordering", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errRes)
// 		return
// 	}

// 	success := response.ClientResponse(http.StatusOK, "succesfully added order", Order, nil)
// 	c.JSON(http.StatusOK, success)
// }

// @Summary Order Items from cart
// @Description Order all products which is currently present inside  the cart
// @Tags User Order
// @Accept json
// @Produce json
// @Security Bearer
// @Param orderBody body models.OrderFromCart true "Order details"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order [post]
func OrderItemsFromCart(c *gin.Context) {

	id, _ := c.Get("user_id")
	userID := id.(int)

	var orderFromCart models.OrderFromCart
	if err := c.ShouldBindJSON(&orderFromCart); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "bad request", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	orderSuccessResponse, err := usecase.OrderItemsFromCart(orderFromCart, userID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not do the order", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully created the order", orderSuccessResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

// @Summary Get Order Details to user side
// @Description Get all order details done by user to user side
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query string true "Page"
// @Param count query string true "Count"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order/page   [GET]
func GetOrderDetails(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "page number not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("count", "10"))
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
	successRes := response.ClientResponse(http.StatusOK, "Full Order Details", OrderDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary Cancel order
// @Description Cancel order by the user using order ID
// @Tags User Order Management
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string true "Order ID"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/order   [PUT]
func CancelOrder(c *gin.Context) {
	orderID := c.Query("id")
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

// @Summary		Checkout section
// @Description	Add products to carts  for the purchase
// @Tags			User Order Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/checkout    [GET]
func CheckOut(c *gin.Context) {
	userID, _ := c.Get("user_id")
	checkoutDetails, err := usecase.Checkout(userID.(int))

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Checkout Page loaded successfully", checkoutDetails, nil)
	c.JSON(http.StatusOK, successRes)
}

// @Summary		Checkout section
// @Description	Add products to carts  for the purchase
// @Tags			User Order Management
// @Accept			json
// @Produce		    json
// @Param			order_id	query	string	true	"order id"
// @Param			payment	query	string	true	"payment"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/place-order     [GET]
func PlaceOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userId := userID.(int)
	orderId := c.Query("order_id")
	paymentMethod := c.Query("payment")
	fmt.Println("payment is ", paymentMethod, "order id is ", orderId)
	if paymentMethod == "cash_on_delivery" {
		Invoice, err := usecase.ExecutePurchaseCOD(userId, orderId)
		if err != nil {
			errorRes := response.ClientResponse(http.StatusInternalServerError, "error in making code ", nil, err.Error())
			c.JSON(http.StatusInternalServerError, errorRes)
			return
		}
		successRes := response.ClientResponse(http.StatusOK, "Placed Order with cash on delivery", Invoice, nil)
		c.JSON(http.StatusOK, successRes)
	}
}
