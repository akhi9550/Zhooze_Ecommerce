package domain

type Wallet struct {
	ID     int     `json:"id"  gorm:"unique;not null"`
	UserID int     `json:"user_id"`
	Users  User    `json:"-" gorm:"foreignkey:UserID"`
	Amount float64 `json:"amount" gorm:"default:0"`
}
type WalletHistory struct {
	ID      int     `json:"id"  gorm:"unique;not null"`
	UserID  int     `json:"user_id"`
	OrderID int     `json:"order_id"`
	Reason  string  `json:"reason"`
	Amount  float64 `json:"amount"`
}
