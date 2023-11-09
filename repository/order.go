package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"fmt"

	"Zhooze/utils/models"
	"errors"
)

func DoesCartExist(userID int) (bool, error) {

	var exist bool
	err := db.DB.Raw("select exists(select 1 from carts where user_id = ?)", userID).Scan(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, nil
}
func AddressExist(orderBody models.OrderIncoming) (bool, error) {

	var count int
	if err := db.DB.Raw("select count(*) from addresses where user_id = ? and id = ?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
func CheckOrderID(orderId string) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM orders WHERE id = ?", orderId).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDetails []models.OrderDetails
	fmt.Println("userid is", userId, "page is ", page, "count is ", count, "offset is", offset)
	db.DB.Raw("SELECT id,final_price,shipment_status,payment_status FROM orders WHERE user_id = ? limit ? offset ? ", userId, count, offset).Scan(&orderDetails)
	fmt.Println("order details is ", orderDetails)

	var fullOrderDetails []models.FullOrderDetails
	// for each order select all the associated products and their details
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails
		db.DB.Raw("SELECT order_items.product_id,products.name as product_name,order_items.quantity,order_items.total_price from order_items inner join products on order_items.product_id = products.id where order_items.id = ?", od.OrderId).Scan(&orderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: od, OrderProductDetails: orderProductDetails})

	}

	return fullOrderDetails, nil

}

func GetOrderDetail(orderId int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := db.DB.Raw("select id,final_price,shipment_status,payment_status from orders where id = ?", orderId).Scan(&OrderDetails).Error; err != nil {
		return models.OrderDetails{}, err
	}
	return OrderDetails, nil
}

func GetShipmentStatus(orderId string) (string, error) {
	var status string
	err := db.DB.Raw("SELECT shipment_status FROM orders WHERE id= ?", orderId).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
func UserOrderRelationship(orderID string, userID int) (int, error) {

	var testUserID int
	err := db.DB.Raw("select user_id from orders where id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}

func ApproveOrder(orderID string) error {
	err := db.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = 'true' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelOrders(orderID string) error {
	status := "cancelled"
	err := db.DB.Exec("UPDATE orders SET shipment_status = ? , approval='false' WHERE id = ? ", status, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = db.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ? ", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = db.DB.Exec("UPDATE orders SET payment_status = 'refunded' WHERE id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := db.DB.Raw("SELECT product_id,quantity FROM order_items WHERE id = ?", orderID).Scan(&OrderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}
	return OrderProductDetails, nil
}

func UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {

	for _, od := range orderProducts {

		var quantity int
		if err := db.DB.Raw("select quantity from products where id = ?", od.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}

		od.Stock += quantity
		if err := db.DB.Exec("update products set quantity = ? where id = ?", od.Stock, od.ProductId).Error; err != nil {
			return err
		}
	}
	return nil

}

func GetOrderDetailsByOrderId(orderID string) (models.CombinedOrderDetails, error) {

	var orderDetails models.CombinedOrderDetails
	err := db.DB.Raw("SELECT orders.id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.state,addresses.pin,addresses.street,addresses.city from orders INNER JOIN users ON orders.user_id = users.id INNER JOIN addresses on users.id = addresses.user_id WHERE id = ?", orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, nil
	}

	return orderDetails, nil
}

func CreateOrder(orderDetails domain.Order) error {
	fmt.Println("mmmmmmmmm", orderDetails)
	err := db.DB.Create(&orderDetails).Error
	if err != nil {
		return err
	}
	return nil
}

func AddOrderItems(orderItemDetails domain.OrderItem, UserID int, ProductID uint, Quantity float64) error {
	fmt.Println("❌❌❌❌❌❌🥱", "details", orderItemDetails, "product", ProductID, "user", UserID, "quantity", Quantity)

	// after creating the order delete all cart items and also update the quantity of the product
	err := db.DB.Omit("id").Create(&orderItemDetails).Error
	if err != nil {
		return err
	}

	err = db.DB.Exec("DELETE FROM carts WHERE user_id = ? and product_id = ?", UserID, ProductID).Error
	if err != nil {
		return err
	}

	err = db.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", Quantity, ProductID).Error
	if err != nil {
		return err
	}

	return nil

}

func GetBriefOrderDetails(orderID int) (domain.OrderSuccessResponse, error) {
	fmt.Println("🤦‍♂️🤦‍♂️🤦‍♂️🤦‍♂️🤦‍♂️🤦‍♂️", orderID)
	var orderSuccessResponse domain.OrderSuccessResponse
	db.DB.Raw("SELECT id as order_id,shipment_status as order_status FROM orders WHERE id = ?", orderID).Scan(&orderSuccessResponse)
	return orderSuccessResponse, nil

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
	if err := db.DB.Raw("SELECT COUNT(*) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE cart_items.cart_id = ? AND carts.user_id = ?", cartID, UserID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func FindOrderStock(cartID int) (int, error) {
	var count int
	if err := db.DB.Raw("SELECT quantity FROM cart_items WHERE cart_id = ?", cartID).Scan(&count).Error; err != nil {
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
	var p int
	if err := db.DB.Raw("SELECT product_id FROM cart_items WHERE cart_id = ?", cartID).Scan(&p).Error; err != nil {
		return 0, err
	}
	return p, nil
}
func CartEmpty(cartID int) error {
	if err := db.DB.Exec("DELETE FROM cart_items WHERE cart_id = ?", cartID).Error; err != nil {
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
	var price float64
	if err := db.DB.Raw("SELECT sum(total_price) FROM cart_items WHERE cart_id = $1", cartID).Scan(&price).Error; err != nil {
		return 0, err
	}
	return price, nil

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

func GetProductsInCart(cart_id int) ([]int, error) {

	var cart_products []int

	if err := db.DB.Raw("select product_id from cart_items where cart_id=?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}

	return cart_products, nil

}
func FindProductNames(product_id int) (string, error) {

	var product_name string

	if err := db.DB.Raw("select name from products where id=?", product_id).Scan(&product_name).Error; err != nil {
		return "", err
	}

	return product_name, nil

}

func FindCartQuantity(cart_id, product_id int) (int, error) {

	var quantity int

	if err := db.DB.Raw("select quantity from cart_items where cart_id=$1 and product_id=$2", cart_id, product_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}

	return quantity, nil

}

func FindPrice(product_id int) (float64, error) {

	var price float64

	if err := db.DB.Raw("select price from products where id=?", product_id).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil

}
func FindStock(id int) (int, error) {
	var stock int
	err := db.DB.Raw("SELECT stock FROM prodcuts WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
