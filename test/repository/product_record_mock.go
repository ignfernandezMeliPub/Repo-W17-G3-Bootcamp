package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRecordRepository struct {
	mock.Mock
}

func (m *MockProductRecordRepository) GetAllProductRecords() ([]models.ProductRecord, error) {
	args := m.Called()
	return args.Get(0).([]models.ProductRecord), args.Error(1)
}

func (m *MockProductRecordRepository) CreateProductRecord(productRecord models.ProductRecord) (models.ProductRecord, error) {
	args := m.Called(productRecord)
	return args.Get(0).(models.ProductRecord), args.Error(1)
}
