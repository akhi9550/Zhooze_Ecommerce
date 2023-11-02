package usecase

import (
	"Zhooze/config"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"

	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func UsersSignUp(user models.UserSignUp) (*models.TokenUser, error) {
	email, err := repository.CheckUserExistsByEmail(user.Email)
	fmt.Println(email)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email != nil {
		return &models.TokenUser{}, errors.New("user with this email is already exists")
	}

	phone, err := repository.CheckUserExistsByPhone(user.Phone)
	fmt.Println(phone, nil)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if phone != nil {
		return &models.TokenUser{}, errors.New("user with this phone is already exists")
	}

	hashPassword, err := helper.PasswordHash(user.Password)
	if err != nil {
		return &models.TokenUser{}, errors.New("error in hashing password")
	}
	user.Password = hashPassword
	userData, err := repository.UserSignUp(user)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not add the user")
	}
	accessToken, err := helper.GenerateAccessToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create access token due to error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create refresh token due to error")
	}
	return &models.TokenUser{
		Users:        userData,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func UsersLogin(user models.LoginDetail) (*models.TokenUser, error) {
	email, err := repository.CheckUserExistsByEmail(user.Email)
	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email == nil {
		return &models.TokenUser{}, errors.New("email doesn't exist")
	}
	userdeatils, err := repository.FindUserByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userdeatils.Password), []byte(user.Password))
	if err != nil {
		return &models.TokenUser{}, errors.New("password not matching")
	}
	var userDetails models.UserDetailsResponse
	err = copier.Copy(&userDetails, &userdeatils)
	if err != nil {
		return &models.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(userDetails)
	if err != nil {
		return &models.TokenUser{}, errors.New("couldn't create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userDetails)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldn't create refreshtoken due to internal error")
	}
	return &models.TokenUser{
		Users:        userDetails,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func AddAddress(userID int, address models.AddressInfo) error {
	err := repository.AddAddress(userID, address)
	if err != nil {
		return err
	}
	return nil
}
func GetAllAddress(userId int) ([]models.AddressInfoResponse, error) {
	addressInfo, err := repository.GetAllAddress(userId)
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}
	return addressInfo, nil

}
func UserDetails(userID int) (models.UsersProfileDetails, error) {
	return repository.UserDetails(userID)
}
func UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error) {
	userExist := repository.CheckUserAvailabilityWithUserID(userID)
	if !userExist {
		return models.UsersProfileDetails{}, errors.New("user doesn't exist")
	}
	if userDetails.Email != "" {
		repository.UpdateUserEmail(userDetails.Email, userID)
	}
	if userDetails.Firstname != "" {
		repository.UpdateFirstName(userDetails.Firstname, userID)
	}
	if userDetails.Lastname != "" {
		repository.UpdateLastName(userDetails.Lastname, userID)
	}

	if userDetails.Phone != "" {
		repository.UpdateUserPhone(userDetails.Phone, userID)
	}
	return repository.UserDetails(userID)
}

func UpdateAddress(addressDetails models.AddressInfo, addressID, userID int) (models.AddressInfoResponse, error) {
	addressExist := repository.CheckAddressAvailabilityWithAddressID(addressID, userID)
	if !addressExist {
		return models.AddressInfoResponse{}, errors.New("address doesn't exist")
	}
	if addressDetails.Name != "" {
		repository.UpdateName(addressDetails.Name, addressID)
	}
	if addressDetails.HouseName != "" {
		repository.UpdateHouseName(addressDetails.HouseName, addressID)
	}
	if addressDetails.Street != "" {
		repository.UpdateStreet(addressDetails.Street, addressID)
	}
	if addressDetails.City != "" {
		repository.UpdateCity(addressDetails.City, addressID)
	}
	if addressDetails.State != "" {
		repository.UpdateState(addressDetails.State, addressID)
	}
	if addressDetails.Pin != "" {
		repository.UpdatePin(addressDetails.Pin, addressID)
	}
	return repository.AddressDetails(addressID)
}

func DeleteAddress(addressID, userID int) error {
	addressExist, err := repository.AddressExistInUserProfile(addressID, userID)
	if err != nil {
		return err
	}
	if !addressExist {
		return errors.New("address does not exist in user profile")
	}
	err = repository.RemoveFromUserProfile(userID, addressID)
	if err != nil {
		return err
	}
	return nil
}

func ChangePassword(id int, old string, password string, repassword string) error {
	userPassword, err := repository.GetPassword(id)
	if err != nil {
		return errors.New("internal error")
	}
	err = helper.CompareHashAndPassword(userPassword, old)
	if err != nil {
		return errors.New("password incorrect")
	}
	if password != repassword {
		return errors.New("password doesn't match")
	}
	newpassword, err := helper.PasswordHash(password)
	if err != nil {
		return errors.New("error in hashing password")
	}
	return repository.ChangePassword(id, string(newpassword))
}
func UpdateQuantityAdd(id, prodcut_id int) error {

	err := repository.UpdateQuantityAdd(id, prodcut_id)
	if err != nil {
		return err
	}

	return nil

}
func UpdateTotalPriceAdd(id, product_id int) error {
	err := repository.UpdateTotalPrice(id, product_id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateQuantityless(id, prodcut_id int) error {

	err := repository.UpdateQuantityless(id, prodcut_id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateTotalPriceLess(id, product_id int) error {
	err := repository.UpdateTotalPrice(id, product_id)
	if err != nil {
		return err
	}
	return nil
}

func ForgotPasswordSend(phone string) error {
	cfg, _ := config.LoadConfig()
	ok := repository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	_, err := helper.TwilioSendOTP(phone, cfg.SERVICESSID)
	if err != nil {
		return errors.New("error ocurred while generating OTP")
	}
	return nil
}

func ForgotPasswordVerifyAndChange(model models.ForgotVerify) error {
	cfg, _ := config.LoadConfig()
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err := helper.TwilioVerifyOTP(cfg.SERVICESSID, model.Otp, model.Phone)
	if err != nil {
		return errors.New("error while verifying")
	}

	id, err := repository.FindIdFromPhone(model.Phone)
	if err != nil {
		return errors.New("cannot find user from mobile number")
	}

	newpassword, err := helper.PasswordHashing(model.NewPassword)
	if err != nil {
		return errors.New("error in hashing password")
	}

	// if user is authenticated then change the password i the database
	if err := repository.ChangePassword(id, string(newpassword)); err != nil {
		return errors.New("could not change password")
	}

	return nil
}
