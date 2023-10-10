package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
)

func LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {
	var details domain.Admin
	if err := db.DB.Raw("SELECT * FROM users WHERE email=? AND isadmin=true", adminDetails.Email).Scan(&details).Error; err != nil {
		return domain.Admin{}, err
	}
	return details, nil
}
func DashBoardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := db.DB.Raw("SELECT COUNT(*) FROM users").Scan(&userDetails.TotalUsers).Error
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
	err = db.DB.Raw("SELECT COUNT(*) FROM products WHERE stock=0").Scan(&productDetails.OutofStockProduct).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	return productDetails, nil
}
