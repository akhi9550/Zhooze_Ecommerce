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
	if err := db.DB.Raw("SELECT * FROM addresses WHERE user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
func UserDetails(userID int) (models.UsersProfileDetails, error) {
	var userDetails models.UsersProfileDetails
	err := db.DB.Raw("SELECT u.firstname,u.lastname,u.email,u.phone FROM users u WHERE u.id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Email, &userDetails.Phone)
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
func CheckAddressAvailabilityWithAddressID(addressID, userID int) bool {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM addresses WHERE id = ? AND user_id = ?", addressID, userID).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func UpdateName(name string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET name= ? WHERE id = ?", name, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateHouseName(HouseName string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET house_name= ? WHERE id = ?", HouseName, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateStreet(street string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET street= ? WHERE id = ?", street, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateCity(city string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET city= ? WHERE id = ?", city, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateState(state string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET state= ? WHERE id = ?", state, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdatePin(pin string, addressID int) error {
	err := db.DB.Exec("UPDATE addresses SET pin= ? WHERE id = ?", pin, addressID).Error
	if err != nil {
		return err
	}
	return nil
}
func AddressDetails(addressID int) (models.AddressInfoResponse, error) {
	var addressDetails models.AddressInfoResponse
	err := db.DB.Raw("SELECT a.id, a.name, a.house_name, a.street, a.city, a.state, a.pin FROM addresses a WHERE a.id = ?", addressID).Row().Scan(&addressDetails.ID, &addressDetails.Name, &addressDetails.HouseName, &addressDetails.Street, &addressDetails.City, &addressDetails.State, &addressDetails.Pin)
	if err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressDetails, nil
}

func ChangePassword(id int, password string) error {
	err := db.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", password, id).Error
	if err != nil {
		return err
	}
	return nil
}
func GetPassword(id int) (string, error) {
	var userPassword string
	err := db.DB.Raw("SELECT password FROM users WHERE id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil
}
func UpdateQuantityAdd(id, prdt_id int) error {
	err := db.DB.Exec("UPDATE Carts SET quantity = quantity + 1 WHERE id=$1 AND product_id = $2 ", id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateQuantityless(id, prdt_id int) error {
	err := db.DB.Exec("UPDATE Carts SET quantity = quantity - 1 WHERE id=$1 AND product_id = $2 ", id, prdt_id).Error
	if err != nil {
		return err
	}
	return nil
}
func FindUserByMobileNumber(phone string) bool {

	var count int
	if err := db.DB.Raw("SELECT count(*) FROM users WHERE phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0

}
func FindIdFromPhone(phone string) (int, error) {
	var id int
	if err := db.DB.Raw("SELECT id FROM users WHERE phone=?", phone).Scan(&id).Error; err != nil {
		return id, err
	}
	return id, nil
}
