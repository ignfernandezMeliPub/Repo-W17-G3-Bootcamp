package service

import (
	"app/internal/repository/product_record_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"time"
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

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", *productRecord.LastUpdateDate)
	if err != nil {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "last_update_date", Value: *productRecord.LastUpdateDate, ExtraInfo: "Date format must be YYYY-MM-DD"}
	}

	if *productRecord.PurchasePrice <= 0.0 {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "purchase_price", Value: *productRecord.PurchasePrice, ExtraInfo: "Purchase price must be greater than 0"}
	}

	if *productRecord.SalePrice <= 0.0 {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "sale_price", Value: *productRecord.SalePrice, ExtraInfo: "Sale price must be greater than 0"}
	}

	productRecordModel := models.ProductRecord{
		LastUpdateDate: *productRecord.LastUpdateDate, //poner time.now()?
		PurchasePrice:  *productRecord.PurchasePrice,
		SalePrice:      *productRecord.SalePrice,
		ProductID:      *productRecord.ProductID,
	}
	return s.productRecordRepository.CreateProductRecord(productRecordModel)
}
