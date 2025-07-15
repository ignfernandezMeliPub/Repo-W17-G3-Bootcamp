package product_batch_repository

import "app/pkg/models"

type ProductBatchRepository interface {
	CreateProductBatch(productBatch models.ProductBatch) (models.ProductBatch, error)
}
