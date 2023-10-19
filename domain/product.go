package domain

type Product struct {
	ID            uint     `json:"id" gorm:"unique;not null"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	CategoryID    uint     `json:"category_id"`
	Category      Category `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	SKU           string   `json:"sku"`
	Size          int      `json:"size"`
	BrandID       uint     `json:"brand_id" gorm:"foreignkey"`
	Stock         int      `json:"stock"`
	Price         float64  `json:"price"`
	ProductStatus string   `json:"product_status"`
	IsDeleted     bool     `json:"is_deleted" gorm:"default:false"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category_name" gorm:"unique; not null"`
}
type Brand struct {
	ID    uint   `json:"id" gorm:"unique; not null"`
	Brand string `json:"brand_name"`
}
