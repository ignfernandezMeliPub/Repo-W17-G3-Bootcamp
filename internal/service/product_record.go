package service

import (
	"app/internal/repository/repositories/product_record_repository"
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
	parsedDate, err := time.Parse("2006-01-02", *productRecord.Data.LastUpdateDate)
	if err != nil {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "last_update_date", Value: *productRecord.Data.LastUpdateDate, ExtraInfo: "Date format must be YYYY-MM-DD"}
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.After(today) {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{
			Argument:  "last_update_date",
			Value:     *productRecord.Data.LastUpdateDate,
			ExtraInfo: "Date cannot be in the future",
		}
	}

	if *productRecord.Data.PurchasePrice <= 0.0 {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "purchase_price", Value: *productRecord.Data.PurchasePrice, ExtraInfo: "Purchase price must be greater than 0"}
	}

	if *productRecord.Data.SalePrice <= 0.0 {
		return models.ProductRecord{}, &custom_errors.InvalidArgValueErr{Argument: "sale_price", Value: *productRecord.Data.SalePrice, ExtraInfo: "Sale price must be greater than 0"}
	}

	productRecordModel := models.ProductRecord{
		LastUpdateDate: *productRecord.Data.LastUpdateDate,
		PurchasePrice:  *productRecord.Data.PurchasePrice,
		SalePrice:      *productRecord.Data.SalePrice,
		ProductID:      *productRecord.Data.ProductID,
	}
	return s.productRecordRepository.CreateProductRecord(productRecordModel)
}
