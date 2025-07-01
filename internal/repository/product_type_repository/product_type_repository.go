package product_type_repository

import "app/pkg/models"

type ProductTypeRepository interface {
	FindProductTypeById(id int) (models.ProductType, error)
}
