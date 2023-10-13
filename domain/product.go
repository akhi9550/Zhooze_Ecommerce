package domain

type Products struct {
	ID          uint     `json:"id" gorm:"unique;not null"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"-" gorm:"foreignkey:GenreID;constraint:OnDelete:CASCADE"`
	SKU         string   `json:"sku"`
	Size        int      `json:"size"`
	BrandID     uint     `json:"brand_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	ProductStatus string `json:"product_status"`
	IsDeleted     bool   `json:"is_deleted" gorm:"default:false"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Category string `json:"category"`
}
