package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
)

func AddToCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, _, err := repository.CheckProduct(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}

	if !ok {
		return models.CartResponse{}, errors.New("product Does not exist")
	}

	QuantityOfProductInCart, err := repository.QuantityOfProductInCart(user_id, product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	quantityOfProduct, err := repository.GetQuantityFromProductID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	if quantityOfProduct == 0 {
		return models.CartResponse{}, errors.New("out of stock")
	}
	if quantityOfProduct == QuantityOfProductInCart {
		return models.CartResponse{}, errors.New("stock limit exceeded")
	}
	cart_id,err:=repository.CreateCart(user_id)
	if err!=nil{
		return models.CartResponse{}, err
	}
	productPrice, err := repository.GetPriceOfProductFromID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	if QuantityOfProductInCart == 0 {
		err := repository.AddItemIntoCart(cart_id, product_id, 1, productPrice)
		if err != nil {

			return models.CartResponse{}, err
		}

	} else {
		currentTotal, err := repository.TotalPriceForProductInCart(user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		err = repository.UpdateCart(QuantityOfProductInCart+1, currentTotal+productPrice, user_id, product_id)
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

// 	// if cartId <= 0 {
// 	// 	newCart, err := repository.MakeNewCart(user_id)
// 	// 	if err != nil {
// 	// 		return models.CartResponse{}, err
// 	// 	}
// 	// 	cartId = newCart.Id
// 	// 	fmt.Println("cart id", cartId, newCart)
// 	// }

//		ok, _, err := repository.CheckProduct(product_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		if !ok {
//			return models.CartResponse{}, errors.New("product doesnot exist")
//		}
//		QuantityofProductInCart, err := repository.QuantityOfProductInCart(user_id, product_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		QuantityofProduct, err := repository.GetQuantityFromProductID(product_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		if QuantityofProduct == 0 {
//			return models.CartResponse{}, errors.New("out of stock")
//		}
//		if QuantityofProduct == QuantityofProductInCart {
//			return models.CartResponse{}, errors.New("stock limit exceeded")
//		}
//		cartId,err:=repository
//		productPrice, err := repository.GetPriceOfProductFromID(product_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		fmt.Println("👺prodcut price", productPrice)
//		if QuantityofProductInCart == 0 {
//			err := repository.AddToCart(user_id, product_id, 1, cartId, productPrice)
//			if err != nil {
//				return models.CartResponse{}, err
//			}
//		} else {
//			currentTotal, err := repository.TotalPriceForProductInCart(user_id, product_id)
//			if err != nil {
//				return models.CartResponse{}, err
//			}
//			err = repository.UpdateCart(QuantityofProductInCart+1, currentTotal+productPrice, user_id, product_id)
//			if err != nil {
//				return models.CartResponse{}, err
//			}
//		}
//		cartDetails, err := repository.DisplayCart(user_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		cartTotal, err := repository.GetTotalPrice(user_id)
//		if err != nil {
//			return models.CartResponse{}, err
//		}
//		return models.CartResponse{
//			UserName:   cartTotal.UserName,
//			TotalPrice: cartTotal.TotalPrice,
//			Cart:       cartDetails,
//		}, nil
//	}
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
	// cartDetails.Quantity = cartDetails.Quantity - 1
	// if cartDetails.Quantity == 0 {
	if err := repository.RemoveProductFromCart(user_id, product_id); err != nil {
		return models.CartResponse{}, err

	}
	fmt.Println("👺👺👺", cartDetails.TotalPrice)
	// if cartDetails.Quantity != 0 {
	// 	product_price, err := repository.GetPriceOfProductFromID(product_id)
	// 	if err != nil {
	// 		return models.CartResponse{}, err
	// 	}
	// 	fmt.Println("😎😎😎", product_price)
	// 	cartDetails.TotalPrice = cartDetails.TotalPrice - product_price
	// 	err = repository.UpdateCartDetails(cartDetails, user_id, product_id)
	// 	if err != nil {
	// 		return models.CartResponse{}, err
	// 	}
	// }
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
