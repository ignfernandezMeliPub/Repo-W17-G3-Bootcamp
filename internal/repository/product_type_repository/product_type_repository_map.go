package product_type_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type ProductTypeRepositoryMap struct {
	database map[int]models.ProductType
}

func NewProductTypeRepositoryMap(database map[int]models.ProductType) *ProductTypeRepositoryMap {

	if database == nil {
		database = map[int]models.ProductType{}
	}

	return &ProductTypeRepositoryMap{database: database}

}

func (r *ProductTypeRepositoryMap) FindProductTypeById(id int) (models.ProductType, error) {

	productType, found := r.database[id]

	if !found {
		return productType, &custom_errors.ResourceNotFoundError{}
	}
	return productType, nil

}
