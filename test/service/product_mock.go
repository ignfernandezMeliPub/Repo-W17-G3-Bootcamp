package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductService) GetProductById(id int) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductService) CreateProduct(product models.ProductRequest) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProductById(product models.ProductPatchRequest) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductService) DeleteProductById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductService) GetReportRecords(id *int) ([]models.ProductRecordReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.ProductRecordReport), args.Error(1)
}
