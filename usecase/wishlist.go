package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func GetWishList(userID int) ([]models.WishListResponse, error) {
	wishList, err := repository.GetWishList(userID)
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishList, err
}

func AddToWishlist(userID, productID int) error {
	productExist, err := repository.DoesProductExist(productID)
	if err != nil {
		return err
	}
	if !productExist {
		return errors.New("product does not exist")
	}
	productExistInWishList, err := repository.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if productExistInWishList {
		return errors.New("product already exist in wishlist")
	}
	err = repository.AddToWishlist(userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromWishlist(productID, userID int) error {
	productExistInWishList, err := repository.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if !productExistInWishList {
		return errors.New("product does not exist in wishlist")
	}
	err = repository.RemoveFromWishList(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
