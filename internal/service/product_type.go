package service

import (
	"app/internal/repository/repositories/product_type_repository"
	"app/pkg/models"
)

type ProductTypeServices interface {
	GetProductTypeById(int) (models.ProductType, error)
	IsValidProductType(int) bool
}

type ProductTypeServiceImpl struct {
	ProductTypeRepository product_type_repository.ProductTypeRepository
}

func NewProductTypeService(rp product_type_repository.ProductTypeRepository) *ProductTypeServiceImpl {
	return &ProductTypeServiceImpl{ProductTypeRepository: rp}
}

func (s *ProductTypeServiceImpl) GetProductTypeById(id int) (models.ProductType, error) {
	return s.ProductTypeRepository.GetProductTypeById(id)
}

func (s *ProductTypeServiceImpl) IsValidProductType(id int) bool {
	return s.ProductTypeRepository.IsValidProductType(id)
}
