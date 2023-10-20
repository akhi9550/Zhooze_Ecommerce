package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"strconv"
)

func LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {
	var details domain.Admin
	if err := db.DB.Raw("SELECT * FROM users WHERE email=? AND isadmin= true", adminDetails.Email).Scan(&details).Error; err != nil {
		return domain.Admin{}, err
	}
	return details, nil
}
func DashBoardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE isadmin='false'").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	err = db.DB.Raw("SELECT COUNT(*)  FROM users WHERE blocked=true").Scan(&userDetails.BlockedUser).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	return userDetails, nil
}

func DashBoardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := db.DB.Raw("SELECT COUNT(*) FROM products").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	err = db.DB.Raw("SELECT COUNT(*) FROM products WHERE quantity=0").Scan(&productDetails.OutofStockProduct).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	return productDetails, nil
}
func ShowAllUsersIn() ([]models.UserDetailsResponse, error) {
	var user []models.UserDetailsResponse
	err := db.DB.Raw("SELECT * FROM users WHERE isadmin='false'").Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func GetUserByID(id string) (domain.User, error) {
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, err
	}
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE id=?", user_id).Scan(&count).Error; err != nil {

		return domain.User{}, err
	}
	if count < 1 {
		return domain.User{}, errors.New("user for the given id does not exist")

	}
	var userDetails domain.User
	if err := db.DB.Raw("SELECT * FROM users WHERE id=?", user_id).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}
	return userDetails, nil
}

func UpdateBlockUserByID(user domain.User) error {
	err := db.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}
