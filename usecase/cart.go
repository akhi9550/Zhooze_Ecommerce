package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func AddToCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, _, err := repository.CheckProduct(product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("product doesnot exist")
	}
	QuantityofProductInCart, err := repository.QuantityOfProductInCart(user_id, product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	QuantityofProduct, err := repository.GetQuantityFromProductID(product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if QuantityofProduct == 0 {
		return models.CartResponse{}, errors.New("out of stock")
	}
	if QuantityofProduct == QuantityofProductInCart {
		return models.CartResponse{}, errors.New("stock limit exceeded")
	}
	productPrice, err := repository.GetPriceOfProductFromID(product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if QuantityofProductInCart == 0 {
		err := repository.AddItemIntoCart(user_id, product_id, 1, productPrice)
		if err != nil {
			return models.CartResponse{}, err
		}
	} else {
		currentTotal, err := repository.TotalPriceForProductInCart(user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		err = repository.UpdateCart(QuantityofProductInCart+1, currentTotal+productPrice, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
	}
	cartDetails, err := repository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := repository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil
}
func RemoveFromCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, err := repository.ProductExist(user_id, product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("product doesn't exist in the cart")
	}
	var cartDetails struct {
		Quantity   int
		TotalPrice float64
	}
	cartDetails, err = repository.GetQuantityAndProductDetails(user_id, product_id, cartDetails)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartDetails.Quantity = cartDetails.Quantity - 1
	if cartDetails.Quantity == 0 {
		if err := repository.RemoveProductFromCart(user_id, product_id); err != nil {
			return models.CartResponse{}, err
		}
	}
	if cartDetails.Quantity != 0 {
		product_price, err := repository.GetPriceOfProductFromID(product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		cartDetails.TotalPrice = cartDetails.TotalPrice - product_price
		err = repository.UpdateCartDetails(cartDetails, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
	}
	updatedCart, err := repository.CartAfterRemovalOfProduct(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := repository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       updatedCart,
	}, nil
}
func DisplayCart(user_id int) (models.CartResponse, error) {
	cart, err := repository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := repository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cart,
	}, nil
}
func EmptyCart(userID int) (models.CartResponse, error) {
	ok, err := repository.CartExist(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("cart already empty")
	}
	if err := repository.EmptyCart(userID); err != nil {
		return models.CartResponse{}, err
	}

	cartTotal, err := repository.GetTotalPrice(userID)

	if err != nil {
		return models.CartResponse{}, err
	}

	cartResponse := models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       []models.Cart{},
	}

	return cartResponse, nil

}
