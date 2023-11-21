package models

type WalletAmount struct {
	Amount float64 `json:"amount"`
}

type WalletHistory struct {
	ID      int     `json:"id"  gorm:"unique;not null"`
	OrderID int     `json:"order_id"`
	Reason  string  `json:"reason"`
	Amount  float64 `json:"amount"`
}
