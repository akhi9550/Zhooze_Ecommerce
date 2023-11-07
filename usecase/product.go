package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	products, err := repository.ShowAllProducts(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range products {
		p := &products[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	return products, nil
}
func ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	products, err := repository.ShowAllProductsFromAdmin(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range products {
		p := &products[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	return products, nil
}
func FilterCategory(data map[string]int) ([]models.ProductBrief, error) {
	err := repository.CheckValidateCategory(data)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	var ProductFromCategory []models.ProductBrief
	for _, id := range data {
		product, err := repository.GetProductFromCategory(id)
		if err != nil {
			return []models.ProductBrief{}, err
		}
		for _, products := range product {
			stock, err := repository.GetQuantityFromProductID(int(products.ID))
			if err != nil {
				return []models.ProductBrief{}, err
			}
			if stock <= 0 {
				products.ProductStatus = "out of stock"
			} else {
				products.ProductStatus = "in stock"
			}
			if products.ID != 0 {
				ProductFromCategory = append(ProductFromCategory, products)
			}
		}
	}
	return ProductFromCategory, nil
}

//	func SeeAllProducts() ([]models.ProductBrief, error) {
//		products, err := repository.SeeAllProducts()
//		if err != nil {
//			return []models.ProductBrief{}, err
//		}
//		for i := range products {
//			p := &products[i]
//			if p.Stock <= 0 {
//				p.ProductStatus = "out of stock"
//			} else {
//				p.ProductStatus = "in stock"
//			}
//		}
//		return products, nil
//	}
func AddProducts(product models.Product) (domain.Product, error) {
	exist := repository.ProductAlreadyExist(product.Name)
	if exist {
		return domain.Product{}, errors.New("product already exist")
	}
	productResponse, err := repository.AddProducts(product)
	if err != nil {
		return domain.Product{}, err
	}
	stock := repository.StockInvalid(productResponse.Name)
	if !stock {
		return domain.Product{}, errors.New("stock is invalid input")
	}
	return productResponse, nil
}
func DeleteProducts(id string) error {
	err := repository.DeleteProducts(id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateProduct(pid int, stock int) (models.ProductUpdateReciever, error) {
	if stock <= 0 {
		return models.ProductUpdateReciever{}, errors.New("stock doesnot update invalid input")
	}
	result, err := repository.CheckProductExist(pid)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	if !result {
		return models.ProductUpdateReciever{}, errors.New("there is no product as you mentioned")
	}
	newcat, err := repository.UpdateProduct(pid, stock)
	if err != nil {
		return models.ProductUpdateReciever{}, err
	}
	return newcat, err

}
func AddImage(c *gin.Context, file *multipart.FileHeader, productID int) (domain.Image, error) {
	if err := c.SaveUploadedFile(file, "uploads/"+file.Filename); err != nil {
		return domain.Image{}, errors.New("failed to saving file")
	}
	baseUrl := "http://localhost:8000"
	uploadedURL := baseUrl + "/uploads/" + file.Filename
	url, err := repository.SendUrl( uploadedURL,productID)
	if err != nil {
		return domain.Image{}, err
	}
	return url, nil

}
