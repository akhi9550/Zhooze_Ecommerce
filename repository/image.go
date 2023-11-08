package repository

import "Zhooze/db"

func GetImageUrl(productImageID int) (string, error) {
	var imageUrl string
	if err := db.DB.Raw("SELECT url FROM images WHERE id = ?", productImageID).Scan(&imageUrl).Error; err != nil {
		return "", err
	}
	return imageUrl, nil
}
