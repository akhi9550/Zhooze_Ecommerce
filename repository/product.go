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
	err := db.DB.Raw("SELECT stock FROM products WHERE id= ?", id).Scan(&quantity).Error
	if err != nil {
		return 0.0, err
	}
	return quantity, nil
}
func GetPriceOfProductFromID(prodcut_id int) (float64, error) {
	var productPrice float64

	if err := db.DB.Raw("select price from products where id = ?", prodcut_id).Scan(&productPrice).Error; err != nil {
		return 0.0, err
	}
	return productPrice, nil
}
func SeeAllProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := db.DB.Raw("SELECT * FROM products").Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
func AddProducts(product domain.Product) (domain.Product, error) {
	var p models.ProductReceiver
	err := db.DB.Raw("INSERT INTO products (name, description, category_id, sku, size, brand_id, stock, price,product_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING name, description, category_id, sku, size, brand_id, stock, price,product_status", product.Name, product.Description, product.CategoryID, product.SKU, product.Size, product.BrandID, product.Stock, product.Price, product.ProductStatus).Scan(&p).Error
	if err != nil {
		return domain.Product{}, err
	}
	var ProductResponses domain.Product

	err = db.DB.Raw("SELECT * FROM products WHERE products.name = ?", p.Name).Scan(&ProductResponses).Error
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
func CheckProductExist(pid int) (bool, error) {
	var a int
	err := db.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", pid).Scan(&a).Error
	if err != nil {
		return false, err
	}
	if a == 0 {
		return false, err
	}
	return true, err
}
func UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error) {
	if db.DB == nil {
		return models.ProductUpdateReciever{}, errors.New("database connection is nil")
	}
	if err := db.DB.Exec("UPDATE products SET stock = stock + $1 WHERE id = $2", stock, pid).Error; err != nil {
		return models.ProductUpdateReciever{}, err
	}
	var newdetails models.ProductUpdateReciever
	var newQuantity int
	if err := db.DB.Raw("SELECT stock FROM products WHERE id =?", pid).Scan(&newQuantity).Error; err != nil {
		return models.ProductUpdateReciever{}, err
	}
	newdetails.ProductID = pid
	newdetails.Stock = newQuantity
	return newdetails, nil
}
