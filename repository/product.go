package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
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
func ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
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
func SeeAllProducts() ([]models.ProductBrief, error) {
	var products []models.ProductBrief
	err := db.DB.Raw("SELECT * FROM products").Scan(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func ShowIndividualProducts(id string) (*models.ProductBrief, error) {
	var product models.ProductBrief
	result := db.DB.Raw(`SELECT p.id,  p.name,  p.sku, p.size , c.category,  p.stock, p.price FROM products p JOIN  categories c ON p.category_id = c.id WHERE  p.sku=?`, id).Scan(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &product, nil
}
func ProductAlreadyExist(Name string) bool {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM products WHERE name = ?", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func StockInvalid(Name string) bool {
	var count int
	if err := db.DB.Raw("SELECT SUM(stock) FROM products WHERE name = ? AND stock >= 0", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func AddProducts(product models.Product) (domain.Product, error) {
	var p domain.Product
	query := `
    INSERT INTO products (name, description, category_id, sku, size, stock, price)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING name, description, category_id, sku, size, stock, price`
	err := db.DB.Raw(query, product.Name, product.Description, product.CategoryID, product.SKU, product.Size, product.Stock, product.Price).Scan(&p).Error
	fmt.Println("dkddkdkd", p)
	if err != nil {
		log.Println(err.Error())
		return domain.Product{}, err
	}
	var ProductResponses domain.Product
	err = db.DB.Raw("SELECT * FROM products WHERE name = ?", p.Name).Scan(&ProductResponses).Error
	if err != nil {
		log.Println(err.Error())
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
	if stock <= 0 {
		return models.ProductUpdateReciever{}, errors.New("stock doesnot update invalid input")
	}
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
func DoesProductExist(productID int) (bool, error) {
	var count int
	err := db.DB.Raw("select count(*) from products where id = ?", productID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
