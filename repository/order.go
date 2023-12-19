package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/helper"
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
	if err := db.DB.Raw("SELECT COUNT(*) FROM addresses WHERE user_id = ? AND id = ?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
func PaymentExist(orderBody models.OrderIncoming) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM payment_methods WHERE id = ?", orderBody.PaymentID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func CheckOrderID(orderId int) (bool, error) {
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
	db.DB.Raw("SELECT id as order_id,final_price,shipment_status,payment_status FROM orders WHERE user_id = ? LIMIT ? OFFSET ? ", userId, count, offset).Scan(&orderDetails)
	var fullOrderDetails []models.FullOrderDetails
	for _, od := range orderDetails {
		var orderProductDetails []models.OrderProductDetails
		db.DB.Raw(`SELECT
		order_items.product_id,
		products.name AS product_name,
		order_items.quantity,
		order_items.total_price
	    FROM
		order_items
	    INNER JOIN
		products ON order_items.product_id = products.id
	    WHERE
		order_items.order_id = $1 `, od.OrderId).Scan(&orderProductDetails)
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

func GetShipmentStatus(orderID int) (string, error) {
	var status string
	err := db.DB.Raw("SELECT shipment_status FROM orders WHERE id= ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
func PaymentStatus(orderID int) (string, error) {
	var status string
	err := db.DB.Raw("SELECT payment_status FROM orders WHERE id= ?", orderID).Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
func UserOrderRelationship(orderID int, userID int) (int, error) {

	var testUserID int
	err := db.DB.Raw("select user_id from orders where id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}

func ApproveOrder(orderID int) error {
	err := db.DB.Exec("UPDATE orders SET shipment_status = 'order placed' , approval = 'true' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func CancelOrders(orderID int) error {
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

func PaymentMethodID(orderID int) (int, error) {
	var a int
	err := db.DB.Raw("SELECT payment_method_id FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}
func GetProductDetailsFromOrders(orderID int) ([]models.OrderProducts, error) {
	var OrderProductDetails []models.OrderProducts
	if err := db.DB.Raw("SELECT product_id,quantity as stock FROM order_items WHERE order_id = ?", orderID).Scan(&OrderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}
	return OrderProductDetails, nil
}

func UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {

	for _, od := range orderProducts {

		var quantity int
		if err := db.DB.Raw("SELECT stock FROM products WHERE id = ?", od.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}

		od.Stock += quantity
		if err := db.DB.Exec("UPDATE products SET stock = ? WHERE id = ?", od.Stock, od.ProductId).Error; err != nil {
			return err
		}
	}
	return nil

}
func TotalAmountFromOrder(orderID int) (float64, error) {
	var total float64
	err := db.DB.Raw("SELECT final_price FROM orders WHERE id = ?", orderID).Scan(&total).Error
	if err != nil {
		return 0.0, err
	}
	return total, nil
}
func UserIDFromOrder(orderID int) (int, error) {
	var a int
	err := db.DB.Raw("SELECT user_id FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return 0, err
	}
	return a, nil
}
func UpdateAmountToWallet(userID int, amount float64) error {
	err := db.DB.Exec("UPDATE wallets SET amount = amount + ? WHERE user_id = ?", amount, userID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateHistory(userID, orderID int, amount float64, reason string) error {
	err := db.DB.Exec("INSERT INTO wallet_histories (user_id ,order_id ,description ,amount) VALUES (?,?,?,?)", userID, orderID, reason, amount).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateHistoryForDebit(userID, orderID int, amount float64, reason string) error {
	err := db.DB.Exec("INSERT INTO wallet_histories (user_id ,order_id ,description ,amount) VALUES (?,?,?,?)", userID, orderID, reason, amount).Error
	if err != nil {
		return err
	}
	err = db.DB.Exec("UPDATE wallet_histories SET is_credited = 'false' where user_id = ? AND order_id = ?", userID, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func PaymentAlreadyPaid(orderID int) (bool, error) {
	var a bool
	err := db.DB.Raw("SELECT shipment_status = 'processing' AND payment_status = 'paid' OR shipment_status = 'order placed' FROM orders WHERE id = ?", orderID).Scan(&a).Error
	if err != nil {
		return false, err
	}
	return a, nil
}

func GetOrderDetailsByOrderId(orderID int) (models.CombinedOrderDetails, error) {

	var orderDetails models.CombinedOrderDetails
	err := db.DB.Raw(`SELECT
    orders.id as order_id,
    orders.final_price,
    orders.shipment_status,
    orders.payment_status,
    users.firstname,
    users.email,
    users.phone,
    addresses.house_name,
    addresses.state,
    addresses.street,
    addresses.city,
    addresses.pin
FROM
    orders
INNER JOIN
    users ON orders.user_id = users.id
INNER JOIN
    addresses ON users.id = addresses.user_id
WHERE
    orders.id = ?`, orderID).Scan(&orderDetails).Error
	if err != nil {
		return models.CombinedOrderDetails{}, nil
	}

	return orderDetails, nil
}

func OrderItems(ob models.OrderIncoming, price float64) (int, error) {
	var id int
	query := `
    INSERT INTO orders (created_at , user_id , address_id , payment_method_id , final_price)
    VALUES (NOW(),?, ?, ?, ?)
    RETURNING id`
	db.DB.Raw(query, ob.UserID, ob.AddressID, ob.PaymentID, price).Scan(&id)
	return id, nil
}
func UpdateWallectAfterOrder(userID int, amount float64) error {
	err := db.DB.Exec("UPDATE wallets SET amount = amount - ? WHERE user_id = ?", amount, userID).Error
	if err != nil {
		return err
	}
	return nil
}
func OrderExist(orderID int) error {
	err := db.DB.Raw("SELECT id FROM orders WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateOrder(orderID int) error {
	err := db.DB.Exec("UPDATE orders SET Shipment_status = 'processing' WHERE id = ?", orderID).Error
	if err != nil {
		return err
	}
	return nil
}

func AddOrderProducts(order_id int, cart []models.Cart) error {
	query := `
    INSERT INTO order_items (order_id,product_id,quantity,total_price)
    VALUES (?, ?, ?, ?) `
	for _, v := range cart {
		var productID int
		if err := db.DB.Raw("SELECT id FROM products WHERE name = $1", v.ProductName).Scan(&productID).Error; err != nil {
			return err
		}
		if err := db.DB.Exec(query, order_id, productID, v.Quantity, v.TotalPrice).Error; err != nil {
			return err
		}
	}
	return nil
}
func UpdateCartAfterOrder(userID, productID int, quantity float64) error {
	err := db.DB.Exec("DELETE FROM carts WHERE user_id = ? and product_id = ?", userID, productID).Error
	if err != nil {
		return err
	}

	err = db.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", quantity, productID).Error
	if err != nil {
		return err
	}

	return nil
}
func GetBriefOrderDetails(orderID, paymentID int) (domain.OrderSuccessResponse, error) {
	if paymentID == 3 {
		err := db.DB.Exec("UPDATE orders SET shipment_status ='processing' , payment_status ='paid' WHERE id = ?", orderID).Error
		if err != nil {
			return domain.OrderSuccessResponse{}, err

		}
	}
	var orderSuccessResponse domain.OrderSuccessResponse
	err := db.DB.Raw(`SELECT id as order_id,shipment_status FROM orders WHERE id = ?`, orderID).Scan(&orderSuccessResponse).Error
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

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
func GetAllOrderDetailsBrief(page, count int) ([]models.CombinedOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDatails []models.CombinedOrderDetails
	err := db.DB.Raw("SELECT orders.id as order_id,orders.final_price,orders.shipment_status,orders.payment_status,users.firstname,users.email,users.phone,addresses.house_name,addresses.street,addresses.city,addresses.state,addresses.pin FROM orders INNER JOIN users ON orders.user_id = users.id INNER JOIN addresses ON orders.address_id = addresses.id limit ? offset ?", count, offset).Scan(&orderDatails).Error
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
	err := db.DB.Exec("UPDATE orders SET final_price = ? WHERE id = ?", Price, orderID).Error
	if err != nil {
		return err
	}
	return nil
}
func GetOrder(orderID int) (domain.Orders, error) {
	var order domain.Orders
	err := db.DB.Raw("SELECT * FROM orders WHERE id = ?", orderID).Scan(&order).Error
	if err != nil {
		return domain.Orders{}, err
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
func TotalAmountInCart(userID int) (float64, error) {
	var price float64
	if err := db.DB.Raw("SELECT SUM(total_price) FROM carts WHERE  user_id= $1", userID).Scan(&price).Error; err != nil {
		return 0, err
	}
	return price, nil
}
func GetCouponDiscountPrice(UserID int, Total float64) (float64, error) {
	discountPrice, err := helper.GetCouponDiscountPrice(UserID, Total, db.DB)
	if err != nil {
		return 0.0, err
	}

	return discountPrice, nil

}
func GetReferralDiscountPrice(FinalPrice float64, UserID int) (float64, error) {
	discountPrice, err := helper.GetReferralDiscountPrice(FinalPrice, UserID, db.DB)
	if err != nil {
		return 0.0, err
	}

	return discountPrice, nil

}
func UpdateRefferal(TotalPrice float64, userID int) error {
	if TotalPrice != 0.0 {
		err := db.DB.Exec("UPDATE referrals SET referral_amount = 0 WHERE user_id = ?", userID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func UpdateCouponDetails(discount_price float64, UserID int) error {

	if discount_price != 0.0 {
		err := db.DB.Exec("update used_coupons set used = true where user_id = ?", UserID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// ///////////////
func WallectAmount(userID int) (float64, error) {
	var a float64
	err := db.DB.Raw("SELECT amount FROM wallets WHERE user_id = $1", userID).Scan(&a).Error
	if err != nil {
		return 0.0, err
	}
	return a, nil
}

// ////////////////////////
func GetAllAddresses(userID int) ([]models.AddressInfoResponse, error) {
	var addressResponse []models.AddressInfoResponse
	err := db.DB.Raw(`SELECT * FROM addresses WHERE user_id = $1`, userID).Scan(&addressResponse).Error
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	return addressResponse, nil
}
func GetAllPaymentOption(userID int) ([]models.PaymentDetails, error) {
	var fullPaymentDetails []models.PaymentDetails
	var paymentMethods []models.PaymentDetail
	err := db.DB.Raw("SELECT * FROM payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}
	var a float64
	err = db.DB.Raw("SELECT amount FROM wallets WHERE user_id = ?", userID).Scan(&a).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}
	fullPaymentDetails = append(fullPaymentDetails, models.PaymentDetails{PaymentDetail: paymentMethods, WallectAmount: a})

	return fullPaymentDetails, nil

}
func GetAddressFromOrderId(orderID int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	var addressId int
	if err := db.DB.Raw("SELECT address_id FROM orders WHERE id =?", orderID).Scan(&addressId).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("first in orders")
	}
	if err := db.DB.Raw("SELECT * FROM addresses WHERE id=?", addressId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, errors.New("second  in address")
	}
	return addressInfoResponse, nil
}
func GetOrderDetailOfAproduct(orderID int) (models.OrderDetails, error) {
	var OrderDetails models.OrderDetails

	if err := db.DB.Raw("SELECT id,final_price,shipment_status,payment_status FROM orders WHERE id = ?", orderID).Scan(&OrderDetails).Error; err != nil {
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
func GetDetailedOrderThroughId(orderId int) (models.CombinedOrderDetails, error) {
	var body models.CombinedOrderDetails
	query := `
	SELECT 
        o.id AS order_id,
        o.final_price AS final_price,
        o.shipment_status AS shipment_status,
        o.payment_status AS payment_status,
        u.firstname AS firstname,
        u.email AS email,
        u.phone AS phone,
        a.house_name AS house_name,
        a.street AS street,
        a.city AS city,
		a.state AS state,
        a.pin AS pin
	FROM orders o
	JOIN users u ON o.user_id = u.id
	JOIN addresses a ON o.address_id = a.id 
	WHERE o.id = ?
	`
	if err := db.DB.Raw(query, orderId).Scan(&body).Error; err != nil {
		err = errors.New("error in getting detailed order through id in repository: " + err.Error())
		return models.CombinedOrderDetails{}, err
	}
	fmt.Println("body in repo", body.OrderId)
	return body, nil
}

func GetItemsByOrderId(orderId int) ([]models.Invoice, error) {
	var items []models.Invoice

	query := `
	SELECT oi.id AS id,product_id, oi.quantity, oi.total_price, o.id AS order_id, o.created_at,
	FROM orders o
	JOIN order_items oi ON o.id = oi.order_id
	WHERE o.id = ?;
	`

	if err := db.DB.Raw(query, orederId).Scan(&items).Error; err != nil {
		return []models.Invoice{}, err
	}

	return items, nil
}
