package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
)

func AddProductOffer(productOffer models.ProductOfferReceiver) error {

	return repository.AddProductOffer(productOffer)

}
func AddCategoryOffer(categoryOffer models.CategoryOfferReceiver) error {

	return repository.AddCategoryOffer(categoryOffer)

}
