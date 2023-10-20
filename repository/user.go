package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"

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
func AddAddress(userID int, address models.AddressInfo) error {
	err := db.DB.Exec("INSERT INTO addresses(user_id,name,house_name,street,city,state,pin)VALUES(?,?,?,?,?,?,?)", userID, address.Name, address.HouseName, address.Street, address.City, address.State, address.Pin).Error
	if err != nil {
		return errors.New("could not add address")
	}
	return nil
}
func GetAllAddress(userId int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	if err := db.DB.Raw("select * from addresses where user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
func UserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	err := db.DB.Raw("SELECT users.firstname,users.lastname,users.email,user.phone FROM users WHERE users.id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Email, &userDetails.Phone)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return userDetails, nil
}
func CheckUserAvailabilityWithUserID(userID int) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE id= ?", userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func UpdateUserEmail(email string, userID int) error {
	err := db.DB.Exec("UPDATE users SET email= ? WHERE id = ?", email, userID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateUserPhone(phone string, userID int) error {
	if err := db.DB.Exec("UPDATE users SET phone = ? WHERE id = ?", phone, userID).Error; err != nil {
		return err
	}
	return nil
}
func UpdateFirstName(name string, userID int) error {

	err := db.DB.Exec("update users set firstname = ? where id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func UpdateLastName(name string, userID int) error {

	err := db.DB.Exec("update users set lastname = ? where id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
