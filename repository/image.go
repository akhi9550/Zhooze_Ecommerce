package repository

import "Zhooze/db"

func GetImageUrl(productImageID int) (string, error) {
	var imageUrl string
	if err := db.DB.Raw("SELECT product_image_url FROM product_images WHERE id = ?", productImageID).Scan(&imageUrl).Error; err != nil {
		return "", err
	}
	return imageUrl, nil
}
