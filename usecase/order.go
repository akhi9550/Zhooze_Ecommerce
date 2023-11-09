package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
)

func OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderBody.UserID = uint(userID)
	cartExist, err := repository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := repository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}

	// get all items a slice of carts
	cartItems, err := repository.GetAllItemsFromCart(int(orderBody.UserID))
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	var orderDetails domain.Order
	var orderItemDetails domain.OrderItem
	// add general order details - that is to be added to orders table
	orderDetails = helper.CopyOrderDetails(orderDetails, orderBody)

	// get grand total iterating through each products in carts
	for _, c := range cartItems {
		orderDetails.GrandTotal += c.TotalPrice
	}
	discount_price, err := repository.GetCouponDiscountPrice(int(orderBody.UserID), orderDetails.GrandTotal)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	err = repository.UpdateCouponDetails(discount_price, orderDetails.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderDetails.FinalPrice = orderDetails.GrandTotal - discount_price
	if orderBody.PaymentID == 2 {
		orderDetails.PaymentStatus = "not paid"
		orderDetails.ShipmentStatus = "pending"
	}

	
	err = repository.CreateOrder(orderDetails)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	for _, c := range cartItems {
		// for each order save details of products and associated details and use order_id as foreign key ( for each order multiple product will be there)
		orderItemDetails.OrderID = orderDetails.OrderId
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = int(c.Quantity)
		orderItemDetails.TotalPrice = c.TotalPrice

		err := repository.AddOrderItems(orderItemDetails, orderDetails.UserID, c.ProductID, c.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}

	}

	err = repository.UpdateUsedOfferDetails(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := repository.GetBriefOrderDetails(orderDetails.OrderId)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	return orderSuccessResponse, nil

}
func GetOrderDetails(userID int, page int, count int) ([]models.CombinedOrderDetails, error) {
	OrderDetails, err := repository.GetOrderDetails(userID, page, count)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return OrderDetails, nil
}

func CancelOrders(orderID string, userID int) error {
	ok, err := repository.CheckOrderID(orderID)
	fmt.Println(ok)
	if !ok {
		return err
	}
	userTest, err := repository.UserOrderRelationship(orderID, userID)
	if err != nil {
		return err
	}
	if userTest != userID {
		return errors.New("the order is not comes by this user")
	}
	orderProductDetails, err := repository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "deliverd" {
		return errors.New("item already delivered, cannot cancel")
	}
	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in " + message + ", so no point in cancelling")
	}
	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = repository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	err = repository.UpdateStockOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	return nil
}
func Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := repository.GetAllAddresses(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := repository.GetAllPaymentOption()
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	cartItems, err := repository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := repository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil
}
func ExecutePurchaseCOD(userID int, orderID string) (models.Invoice, error) {
	ok, err := repository.CartExist(userID)
	if err != nil {
		return models.Invoice{}, err
	}
	if !ok {
		return models.Invoice{}, errors.New("cart doesn't exist")
	}
	err = repository.EmptyCart(userID)
	if err != nil {
		return models.Invoice{}, err
	}

	address, err := repository.GetAddressFromOrderId(orderID)
	if err != nil {
		return models.Invoice{}, err
	}
	orderDetails, err := repository.GetOrderDetailOfAproduct(orderID)
	if err != nil {
		return models.Invoice{}, err
	}

	Invoice := models.Invoice{
		OrderDetails: orderDetails,
		AddressInfo:  address,
	}
	return Invoice, nil

}
