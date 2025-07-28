package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(p models.Product) (models.Product, error) {
	args := m.Called(p)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductById(id int) (models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductByCode(code string) (models.Product, error) {
	args := m.Called(code)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProductById(product models.Product) (models.Product, error) {
	args := m.Called(product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) DeleteProductById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) GetReportRecords(id *int) ([]models.ProductRecordReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.ProductRecordReport), args.Error(1)
}
