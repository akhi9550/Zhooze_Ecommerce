package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
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
func ShowAllUsers() ([]models.UserDetailsResponse, error) {
	users, err := repository.ShowAllUsersIn()
	if err != nil {
		return []models.UserDetailsResponse{}, err
	}
	return users, nil
}
func BlockedUser(id string) error {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = repository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
func UnBlockedUser(id string) error {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}
	err = repository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
