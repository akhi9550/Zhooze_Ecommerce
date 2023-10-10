package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CheckUserExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	res := db.DB.Where(&domain.User{Email: email}).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return &user, nil
}
func CheckUserExistsByPhone(phone string) (*domain.User, error) {
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
func UserSignUp(user models.UserSignUp) (models.UserDetailsResponse, error) {
	var SignupDetail models.UserDetailsResponse
	err := db.DB.Raw("INSERT INTO users(firstname,lastname,email,password,phone)VALUES(?,?,?,?,?)RETURNING id,firstname,lastname,email,password,phone", user.Firstname, user.Lastname, user.Email, user.Password, user.Phone).Scan(&SignupDetail).Error
	if err != nil {
		fmt.Println("Respository error:", err)
		return models.UserDetailsResponse{}, err
	}
	return SignupDetail, nil
}
func FindUserByEmail(user models.LoginDetail) (models.UserLoginResponse, error) {
	var userDetails models.UserLoginResponse
	err := db.DB.Raw("SELECT * FROM users WHERE email=? and blocked=false", user.Email).Scan(&userDetails).Error
	if err != nil {
		return models.UserLoginResponse{}, errors.New("error checking user details")
	}
	return userDetails, nil
}
