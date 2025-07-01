package service

import (
	"app/internal/repository/product_type_repository"
	"app/pkg/models"
)

type ProductTypeServices interface {
	GetProductTypeById(int) (models.ProductType, error)
}

type ProductTypeServiceImpl struct {
	ProductTypeRepository product_type_repository.ProductTypeRepository
}

func NewProductTypeService(rp product_type_repository.ProductTypeRepository) ProductTypeServiceImpl {
	return ProductTypeServiceImpl{ProductTypeRepository: rp}
}

func (s *ProductTypeServiceImpl) GetProductTypeById(id int) (models.ProductType, error) {
	return s.ProductTypeRepository.FindProductTypeById(id)
}
