package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"fmt"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {
	adminCompareDetails, err := repository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		fmt.Println(err)
		return domain.TokenAdmin{}, err
	}
	var adminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	tokenString, err := helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	return domain.TokenAdmin{
		Admin: adminDetailsResponse,
		Token: tokenString,
	}, nil
}

func DashBoard() (models.CompleteAdminDashboard, error) {
	userDetails, err := repository.DashBoardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := repository.DashBoardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
	}, nil
}
