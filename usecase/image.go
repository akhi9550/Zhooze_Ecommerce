package usecase

import "Zhooze/repository"

func CropImage(productImageId int) (string, error) {
	imageUrl, err := repository.GetImageUrl(productImageId)
	if err!=nil{
		return "",err
	}
	return imageUrl,nil
}