package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
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
	err := db.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = true WHERE order_id = ?", order_id).Error
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
func UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {
	for _, ok := range orderProducts {
		var quantity int
		if err := db.DB.Raw("SELECT quantity FROM products WHERE id = ?", ok.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}
		ok.Quantity += quantity
		if err := db.DB.Exec("UPDATE products SET quantity  = ? WHERE id = ?", ok.Quantity, ok.ProductId).Error; err != nil {
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
	err := db.DB.Raw("select orders.order_id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin from orders inner join users on orders.user_id = users.id inner join addresses on users.id = addresses.user_id limit ? offset ?", 2, offset).Scan(&orderDatails).Error
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
