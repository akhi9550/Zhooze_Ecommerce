package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/utils/models"
	"errors"
	"log"
	"strconv"
)

func ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	if page <= 0 {
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

	if err := db.DB.Raw("SELECT price FROM products WHERE id = ?", prodcut_id).Scan(&productPrice).Error; err != nil {
		return 0.0, err
	}
	return productPrice, nil
}

func ProductAlreadyExist(Name string) bool {
	var count int
	if err := db.DB.Raw("SELECT count(*) FROM products WHERE name = ?", Name).Scan(&count).Error; err != nil {
		return false
	}
	return count > 0
}
func FindCategoryID(id int) (int, error) {
	var a int
	if err := db.DB.Raw("SELECT category_id FROM products WHERE id = ?", id).Scan(&a).Error; err != nil {
		return 0.0, err
	}
	return a, nil
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
    INSERT INTO products (name, description, category_id, size, stock, price)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING name, description, category_id, size, stock, price`
	err := db.DB.Raw(query, product.Name, product.Description, product.CategoryID, product.Size, product.Stock, product.Price).Scan(&p).Error
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
	if err := db.DB.Exec("DELETE FROM images WHERE product_id = ?", product_id).Error; err != nil {
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

func UpdateProductImage(productID int, url string) error {
	err := db.DB.Exec("INSERT INTO images (product_id,url) VALUES ($1,$2) RETURNING * ", productID, url).Error
	if err != nil {
		return errors.New("error while insert image to database")
	}
	return nil
}
func DisplayImages(productID int) (domain.Product, []domain.Image, error) {
	var product domain.Product
	var image []domain.Image
	err := db.DB.Raw(`SELECT * FROM products WHERE product_id = $1`, productID).Scan(&product).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	err = db.DB.Raw(`SELECT * FROM images WHERE product_id = $1`, productID).Scan(&image).Error
	if err != nil {
		return domain.Product{}, []domain.Image{}, err
	}
	return product, image, nil
}
func GetImage(productID int) ([]string, error) {
	var url []string
	if err := db.DB.Raw(`SELECT url FROM Images WHERE product_id=?`, productID).Scan(&url).Error; err != nil {
		return url, err
	}
	return url, nil

}
