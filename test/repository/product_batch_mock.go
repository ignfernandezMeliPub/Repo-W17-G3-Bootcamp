package repository

import (
	"app/pkg/models"
	"github.com/stretchr/testify/mock"
)

type ProductBatchMock struct {
	mock.Mock
}

func NewProductBatchMock() *ProductBatchMock {
	return &ProductBatchMock{}
}

func (p *ProductBatchMock) CreateProductBatch(productBatch models.ProductBatch) (models.ProductBatch, error) {
	args := p.Called(productBatch)
	return args.Get(0).(models.ProductBatch), args.Error(1)
}
