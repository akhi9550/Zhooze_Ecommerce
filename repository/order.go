package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"fmt"

	"Zhooze/utils/models"
	"errors"
)

func CheckOrderID(orderId string) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetShipmentStatus(orderId string) (string, error) {
	var status string
	err := db.DB.Raw("SELECT shipment_status FROM orders WHERE id= ?", orderId).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}

func ApproveOrder(order_id string) error {
	err := db.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = 'true' WHERE id = ?", order_id).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelOrders(order_id string) error {
	status := "cancelled"
	err := db.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false' WHERE id = ? ", status, order_id).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = db.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ? ", order_id).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = db.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE id = ?", order_id).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func GetProductDetailsFromOrders(order_id string) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := db.DB.Raw("SELECT product_id,quantity FROM carts WHERE id = ?", order_id).Scan(&OrderProductDetails).Error; err != nil {
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
	err := db.DB.Raw("SELECT orders.id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin FROM orders INNER JOIN users ON orders.user_id = users.id INNER JOIN addresses ON orders.address_id = addresses.id limit ? offset ?", 2, offset).Scan(&orderDatails).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, nil
	}
	return orderDatails, nil

}

// func GetPaymentStatus(orderID string) (string, error) {
// 	var paymentStatus string
// 	err := db.DB.Raw("SELECT payment_status FROM orders WHERE id = ?", orderID).Scan(&paymentStatus).Error
// 	if err != nil {
// 		return "", err
// 	}
// 	return paymentStatus, nil
// }
// func RefundOrder(paymentStatus string, orderID string) error {
// 	err := db.DB.Exec("UPDATE orders SET payment_status = ?, shipment_status = 'returned' WHERE id = ?", paymentStatus, orderID).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func OrderItemsFromCart(cartID int) (domain.Order, error) {
	var cart domain.Cart
	if err := db.DB.Raw("SELECT * FROM carts WHERE id = ?", cartID).Scan(&cart).Error; err != nil {
		return domain.Order{}, err
	}
	fmt.Println(cart)
	var orderItems domain.Order
	err := db.DB.Raw(`
    INSERT INTO orders  (created_at, user_id, address_id, cart_id)
    SELECT NOW(), c.user_id, a.id , c.id
    FROM carts c
    JOIN addresses a ON c.user_id = a.user_id
    WHERE c.id = ?
    RETURNING orders.id,orders.user_id, orders.address_id, orders.cart_id`, cartID).Scan(&orderItems).Error
	if err != nil {
		return domain.Order{}, err
	}
	return orderItems, nil

}
func AddpaymentMethod(paymentID int, orderID uint) error {
	fmt.Println("payment id : ", orderID)
	err := db.DB.Exec(`UPDATE orders SET payment_method_id = $1 WHERE id = $2`, paymentID, orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckAddressAvailabilityWithID(addressID, userID int) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?", addressID, userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func CheckCartAvailabilityWithID(cartID, UserID int) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE id = ? AND user_id = ?", cartID, UserID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func FindOrderStock(cartID int) (int, error) {
	var count int
	if err := db.DB.Raw("SELECT quantity FROM carts WHERE id = ?", cartID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func AddAmountToOrder(Price float64, orderID uint) error {
	fmt.Println("🤷‍♂️🤷‍♂️🤷‍♂️", Price, orderID)
	err := db.DB.Exec("UPDATE orders SET final_price = ? WHERE id = ?", Price, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func GetOrder(orderID int) (domain.Order, error) {
	var order domain.Order
	err := db.DB.Raw("SELECT * FROM orders WHERE id = ?", orderID).Scan(&order).Error
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func FindProductFromCart(cartID int) (int, error) {
	fmt.Println("carttttt", cartID)
	var p int
	if err := db.DB.Raw("SELECT product_id FROM carts WHERE id = ?", cartID).Scan(&p).Error; err != nil {
		return 0, err
	}
	return p, nil
}
func CartEmpty(cartID int) error {
	if err := db.DB.Exec("DELETE FROM carts WHERE id = ?", cartID).Error; err != nil {
		return err
	}
	return nil

}
func ProductStockMinus(productID, stock int) error {
	err := db.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", stock, productID).Error
	if err != nil {
		return err
	}
	return nil
}
func GetPaymentId(paymentID int) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE id = ? ", paymentID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func TotalAmountInCart(cartID int) (float64, error) {
	fmt.Println("😊😊", cartID)
	var price float64
	if err := db.DB.Raw("SELECT total_price FROM carts WHERE id = $1", cartID).Scan(&price).Error; err != nil {
		return 0, err
	}
	fmt.Println("👺", price)
	return price, nil

}

func GetOrderDetails(userID int, page int, count int) ([]models.CombinedOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orders []models.CombinedOrderDetails
	err := db.DB.Raw(`SELECT  o.id as order_id,o.final_price ,o.shipment_status, o.payment_status,u.firstname as firstname,u.email as email,u.phone as phone,a.house_name as house_name,a.street as street,a.city as city,a.state as state,a.pin as pin FROM orders o JOIN users u on o.user_id = u.id JOIN addresses a on o.address_id = a.id WHERE o.user_id = ? LIMIT ? OFFSET ?`, userID, count, offset).Scan(&orders).Error
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orders, nil
}

func UserOrderRelationship(orderID string, userID int) (int, error) {
	var testUserID int
	err := db.DB.Raw("SELECT user_id FROM orders WHERE id = ?", orderID).Scan(&testUserID).Error
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
	if err := db.DB.Raw("SELECT address_id FROM orders WHERE id =?", orderId).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := db.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second  in address")
	}
	return addressInfoResponse, nil
}
func GetOrderDetailOfAproduct(orderId string) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := db.DB.Raw("SELECT id,final_price,shipment_status,payment_status FROM orders WHERE id = ?", orderId).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}
func GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {
	fmt.Println("ggggg", orderID)
	var orderDetails models.CombinedOrderDetails
	err := db.DB.Raw(`
    SELECT orders.id as order_id, orders.final_price, orders.shipment_status, orders.payment_status,
    users.firstname, users.email, users.phone, addresses.house_name, addresses.state,
    addresses.pin, addresses.street, addresses.city
    FROM orders
    INNER JOIN users ON orders.user_id = users.id
    INNER JOIN addresses ON users.id = addresses.user_id
    WHERE orders.id = $1`, orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, err
	}
	fmt.Println("fffffffff", orderDetails)
	return orderDetails, nil
}
