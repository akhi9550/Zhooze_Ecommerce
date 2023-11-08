package models

type Image struct {
	Url string `JSON:"url" `
}

type ProductBrief struct {
	ID            uint    `json:"id" gorm:"unique;not null"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Size          int     `json:"size"`
	Stock         int     `json:"stock"`
	Price         float64 `json:"price"`
	ProductStatus string  `json:"product_status"`
	ProductImages []Image `json:"product_images"`
}

type ProductReceiver struct {
	Name        string  `json:"name" `
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Size        int     `json:"size"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}
type Product struct {
	Name        string  `json:"name" `
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Size        int     `json:"size"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}
type Category struct {
	Category string `json:"category"`
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
