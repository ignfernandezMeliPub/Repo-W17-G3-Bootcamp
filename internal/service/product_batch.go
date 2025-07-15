package service

import (
	"app/internal/repository/product_batch_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type ProductBatchService interface {
	CreateProductBatch(productBatch models.ProductBatchRequest) (models.ProductBatch, error)
	GetAllProductBatchesBySection() ([]models.ProductBatchResponse, error)
	GetProductBatchBySectionId(sectionId int) (models.ProductBatchResponse, error)
}

type ProductBatchServiceImpl struct {
	rp product_batch_repository.ProductBatchRepository
}

func NewProductBatchService(rp product_batch_repository.ProductBatchRepository) ProductBatchService {
	return &ProductBatchServiceImpl{rp}
}

func (p ProductBatchServiceImpl) CreateProductBatch(pb models.ProductBatchRequest) (prod models.ProductBatch, err error) {
	mandatoryFields := map[string]any{
		"batch_number":        pb.BatchNumber,
		"current_quantity":    pb.CurrentQuantity,
		"current_temperature": pb.CurrentTemperature,
		"due_date":            pb.DueDate,
		"initial_quantity":    pb.InitialQuantity,
		"manufacturing_date":  pb.ManufacturingDate,
		"manufacturing_hour":  pb.ManufacturingHour,
		"minimum_temperature": pb.MinimumTemperature,
		"product_id":          pb.ProductId,
		"section_id":          pb.SectionId,
	}
	for field, value := range mandatoryFields {
		if value == nil {
			return prod, &custom_errors.MandatoryArgMissingErr{Argument: field}
		}
	}
	newProdBatch := models.ProductBatch{
		BatchNumber:        *pb.BatchNumber,
		CurrentQuantity:    *pb.CurrentQuantity,
		CurrentTemperature: *pb.CurrentTemperature,
		DueDate:            *pb.DueDate,
		InitialQuantity:    *pb.InitialQuantity,
		ManufacturingDate:  *pb.ManufacturingDate,
		ManufacturingHour:  *pb.ManufacturingHour,
		MinimumTemperature: *pb.MinimumTemperature,
		ProductId:          *pb.ProductId,
		SectionId:          *pb.SectionId,
	}
	return p.rp.CreateProductBatch(newProdBatch)
}

func (p ProductBatchServiceImpl) GetAllProductBatchesBySection() (prods []models.ProductBatchResponse, err error) {
	return p.rp.GetAllProductBatchesBySection()
}

func (p ProductBatchServiceImpl) GetProductBatchBySectionId(sectionId int) (prods models.ProductBatchResponse, err error) {
	return p.rp.GetProductBatchBySectionId(sectionId)
}
