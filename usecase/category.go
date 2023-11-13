package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func GetCategory() ([]domain.Category, error) {
	category, err := repository.GetCategory()
	if err != nil {
		return []domain.Category{}, err
	}
	return category, nil

}
func AddCategory(category models.Category) (domain.Category, error) {
	exists, err := repository.CheckIfCategoryAlreadyExists(category.Category)
	if err != nil {
		return domain.Category{}, err
	}

	if exists {
		return domain.Category{}, errors.New("category already exists")
	}
	categories, err := repository.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categories, nil
}
func UpdateCategory(current string, new string) (domain.Category, error) {
	categries, err := repository.CheckCategory(current)
	if err != nil {
		return domain.Category{}, err
	}
	if !categries {
		return domain.Category{}, errors.New("category doesn't exist")
	}
	newCate, err := repository.UpdateCategory(current, new)
	if err != nil {
		return domain.Category{}, err
	}
	return newCate, nil
}
func DeleteCategory(id int) error {
	err := repository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
