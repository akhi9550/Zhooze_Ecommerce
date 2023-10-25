package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	products, err := repository.ShowAllProducts(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range products {
		p := &products[i]
		if p.Quantity == 0 {
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
		for _, product := range product {
			quantity, err := repository.GetQuantityFromProductID(int(product.ID))
			if err != nil {
				return []models.ProductBrief{}, err
			}
			if quantity == 0 {
				product.ProductStatus = "out of stock"
			} else {
				product.ProductStatus = "in stock"
			}
			if product.ID != 0 {
				ProductFromCategory = append(ProductFromCategory, product)
			}
		}
	}
	return ProductFromCategory, nil
}
func SeeAllProducts() ([]domain.Product, error) {
	products, err := repository.SeeAllProducts()
	if err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}
func AddProducts(product domain.Product) (domain.Product, error) {
	productResponse, err := repository.AddProducts(product)
	if err != nil {
		return domain.Product{}, err
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
