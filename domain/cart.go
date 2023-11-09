package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     uint    `json:"user_id" gorm:"uniquekey; not null"`
	Users      User    `json:"-" gorm:"foreignkey:UserID"`
	ProductID  uint    `json:"product_id"`
	Products   Product `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

// type Cart struct {
// 	gorm.Model
// 	UserId int  `json:"user_id" gorm:"not null"`
// 	Users  User `json:"-" gorm:"foreignkey:UserId"`
// }

// type CartItems struct {
// 	ID         int     `json:"id" gorm:"primarykey;not null"`
// 	CartID     int     `json:"cart_id" gorm:"not null"`
// 	Cart       Cart    `json:"-" gorm:"foreignkey:CartID;constraint:OnDelete:CASCADE"`
// 	ProductId  int     `json:"product_id"`
// 	Product    Product `json:"-" gorm:"foreignkey:ProductId"`
// 	Quantity   float64 `json:"quantity"`
// 	TotalPrice float64 `json:"total_price"`
// }
