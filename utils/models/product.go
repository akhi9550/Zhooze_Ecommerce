package models

type ProductBrief struct {
	ID          uint    `json:"id" gorm:"unique;not null"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	SKU         string  `json:"sku"`
	Size        int     `json:"size"`
	BrandID     uint    `json:"brand_id"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	// Image         string  `json:"image" gorm:"not null"`
	ProductStatus string `json:"product_status"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Category string `json:"category_name"`
}
