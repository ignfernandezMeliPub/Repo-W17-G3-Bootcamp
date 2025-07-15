package product_record_repository

import (
	"app/pkg/models"
)

type ProductRecordRepository interface {
	CreateProductRecord(productRecord models.ProductRecord) (models.ProductRecord, error)
	GetAllProductRecords() ([]models.ProductRecord, error)
}
