package models

// type Cart struct {
// 	ProductID   uint    `json:"product_id"`
// 	ProductName string  `json:"product_name"`
// 	Quantity    float64 `json:"quantity"`
// 	TotalPrice  float64 `json:"total_price"`
// }

type CartResponse struct {
	UserName   string
	TotalPrice float64
	Cart       []Cart
}
type CartTotal struct {
	UserName       string  `json:"user_name"`
	TotalPrice     float64 `json:"total_price"`
	FinalPrice     float64 `json:"final_price"`
	DiscountReason string
}
type Cart struct {
	CartId     int     `json:"cart_id"`
	ProductId  int     `json:"product_id"`
	Quantity   float64 `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}
type Carts struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}
