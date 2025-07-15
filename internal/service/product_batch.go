package service

import (
	"app/internal/repository/repositories/product_batch_repository"
	"app/pkg/models"
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
