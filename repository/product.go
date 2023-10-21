package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
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
func AddProducts(product models.ProductReceiver) (domain.Product, error) {
	var id int
	err := db.DB.Raw("INSERT INTO products (name, description, category_id, sku, size, brand_id, quantity, price,product_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id", product.Name, product.Description, product.CategoryID, product.SKU, product.Size, product.BrandID, product.Quantity, product.Price, product.ProductStatus).Scan(&id).Error
	if err != nil {
		return domain.Product{}, err
	}
	var ProductResponses domain.Product

	err = db.DB.Raw(`SELECT p.id, p.name, p.description, c.id, p.sku, p.size, b.brand_id, p.quantity, p.price, p.product_status FROM products p INNER JOIN categories c ON p.category_id = c.id INNER JOIN brands b ON p.brand_id = b.id WHERE p.id = ?`, ProductResponses.ID).Scan(&ProductResponses).Error
	if err != nil {
		return domain.Product{}, err
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
func UpdateProduct(id uint, product models.ProductReceiver) error {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("product for given id does not exist")
	}
	var reciever models.ProductBrief
	if err := db.DB.Raw(`update products set name = $1,
	description = $2,
	category_id = $3,
	sku = $4,
	size = $5,
	brand_id = $6,
	quantity = $7,
	price = $8
	where id = $9 returning id`,
		product.Name,
		product.Description,
		product.CategoryID,
		product.SKU,
		product.Size,
		product.BrandID,
		product.Quantity,
		product.Price,
		reciever.ID).Error; err != nil {
		return err
	}
	return nil

}
