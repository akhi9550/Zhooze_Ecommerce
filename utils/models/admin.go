package models

type AdminLogin struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"min=6,max=20"`
}

//	type AdminDetails struct {
//		ID    uint   `json:"id" gorm:"uniquekey; not null"`
//		Name  string `json:"name" gorm:"validate:required"`
//		Email string `json:"email" gorm:"validate:required"`
//	}
//
//	type AdminSignUp struct {
//		Name            string `json:"name" binding:"required" gorm:"validate:required"`
//		Email           string `json:"email" binding:"required" gorm:"validate:required"`
//		Password        string `json:"password" binding:"required" gorm:"validate:required"`
//		ConfirmPassword string `json:"confirmpassword" binding:"required"`
//	}
type AdminDetailsResponse struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"Email"`
}

type DashBoardUser struct {
	TotalUsers  int `json:"Totaluser"`
	BlockedUser int `json:"Blockuser"`
}
type DashBoardProduct struct {
	TotalProducts     int `json:"Totalproduct"`
	OutofStockProduct int `json:"Outofstock"`
}
type CompleteAdminDashboard struct {
	DashboardUser    DashBoardUser
	DashboardProduct DashBoardProduct
}
