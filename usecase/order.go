package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
)

func GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error) {
	OrderDetails, err := repository.GetOrderDetails(userID, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return OrderDetails, nil
}
func CancelOrders(orderID string, userID int) error {
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
	err = repository.UpdateQuantityOfProduct(orderProductDetails)
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

	// get all items from users cart
	cartItems, err := repository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	// get grand total of all the product
	grandTotal, err := repository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,

		Grand_Total: grandTotal.TotalPrice,
		Total_Price: grandTotal.FinalPrice,
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
