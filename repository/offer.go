package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"errors"
	"time"
)

func AddProductOffer(productOffer models.ProductOfferReceiver) error {
	// check if the offer with the offer name already exist in the db
	var count int
	err := db.DB.Raw("select count(*) from product_offers where offer_name = ? and product_id = ?", productOffer.OfferName, productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the offer already exists")
	}

	// if there is any other offer for this product delete that before adding this one
	count = 0
	err = db.DB.Raw("select count(*) from product_offers where product_id = ?", productOffer.ProductID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		err = db.DB.Exec("delete from product_offers where product_id = ?", productOffer.ProductID).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = db.DB.Exec("INSERT INTO product_offers (product_id, offer_name, discount_percentage, start_date, end_date, offer_limit,offer_used) VALUES (?, ?, ?, ?, ?, ?, ?)", productOffer.ProductID, productOffer.OfferName, productOffer.DiscountPercentage, startDate, endDate, productOffer.OfferLimit, 0).Error
	if err != nil {
		return err
	}

	return nil

}
func AddCategoryOffer(categoryOffer models.CategoryOfferReceiver) error {

	// check if the offer with the offer name already exist in the db
	var count int
	err := db.DB.Raw("select count(*) from category_offers where offer_name = ?", categoryOffer.OfferName).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("the offer already exists")
	}

	// if there is any other offer for this category delete that before adding this one
	count = 0
	err = db.DB.Raw("select count(*) from category_offers where category_id = ?", categoryOffer.CategoryID).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {

		err = db.DB.Exec("delete from category_offers where category_id = ?", categoryOffer.CategoryID).Error
		if err != nil {
			return err
		}
	}

	startDate := time.Now()
	endDate := time.Now().Add(time.Hour * 24 * 5)
	err = db.DB.Exec("INSERT INTO category_offers (category_id, offer_name, discount_percentage, start_date, end_date, offer_limit,offer_used) VALUES (?, ?, ?, ?, ?, ?, ?)", categoryOffer.CategoryID, categoryOffer.OfferName, categoryOffer.DiscountPercentage, startDate, endDate, categoryOffer.OfferLimit, 0).Error
	if err != nil {
		return err
	}

	return nil

}
