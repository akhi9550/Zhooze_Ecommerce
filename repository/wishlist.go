package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"errors"
)

func GetWishList(userID int) ([]models.WishListResponse, error) {
	var wishList []models.WishListResponse
	err := db.DB.Raw("SELECT products.id as product_id,products.name as product_name,products.description FROM products INNER JOIN wish_lists ON products.id = wish_lists.product_id WHERE wish_lists.user_id = ?", userID).Scan(&wishList).Error
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishList, nil
}

func ProductExistInWishList(productID, userID int) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM wish_lists WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count).Error
	if err != nil {
		return false, errors.New("error checking user product already present")
	}
	return count > 0, nil
}

func AddToWishlist(userID, productID int) error {
	err := db.DB.Exec("INSERT INTO wish_lists (user_id,product_id) VALUES (?,?)", userID, productID).Error
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromWishList(userID, productID int) error {
	err := db.DB.Exec("DELETE FROM wish_lists WHERE user_id = ? AND product_id = ?", userID, productID).Error
	if err != nil {
		return err
	}
	return nil
}
