package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
)

func AddProductOffer(model models.ProductOfferReceiver) error {
	if err := repository.AddProductOffer(model); err != nil {
		return err
	}

	return nil
}
func GetOffers() ([]domain.ProductOffer, error) {

	offers, err := repository.GetOffers()
	if err != nil {
		return []domain.ProductOffer{}, err
	}
	return offers, nil

}


func AddCategoryOffer(model models.CategoryOfferReceiver) error {
	if err := repository.AddCategoryOffer(model); err != nil {
		return err
	}

	return nil
}
func GetCategoryOffer() ([]domain.CategoryOffer, error) {

	offers, err := repository.GetCategoryOffer()
	if err != nil {
		return []domain.CategoryOffer{}, err
	}
	return offers, nil

}
