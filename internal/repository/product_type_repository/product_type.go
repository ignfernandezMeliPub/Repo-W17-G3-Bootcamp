package product_type_repository

import "app/pkg/models"

type ProductTypeRepository interface {
	GetProductTypeById(id int) (models.ProductType, error)
	IsValidProductType(id int) bool
}
