package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"errors"
	"strconv"
)

func ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var productBrief []models.ProductBrief
	err := db.DB.Raw(`SELECT * FROM products limit ? offset ?`, count, offset).Scan(&productBrief).Error
	if err != nil {
		return nil, err
	}
	return productBrief, nil
}
func CheckValidateCategory(data map[string]int) error {
	for _, id := range data {
		var count int
		err := db.DB.Raw("SELECT COUNT(*) FROM categories WHERE id=?", id).Scan(&count).Error
		if err != nil {
			return err
		}
		if count < 1 {
			return errors.New("doesn't exist")
		}
	}
	return nil
}
func GetProductFromCategory(id int) ([]models.ProductBrief, error) {
	var product []models.ProductBrief
	err := db.DB.Raw(`SELECT * FROM products JOIN categories ON products.category_id=categories.id WHERE categories.id=?`, id).Scan(&product).Error
	if err != nil {
		return []models.ProductBrief{}, err
	}
	return product, nil
}
func GetQuantityFromProductID(id int) (int, error) {
	var quantity int
	err := db.DB.Raw("SELECT quantity FROM products WHERE id= ?", id).Scan(&quantity).Error
	if err != nil {
		return 0.0, err
	}
	return quantity, nil
}
func SeeAllProducts() ([]models.ProductBrief, error) {
	var products []models.ProductBrief
	err := db.DB.Raw("SELECT * FROM products").Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func AddProducts(product models.ProductReceiver) (models.ProductResponse, error) {
	var id int
	err := db.DB.Raw("INSERT INTO products (name, description, category_id,category_name, sku, size, brand_id, quantity, price) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id", product.Name, product.Description, product.CategoryID, product.CategoryName, product.SKU, product.Size, product.BrandID, product.Quantity, product.Price).Scan(&id).Error
	if err != nil {
		return models.ProductResponse{}, err
	}
	var ProductResponses models.ProductResponse
	err = db.DB.Raw(`SELECT id, name, description,category_name, sku, size, brand_id, quantity, price FROM products JOIN categories ON products.category_id = categories.id WHERE products.id=?`, id).Scan(&ProductResponses).Error
	if err != nil {
		return models.ProductResponse{}, err
	}
	return ProductResponses, nil
}

func DeleteProducts(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", product_id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("product for given id does not exist")
	}
	if err := db.DB.Exec("DELETE FROM products WHERE id=?", product_id).Error; err != nil {
		return err
	}
	return nil
}
