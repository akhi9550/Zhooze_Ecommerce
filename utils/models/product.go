package models

type ProductBrief struct {
	ID            uint    `json:"id" gorm:"unique;not null"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	SKU           string  `json:"sku"`
	Size          int     `json:"size"`
	BrandID       uint    `json:"brand_id"`
	Stock         int     `json:"stock"`
	Price         float64 `json:"price"`
	ProductStatus string  `json:"product_status"`
}

type ProductReceiver struct {
	Name          string  `json:"name" `
	Description   string  `json:"description"`
	CategoryID    uint    `json:"category_id"`
	SKU           string  `json:"sku"`
	Size          int     `json:"size"`
	BrandID       uint    `json:"brand_id"`
	Stock         int     `json:"stock"`
	Price         float64 `json:"price"`
	ProductStatus string  `json:"product_status"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category_name"`
}
type SetNewName struct {
	Current string `json:"current"`
	New     string `json:"new"`
}
type ProductUpdate struct {
	ProductId int `json:"product_id"`
	Stock     int `json:"stock"`
}
type ProductUpdateReciever struct {
	ProductID int
	Stock     int
}
