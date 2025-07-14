package service

import (
	"app/internal/repository/product_record_repository"
	"app/pkg/models"
)

type ProductRecordService interface {
	GetAllProductRecords() ([]models.ProductRecord, error)
	CreateProductRecord(productRecord models.ProductRecordRequest) (models.ProductRecord, error)
}

type ProductRecordServiceImpl struct {
	productRecordRepository product_record_repository.ProductRecordRepository
}

func NewProductRecordService(productRecordRepository product_record_repository.ProductRecordRepository) *ProductRecordServiceImpl {
	return &ProductRecordServiceImpl{productRecordRepository: productRecordRepository}
}

func (s *ProductRecordServiceImpl) GetAllProductRecords() ([]models.ProductRecord, error) {
	return s.productRecordRepository.GetAllProductRecords()
}

func (s *ProductRecordServiceImpl) CreateProductRecord(productRecord models.ProductRecordRequest) (models.ProductRecord, error) {
	productRecordModel := models.ProductRecord{
		LastUpdateDate: *productRecord.LastUpdateDate,
		PurchasePrice:  *productRecord.PurchasePrice,
		SalePrice:      *productRecord.SalePrice,
		ProductID:      *productRecord.ProductID,
	}
	return s.productRecordRepository.CreateProductRecord(productRecordModel)
}
