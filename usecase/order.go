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
		return errors.New("the order is not dome by this user")
	}
	orderProductDetails, err := repository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderId)
	if err != nil {
		return err()
	}
	if shipmentStatus == "deliverd" {
		return errors.New("item already delivered, cannot cancel")
	}
	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("item already delivered, cannot cancel")
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
