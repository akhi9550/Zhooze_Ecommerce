package models

type ProductBrief struct {
	ID           uint   `json:"id" gorm:"unique;not null"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	
	SKU          string `json:"sku"`
	Size         int    `json:"size"`
	BrandID      uint   `json:"brand_id"`
	Quantity     int    `json:"quantity"`
	// Stock       int     `json:"stock"`
	Price float64 `json:"price"`
	// Image         string  `json:"image" gorm:"not null"`
	ProductStatus string `json:"product_status"`
}
type ProductResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	CategoryName string  `json:"category_name"`
	SKU          string  `json:"sku"`
	Size         int     `json:"size"`
	BrandID      uint    `json:"brand_id"`
	Quantity     int     `json:"quantity"`
	Stock        int     `json:"stock"`
	Price        float64 `json:"price"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Category string `json:"category_name"`
}
type UpdateProduct struct {
	Quantity  int `json:"quantity" binding:"required"`
	ProductID int `json:"product-id" binding:"required"`
}
