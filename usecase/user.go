package usecase

import (
	"Zhooze/config"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"strconv"

	"errors"
	"fmt"

	"github.com/google/uuid"
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
	// create referral code for the user and send in details of referred id of user if it exist
	id := uuid.New().ID()
	str := strconv.Itoa(int(id))
	userReferral := str[:8]
	err = repository.CreateReferralEntry(userData, userReferral)
	if err != nil {
		return &models.TokenUser{}, err
	}
	if user.ReferralCode != "" {
		// first check whether if a user with that referralCode exist
		referredUserId, err := repository.GetUserIdFromReferrals(user.ReferralCode)
		if err != nil {
			return &models.TokenUser{}, err
		}
		if referredUserId != 0 {
			referralAmount := 150
			err := repository.UpdateReferralAmount(float64(referralAmount), referredUserId, userData.Id)
			if err != nil {
				return &models.TokenUser{}, err
			}
			referreason := "Amount credited for used referral code"
			err = repository.UpdateHistory(userData.Id, 0, float64(referralAmount), referreason)
			if err != nil {
				return &models.TokenUser{}, err
			}
			amount, err := repository.AmountInrefferals(userData.Id)
			if err != nil {
				return &models.TokenUser{}, err
			}
			wallectExist, err := repository.ExistWallect(userData.Id)
			if err != nil {
				return &models.TokenUser{}, err
			}
			if !wallectExist {
				err = repository.NewWallect(userData.Id, amount)
				if err != nil {
					return &models.TokenUser{}, err
				}
			}
			err = repository.UpdateReferUserWallect(amount, referredUserId)
			if err != nil {
				return &models.TokenUser{}, err
			}
			reason := "Amount credited for refer a new person"
			err = repository.UpdateHistory(referredUserId, 0, amount, reason)
			if err != nil {
				return &models.TokenUser{}, err
			}
		}
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

func UpdateQuantityAdd(id, productID int) error {
	productExist, err := repository.ProductExistCart(id, productID)
	if !productExist {
		return errors.New("product doesnot exist cart")
	}
	if err != nil {
		return err
	}
	stock, err := repository.ProductStock(productID)
	if err != nil {
		return err
	}
	if stock <= 0 {
		return errors.New("not available out of stock")
	}
	err = repository.UpdateQuantityAdd(id, productID)
	if err != nil {
		return err
	}
	stockfromcart, err := repository.StockFormCart(productID)
	if err != nil {
		return err
	}
	if stock <= stockfromcart {
		return errors.New("its maximum no more updation")
	}
	productPrice, err := repository.GetPriceOfProductFromID(productID)
	if err != nil {

		return err
	}
	discount_percentage, err := repository.FindDiscountPercentageForProduct(productID)
	if err != nil {
		return errors.New("there was some error in finding the discounted prices")
	}
	var discount float64

	if discount_percentage > 0 {
		discount = (productPrice * float64(discount_percentage)) / 100
	}

	Price := productPrice - discount
	categoryID, err := repository.FindCategoryID(productID)
	if err != nil {
		return err
	}
	discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(categoryID)
	if err != nil {
		return errors.New("there was some error in finding the discounted prices")
	}
	var discountcategory float64

	if discount_percentageCategory > 0 {
		discountcategory = (productPrice * float64(discount_percentageCategory)) / 100
	}

	FinalPrice := Price - discountcategory
	FinalPrice = FinalPrice * float64(stockfromcart)
	err = repository.UpdateTotalPrice(id, productID, FinalPrice)
	if err != nil {
		return err
	}
	return nil

}

func UpdateQuantityless(id, productID int) error {
	productExist, err := repository.ProductExistCart(id, productID)
	if !productExist {
		return errors.New("product doesnot exist cart")
	}
	if err != nil {
		return err
	}
	stock, err := repository.ExistStock(id, productID)
	if err != nil {
		return err
	}
	if stock <= 1 {
		return errors.New("its  minimum")
	}
	err = repository.UpdateQuantityless(id, productID)
	if err != nil {
		return err
	}
	stockfromcart, err := repository.StockFormCart(productID)
	if err != nil {
		return err
	}
	productPrice, err := repository.GetPriceOfProductFromID(productID)
	if err != nil {

		return err
	}
	discount_percentage, err := repository.FindDiscountPercentageForProduct(productID)
	if err != nil {
		return errors.New("there was some error in finding the discounted prices")
	}
	var discount float64

	if discount_percentage > 0 {
		discount = (productPrice * float64(discount_percentage)) / 100
	}

	Price := productPrice - discount
	categoryID, err := repository.FindCategoryID(productID)
	if err != nil {
		return err
	}
	discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(categoryID)
	if err != nil {
		return errors.New("there was some error in finding the discounted prices")
	}
	var discountcategory float64

	if discount_percentageCategory > 0 {
		discountcategory = (productPrice * float64(discount_percentageCategory)) / 100
	}

	FinalPrice := Price - discountcategory
	FinalPrice = FinalPrice * float64(stockfromcart)
	err = repository.UpdateTotalPrice(id, productID, FinalPrice)
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

	// if user is authenticated then change the password in the database
	if err := repository.ChangePassword(id, string(newpassword)); err != nil {
		return errors.New("could not change password")
	}

	return nil
}
