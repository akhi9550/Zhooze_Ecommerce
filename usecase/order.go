package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

func OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	fmt.Println("useriddddddddd", userID)
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = uint(userID)
	fmt.Println("ffffffffffff", orderBody.UserID)
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
	fmt.Println("🤣🤣", cartItems)
	var orderDetails domain.Order
	var orderItemDetails domain.OrderItem
	// add general order details - that is to be added to orders table
	orderDetails = helper.CopyOrderDetails(orderDetails, orderBody)
	// if orderBody.PaymentID == 2 {
	// 	orderDetails.PaymentStatus = "not paid"
	// 	orderDetails.ShipmentStatus = "pending"
	// }
	err = repository.CreateOrder(orderDetails)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	for _, c := range cartItems {
		// for each order save details of products and associated details and use order_id as foreign key (for each order multiple product will be there)
		orderItemDetails.OrderID = orderDetails.ID
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		orderItemDetails.TotalPrice = c.TotalPrice
		fmt.Println("..............", orderItemDetails.OrderID, orderItemDetails.ProductID, orderItemDetails.Quantity, orderItemDetails.TotalPrice)
		err := repository.AddOrderItems(orderItemDetails, orderDetails.UserID, c.ProductID, c.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}

	}
	orderSuccessResponse, err := repository.GetBriefOrderDetails(int(orderDetails.ID))
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	return orderSuccessResponse, nil

}

// func OrderItemsFromCart(userID, CartID, addressID, paymentID int) (domain.Order, error) {
// 	addressExist := repository.CheckAddressAvailabilityWithID(addressID, userID)
// 	if !addressExist {
// 		return domain.Order{}, errors.New("address doesn't exist")
// 	}
// 	paymentMethod := repository.GetPaymentId(paymentID)
// 	if !paymentMethod {
// 		return domain.Order{}, errors.New("paymentmethod doesn't exist")
// 	}
// 	cartExist := repository.CheckCartAvailabilityWithID(CartID, userID)
// 	if !cartExist {
// 		return domain.Order{}, errors.New("cart doesn't exist")
// 	}
// 	totlaAmount, err := repository.TotalAmountInCart(CartID)
// 	if err != nil {
// 		return domain.Order{}, nil
// 	}
// 	orderItems, err := repository.OrderItemsFromCart(CartID)
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	if err := repository.AddpaymentMethod(paymentID, orderItems.ID); err != nil {
// 		return domain.Order{}, err
// 	}
// 	if err := repository.AddAmountToOrder(totlaAmount, orderItems.ID); err != nil {
// 		return domain.Order{}, err
// 	}
// 	stock, err := repository.FindOrderStock(CartID)
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	body, err := repository.GetOrder(int(orderItems.ID))
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	productID, err := repository.FindProductFromCart(CartID)
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	err = repository.CartEmpty(CartID)
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	err = repository.ProductStockMinus(productID, stock)
// 	if err != nil {
// 		return domain.Order{}, err
// 	}
// 	return body, nil
// }

func GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := repository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func CancelOrders(orderId string, userId int) error {
	userTest, err := repository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := repository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = repository.CancelOrders(orderId)
	if err != nil {
		return err
	}
	err = repository.UpdateQuantityOfProduct(orderProductDetails)
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
