package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"mime/multipart"
)

func ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	productDetails, err := repository.ShowAllProducts(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}

	for i := range productDetails {
		p := &productDetails[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	//loop inside products and then calculate discounted price of each then return
	for j := range productDetails {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil
}

func ShowAllProductsFromAdmin(page int, count int) ([]models.ProductBrief, error) {
	productDetails, err := repository.ShowAllProductsFromAdmin(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range productDetails {
		p := &productDetails[i]
		if p.Stock <= 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	for j := range productDetails {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(productDetails[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}
		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(productDetails[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (productDetails[j].Price * float64(discount_percentageCategory)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].DiscountedPrice - categorydiscount
	}
	var updatedproductDetails []models.ProductBrief
	for _, p := range productDetails {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil

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
	for j := range ProductFromCategory {
		discount_percentage, err := repository.FindDiscountPercentageForProduct(int(ProductFromCategory[j].ID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64
		if discount_percentage > 0 {
			discount = (ProductFromCategory[j].Price * float64(discount_percentage)) / 100
		}
		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].Price - discount

		discount_percentageCategory, err := repository.FindDiscountPercentageForCategory(int(ProductFromCategory[j].CategoryID))
		if err != nil {
			return []models.ProductBrief{}, errors.New("there was some error in finding the discounted prices")
		}
		var categorydiscount float64
		if discount_percentageCategory > 0 {
			categorydiscount = (ProductFromCategory[j].Price * float64(discount_percentageCategory)) / 100
		}

		ProductFromCategory[j].DiscountedPrice = ProductFromCategory[j].DiscountedPrice - categorydiscount
	}
	updatedproductDetails := make([]models.ProductBrief, 0)
	for _, p := range ProductFromCategory {
		img, err := repository.GetImage(int(p.ID))
		if err != nil {
			return nil, err
		}
		p.Image = img
		updatedproductDetails = append(updatedproductDetails, p)
	}

	return updatedproductDetails, nil
}

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
func UpdateProductImage(id int, file *multipart.FileHeader) error {

	url, err := helper.AddImageToS3(file)
	if err != nil {
		fmt.Println("error in s3", err)
		return err
	}
	err = repository.UpdateProductImage(id, url)
	if err != nil {
		fmt.Println("error in updation", err)
		return err
	}
	return nil
}
