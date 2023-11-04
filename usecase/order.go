package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
)

func OrderItemsFromCart(userID, CartID, addressID, paymentID int) (domain.Order, error) {
	addressExist := repository.CheckAddressAvailabilityWithID(addressID, userID)
	if !addressExist {
		return domain.Order{}, errors.New("address doesn't exist")
	}
	paymentMethod := repository.GetPaymentId(paymentID)
	if !paymentMethod {
		return domain.Order{}, errors.New("paymentmethod doesn't exist")
	}
	cartExist := repository.CheckCartAvailabilityWithID(CartID, userID)
	if !cartExist {
		return domain.Order{}, errors.New("cart doesn't exist")
	}
	totlaAmount, err := repository.TotalAmountInCart(CartID)
	if err != nil {
		return domain.Order{}, nil
	}
	orderItems, err := repository.OrderItemsFromCart(CartID, addressID, paymentID)
	if err != nil {
		return domain.Order{}, err
	}
	if err := repository.AddAmountToOrder(totlaAmount, orderItems.ID); err != nil {
		return domain.Order{}, err
	}
	body, err := repository.GetOrder(int(orderItems.ID))
	if err != nil {
		return domain.Order{}, err
	}
	return body, nil
}

func GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error) {
	OrderDetails, err := repository.GetOrderDetails(userID, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
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
		return errors.New("the order is not come by this user")
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
		return errors.New("the order is in" + message + ", so no point in cancelling")
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

	// list all address added by the user
	allUserAddress, err := repository.GetAllAddresses(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	// get available payment options
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
