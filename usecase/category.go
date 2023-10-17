package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"errors"
)

func AddCategory(category domain.Category) (domain.Category, error) {
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
func DeleteCategory(id string) error {
	err := repository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
