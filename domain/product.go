package domain

type Product struct {
	ID            uint     `json:"id" gorm:"unique;not null"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	CategoryID    uint     `json:"category_id"`
	Category      Category `json:"-" gorm:"foreignkey:CategoryID;"`
	SKU           string   `json:"sku"`
	Size          int      `json:"size"`
	Stock         int      `json:"stock"`
	Price         float64  `json:"price"`
	Product_image []byte   `json:"product_image"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category" gorm:"unique; not null"`
}
type ProductImages struct {
	ID              uint   `json:"id" gorm:"unique; not null"`
	ProductImageUrl string `json:"product_image_url"`
}
