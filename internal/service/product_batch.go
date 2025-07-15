package service

import (
	"app/internal/repository/product_batch_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"time"
)

type ProductBatchService interface {
	CreateProductBatch(productBatch models.ProductBatchRequest) (models.ProductBatch, error)
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
	dateLayout := "2006-01-02"
	if _, err := time.Parse(dateLayout, *pb.DueDate); err != nil {
		return prod, &custom_errors.InvalidArgValueErr{Argument: "due_date", Value: *pb.DueDate, ExtraInfo: "Incorrect date"}
	}
	if _, err := time.Parse(dateLayout, *pb.ManufacturingDate); err != nil {
		return prod, &custom_errors.InvalidArgValueErr{Argument: "manufacturing_date", Value: *pb.ManufacturingDate, ExtraInfo: "Incorrect date"}
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
