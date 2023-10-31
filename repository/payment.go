package repository

import (
	"Zhooze/db"
)

func CheckPaymentStatus(razorID string, orderID string) (string, error) {
	var paymentStatus string
	err := db.DB.Raw("SELECT payment_status FROM orders WHERE order_id = ?", orderID).Scan(&paymentStatus).Error
	if err != nil {
		return "", err
	}
	return paymentStatus, nil
}
func UpdatePaymentDetails(orderID string, paymentID string) error {
	err := db.DB.Exec("UPDATE razer_pays set payment_id = ? WHERE razor_id= ?", paymentID, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func AddRazorPayDetails(orderID string, razorPayOrderID string) error {
	err := db.DB.Exec("INSERT INTO razer_pays (order_id,razor_id) VALUES (?,?)", orderID, razorPayOrderID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateShipmentAndPaymentByOrderID(shipmentStatus string, paymentStatus string, orderID string) error {
	err := db.DB.Exec("UPDATE orders SET payment_status = ?,shipment_status = ?, WHERE order_id = ?", paymentStatus, shipmentStatus, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
