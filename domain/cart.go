package domain

type Cart struct {
	ID     uint `json:"id" gorm:"primarykey"`
	UserID uint `json:"user_id" gorm:"not null"`
	Users  User `json:"-" gorm:"foreignkey:UserID"`
}

type LineItems struct {
	ID         uint    `json:"id" gorm:"primarykey"`
	CartID     uint    `json:"cart_id" gorm:"not null"`
	Cart       Cart    `json:"-" gorm:"foreignkey:CartID"`
	ProductID  uint    `json:"product_id"`
	Products   Product `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   int     `json:"quantity" gorm:"default:1"`
	TotalPrice float64 `json:"total_price"`
}
