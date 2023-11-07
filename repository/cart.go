package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
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

	if err := db.DB.Raw(`SELECT carts.id as cart_id,carts.user_id,users.firstname as user_name,carts.product_id,products.name as product_name,carts.quantity,carts.total_price FROM carts
	 inner join users 
	 on carts.user_id = users.id 
	 inner join products 
	 on carts.product_id = products.id 
	 where user_id = ?`, userID).Scan(&cartResponse).Error; err != nil {
		return []models.Cart{}, err
	}
	return cartResponse, nil

}
func GetTotalPrice(userID int) (models.CartTotal, error) {

	var cartTotal models.CartTotal
	err := db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.TotalPrice).Error
	if err != nil {
		return models.CartTotal{}, err
	}
	err = db.DB.Raw("SELECT COALESCE(SUM(total_price), 0) FROM carts WHERE user_id = ?", userID).Scan(&cartTotal.FinalPrice).Error
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
	if err := db.DB.Raw("SELECT COUNT(*) FROM carts WHERE user_id = ? ", userID).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil

}
func EmptyCart(userID int) error {
	if err := db.DB.Exec("DELETE FROM orders WHERE cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID).Error; err != nil {
		return err
	}

	// Then delete the cart
	if err := db.DB.Exec("DELETE FROM carts WHERE user_id = ?", userID).Error; err != nil {
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
	err := db.DB.Raw("SELECT quantity FROM carts WHERE user_id = ? AND product_id = ?", userId, productId).Scan(&productQty).Error
	if err != nil {
		return 0, err
	}
	return productQty, nil
}
func AddItemIntoCart(userId int, productId int, Quantity int, productprice float64) error {
	if err := db.DB.Exec("INSERT INTO carts (user_id,product_id,quantity,total_price) VALUES (?,?,?,?)", userId, productId, Quantity, productprice).Error; err != nil {
		return err
	}
	return nil

}
func TotalPriceForProductInCart(userID int, productID int) (float64, error) {

	var totalPrice float64
	if err := db.DB.Raw("SELECT SUM(total_price) as total_price FROM carts WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&totalPrice).Error; err != nil {
		return 0.0, err
	}
	return totalPrice, nil
}
func UpdateCart(quantity int, price float64, userID int, product_id int) error {

	if err := db.DB.Exec("UPDATE carts SET quantity = ?, total_price = ? WHERE user_id = ? AND product_id = ?", quantity, price, userID, product_id).Error; err != nil {
		return err
	}

	return nil

}
func ProductExist(userID int, productID int) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM carts WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error; err != nil {
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
	if err := db.DB.Raw("SELECT quantity,total_price FROM carts WHERE user_id = ? AND product_id = ?", userId, productId).Scan(&cartDetails).Error; err != nil {
		return struct {
			Quantity   int
			TotalPrice float64
		}{}, err
	}
	return cartDetails, nil
}
func RemoveProductFromCart(userID int, product_id int) error {

	if err := db.DB.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ?", uint(userID), uint(product_id)).Error; err != nil {
		return err
	}

	return nil
}
func UpdateCartDetails(cartDetails struct {
	Quantity   int
	TotalPrice float64
}, userId int, productId int) error {
	if err := db.DB.Raw("UPDATE carts SET quantity = ? , total_price = ? WHERE user_id = ? AND product_id = ? ", cartDetails.Quantity, cartDetails.TotalPrice, userId, productId).Scan(&cartDetails).Error; err != nil {
		return err
	}
	return nil

}
func CartAfterRemovalOfProduct(user_id int) ([]models.Cart, error) {
	var cart []models.Cart
	if err := db.DB.Raw("SELECT carts.product_id,products.name as product_name,carts.quantity,carts.total_price FROM carts INNER JOIN products on carts.product_id = products.id WHERE carts.user_id = ?", user_id).Scan(&cart).Error; err != nil {
		return []models.Cart{}, err
	}
	return cart, nil
}
