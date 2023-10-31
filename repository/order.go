package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"errors"
)

func CheckOrderID(orderId string) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM orders WHERE order_id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetShipmentStatus(orderId string) (string, error) {
	var status string
	err := db.DB.Raw("SELECT shipment_status FROM orders WHERE order_id= ?", orderId).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func ApproveOrder(order_id string) error {
	err := db.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = 'true' WHERE order_id = ?", order_id).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelOrders(order_id string) error {
	status := "cancelled"
	err := db.DB.Exec("UPDATE orders SET shipment_status = ? ,approval='false' WHERE order_id = ? ", status, order_id).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = db.DB.Raw("SELECT payment_method_id FROM orders WHERE order_id = ? ", order_id).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = db.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE order_id = ?", order_id).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func GetProductDetailsFromOrders(order_id string) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := db.DB.Raw("SELECT product_id,quantity FROM order_items WHERE order_id = ?", order_id).Scan(&OrderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}
	return OrderProductDetails, nil
}
func UpdateStockOfProduct(orderProducts []models.OrderProducts) error {
	for _, ok := range orderProducts {
		var quantity int
		if err := db.DB.Raw("SELECT stock FROM products WHERE id = ?", ok.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}
		ok.Stock += quantity
		if err := db.DB.Exec("UPDATE products SET stock  = ? WHERE id = ?", ok.Stock, ok.ProductId).Error; err != nil {
			return err
		}
	}
	return nil
}
func GetAllOrderDetailsBrief(page int) ([]models.CombinedOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * 2
	var orderDatails []models.CombinedOrderDetails
	err := db.DB.Raw("SELECT orders.order_id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin FROM orders inner join users on orders.user_id = users.id inner join addresses on users.id = addresses.user_id limit ? offset ?", 2, offset).Scan(&orderDatails).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, nil
	}
	return orderDatails, nil

}
func GetPaymentStatus(orderID string) (string, error) {
	var paymentStatus string
	err := db.DB.Raw("SELECT payment_status FROM orders WHERE order_id = ?", orderID).Scan(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}
func RefundOrder(paymentStatus string, orderID string) error {
	err := db.DB.Exec("UPDATE orders SET payment_status = ?, shipment_status = 'returned' WHERE order_id = ?", paymentStatus, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func GetOrderDetails(userID int, page int, count int) ([]models.FullOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDatails []models.OrderDetails
	db.DB.Raw("SELECT order_id, final_price, shipment_status, payment_status FROM orders WHERE user_id = ? LIMIT ? OFFSET ?", userID, count, offset).Scan(&orderDatails)
	var fullOrderDetails []models.FullOrderDetails
	for _, ok := range orderDatails {
		var OrderProductDetails []models.OrderProductDetails
		db.DB.Raw("SELECT o.product_id,products.name as product_name,o.quantity,o.total_price FROM order_items o inner join products on o.product_id = products.id where o.order_id = ?", ok.OrderId).Scan(&OrderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: ok, OrderProductDetails: OrderProductDetails})
	}
	return fullOrderDetails, nil
}
func UserOrderRelationship(orderID string, userID int) (int, error) {
	var testUserID int
	err := db.DB.Raw("SELECT user_id FROM orders WHERE order_id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}
func GetAllAddresses(userID int) ([]models.AddressInfoResponse, error) {
	var addressResponse []models.AddressInfoResponse
	err := db.DB.Raw(`SELECT * FROM addresses WHERE user_id = $1`, userID).Scan(&addressResponse).Error
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	return addressResponse, nil
}
func GetAllPaymentOption() ([]models.PaymentDetails, error) {
	var paymentMethods []models.PaymentDetails
	err := db.DB.Raw("SELECT * FROM payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}

	return paymentMethods, nil

}
func GetAddressFromOrderId(orderId string) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	var addressId int
	if err := db.DB.Raw("SELECT address_id FROM orders WHERE order_id =?", orderId).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := db.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second  in address")
	}
	return addressInfoResponse, nil
}
func GetOrderDetailOfAproduct(orderId string) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := db.DB.Raw("SELECT order_id,final_price,shipment_status,payment_status FROM orders WHERE order_id = ?", orderId).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}
func GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {
	var orderDetails models.CombinedOrderDetails
	err := db.DB.Raw("SELECT orders.order_id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.state,addresses.pin,addresses.street,addresses.city from orders inner join users on orders.user_id = users.id inner join addresses on users.id = addresses.user_id where order_id = ?", orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, err
	}

	return orderDetails, nil
}
