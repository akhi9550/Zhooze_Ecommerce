package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func AddCategory(category domain.Category) (domain.Category, error) {
	categories, err := repository.AddCategory(category)
	if err != nil {
		return domain.Category{}, err
	}
	return categories, nil
}
func DeleteCategory(id string) error {
	err := repository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
func UpdateCategory(current string, new string) (models.UpdateCategory, error) {
	categries, err := repository.CheckCategory(current)
	if err != nil {
		return models.UpdateCategory{}, err
	}
	if !categries {
		return models.UpdateCategory{}, errors.New("category doesn't exist")
	}
	newcat, err := repository.UpdateCategory(current, new)
	if err != nil {
		return models.UpdateCategory{}, err
	}
	return newcat, nil
}
