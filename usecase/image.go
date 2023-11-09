package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
)

func ShowImages(productID int) ([]models.Image, error) {
	image, err := repository.ShowImages(productID)
	if err != nil {
		return nil, err
	}
	return image, nil
}
func GetImageURL(productImageId int) (string, error) {
	imageUrl, err := repository.GetImageUrl(productImageId)
	if err != nil {
		return "", err
	}
	return imageUrl, nil
}
