package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func FindUserByPhoneNumber(phone string) (*domain.User, error) {
	var user domain.User
	res := db.DB.Where(&domain.User{Phone: phone}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error

	}
	return &user, nil
}

func UserDetailsUsingPhone(phone string) (models.UserDetailsResponse, error) {
	var userDeatils models.UserDetailsResponse
	if err := db.DB.Raw("SELECT * FROM users WHERE phone = ?", phone).Scan(&userDeatils).Error; err != nil {
		return models.UserDetailsResponse{}, err
	}

	return userDeatils, nil
}

func FindUsersByEmail(email string) (bool, error) {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserPhoneByEmail(email string) (string, error) {
	fmt.Println(email)
	var phone string
	if err := db.DB.Raw("SELECT phone FROM users WHERE email = ?", email).Scan(&phone).Error; err != nil {
		return "", nil
	}
	return phone, nil
}
