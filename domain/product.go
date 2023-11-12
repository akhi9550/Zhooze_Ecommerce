package domain

import "time"

type Product struct {
	ID          uint     `json:"id" gorm:"unique;not null"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"-" gorm:"foreignkey:CategoryID;"`
	Size        int      `json:"size"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
}
type Category struct {
	ID       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category" gorm:"unique; not null"`
}
type ProductImages struct {
	ID              uint   `json:"id" gorm:"unique; not null"`
	ProductImageUrl string `json:"product_image_url"`
}
type Image struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductId uint   `json:"product_id"`
	Url       string `JSON:"url" `
}
type ProductOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	ProductID          uint      `json:"product_id"`
	Products           Product  `json:"-" gorm:"foreignkey:ProductID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	OfferLimit         int       `json:"offer_limit"`
	OfferUsed          int       `json:"offer_used"`
}

type CategoryOffer struct {
	ID                 uint      `json:"id" gorm:"unique; not null"`
	CategoryID         uint      `json:"category_id"`
	Category           Category  `json:"-" gorm:"foreignkey:CategoryID"`
	OfferName          string    `json:"offer_name"`
	DiscountPercentage int       `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	OfferLimit         int       `json:"offer_limit"`
	OfferUsed          int       `json:"offer_used"`
}
