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
		BatchNumber:        *pb.Data.BatchNumber,
		CurrentQuantity:    *pb.Data.CurrentQuantity,
		CurrentTemperature: *pb.Data.CurrentTemperature,
		DueDate:            *pb.Data.DueDate,
		InitialQuantity:    *pb.Data.InitialQuantity,
		ManufacturingDate:  *pb.Data.ManufacturingDate,
		ManufacturingHour:  *pb.Data.ManufacturingHour,
		MinimumTemperature: *pb.Data.MinimumTemperature,
		ProductId:          *pb.Data.ProductId,
		SectionId:          *pb.Data.SectionId,
	}
	return p.rp.CreateProductBatch(newProdBatch)
}
