package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
)

func GetImageUrl(productImageID int) (string, error) {
	var imageUrl string
	if err := db.DB.Raw("SELECT url FROM images WHERE id = ?", productImageID).Scan(&imageUrl).Error; err != nil {
		return "", err
	}
	return imageUrl, nil
}
func ShowImages(productID int) ([]models.Image, error) {
	var image []models.Image
	err := db.DB.Raw(`SELECT url FROM images  WHERE images.product_id = $1`, productID).Scan(&image).Error
	if err != nil {
		return nil, err
	}
	return image, nil
}
