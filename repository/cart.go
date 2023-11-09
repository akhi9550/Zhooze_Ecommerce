package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"fmt"
)

func DisplayCart(userID int) ([]models.Cart, error) {

	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return []models.Cart{}, err
	}

	if count == 0 {
		return []models.Cart{}, nil
	}

	var cartResponse []models.Cart

	if err := db.DB.Raw(`SELECT * FROM cart_items ci JOIN carts c ON ci.cart_id =  c.id WHERE c.user_id = $1 `, userID).Scan(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}
func GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE user_id = ?", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE user_id = ?", userID).Scan(&cartTotal.FinalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = db.DB.Raw("SELECT firstname as user_name FROM users WHERE id = ?", userID).Scan(&cartTotal.UserName).Error
	if err != nil {
		return models.CartTotal{}, err
	}

	return cartTotal, nil

}
func CartExist(userID int) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func EmptyCart(userID int) error {
	if err := db.DB.Exec("DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID).Error; err != nil {
		return err
	}
	return nil
}
func CheckProduct(product_id int) (bool, string, error) {
	var count int

	err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id = ?", product_id).Scan(&count).Error

	if err != nil {
		return false, "", err
	}

	if count > 0 {
		var category string
		err := db.DB.Raw("SELECT categories.category FROM categories INNER JOIN products ON products.category_id = categories.id WHERE products.id = ?", product_id).Scan(&category).Error

		if err != nil {
			return false, "", err
		}
		return true, category, nil
	}
	return false, "", nil
}
func QuantityOfProductInCart(userId int, productId int) (int, error) {
	var productQty int
	err := db.DB.Raw("SELECT cart_items.quantity FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE carts.user_id = ? AND cart_items.product_id = ?", userId, productId).Scan(&productQty).Error
	if err != nil {
		return 0, err
	}
	return productQty, nil
}
func AddItemIntoCart(cartId int, productId int, Quantity int, productprice float64) error {
	if err := db.DB.Exec("INSERT INTO cart_items(cart_id,product_id,quantity,total_price) VALUES ($1 ,$2 ,$3 ,$4) ", cartId, productId, Quantity, productprice).Error; err != nil {
		return err
	}
	return nil

}
func CreateCart(userid int) (int, error) {
	var a int
	if err := db.DB.Exec("INSERT INTO carts (user_id) values(?)", userid).Error; err != nil {
		return 0, err
	}
	if err := db.DB.Raw("SELECT id FROM carts WHERE user_id = ?", userid).Scan(&a).Error; err != nil {
		return 0, err
	}
	fmt.Println("🤷‍♂️cart id:", a)
	return a, nil

}

func AddToCart(user_id, product_id, quantity, cartId int, productPrice float64) error {
	query := "INSERT INTO cart_items(cart_id,product_id,quantity,total_price) VALUES ($1 ,$2 ,$3 ,$4) "

	if err := db.DB.Exec(query, cartId, product_id, quantity, productPrice).Error; err != nil {
		return err
	}
	return nil
}

func TotalPriceForProductInCart(userID int, productID int) (float64, error) {

	var totalPrice float64
	if err := db.DB.Raw("SELECT SUM(total_price) as total_price FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE carts.user_id = ? AND cart_items.product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}
func UpdateCart(quantity int, price float64, userID int, product_id int) error {
	fmt.Println("lllllllllllllll", quantity, price, userID, product_id)
	if err := db.DB.Exec(`UPDATE cart_items
	SET quantity = ?, total_price = ? 
	WHERE product_id = ? 
	AND cart_id 
	IN (SELECT id FROM carts WHERE user_id = ?)`, quantity, price, product_id, userID).Error; err != nil {
		return err
	}

	return nil

}
func ProductExist(userID int, productID int) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE carts.user_id = ? AND cart_items.product_id = ?", userID, productID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func GetQuantityAndProductDetails(userId int, productId int, cartDetails struct {
	Quantity   int
	TotalPrice float64
}) (struct {
	Quantity   int
	TotalPrice float64
}, error) {
	if err := db.DB.Raw("SELECT cart_items.quantity,cart_items.total_price FROM cart_items JOIN carts ON cart_items.cart_id = carts.id WHERE carts.user_id = ? AND cart_items.product_id = ?", userId, productId).Scan(&cartDetails).Error; err != nil {
		return struct {
			Quantity   int
			TotalPrice float64
		}{}, err
	}
	return cartDetails, nil
}
func RemoveProductFromCart(userID int, product_id int) error {

	if err := db.DB.Exec("DELETE FROM cart_items WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?) AND product_id = ?", uint(userID), uint(product_id)).Error; err != nil {
		return err
	}

	return nil
}
func UpdateCartDetails(cartDetails struct {
	Quantity   int
	TotalPrice float64
}, userId int, productId int) error {
	if err := db.DB.Raw("UPDATE cart_items AS ci "+"SET ci.quantity = ?, ci.total_price = ? "+"WHERE ci.cart_id IN (SELECT id FROM carts WHERE user_id = ?) "+"AND ci.product_id = ?", cartDetails.Quantity, cartDetails.TotalPrice, userId, productId).Scan(&cartDetails).Error; err != nil {
		return err
	}
	return nil

}
func CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error) {
	var cart []models.Cart
	if err := db.DB.Raw("SELECT ci.cart_id, ci.product_id, p.name as product_name, ci.quantity, ci.total_price "+"FROM cart_items as ci "+"INNER JOIN products as p ON ci.product_id = p.id "+"INNER JOIN carts as c ON ci.cart_id = c.id "+"WHERE c.user_id = ?", user_id).Scan(&cart).Error; err != nil {
		return []models.Cart{}, err
	}
	return cart, nil
}
func MakeNewCart(userId int) (models.Carts, error) {
	var cart models.Carts
	cart.UserId = userId
	if err := db.DB.Create(&cart).Error; err != nil {
		fmt.Println("error while inserting to carts", err)
	}
	if err := db.DB.Last(&cart).Scan(&cart).Error; err != nil {
		return models.Carts{}, err
	}
	fmt.Println("cart:", cart)
	return cart, nil

}
