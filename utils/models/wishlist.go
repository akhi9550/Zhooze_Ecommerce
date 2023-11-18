package models

type WishListResponse struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}
