package models

type UserSignUp struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}
type UserDetailsResponse struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
type TokenUser struct {
	Users        UserDetailsResponse
	AccessToken  string
	RefreshToken string
}
type LoginDetail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}
type AddressInfoResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
}
type AddressInfo struct {
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
type UsersProfileDetails struct {
	Firstname string `json:"firstname" `
	Lastname  string `json:"lastname" `
	Email     string `json:"email" `
	Phone     string `json:"phone" `
}

type UpdatePassword struct {
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

type PaymentDetails struct {
	ID           uint   `json:"id"`
	Payment_Name string `json:"payment_name"`
}

type CheckoutDetails struct {
	AddressInfoResponse []AddressInfoResponse
	Payment_Method      []PaymentDetails
	Cart                []Cart
	Grand_Total         float64
	Total_Price         float64
}
