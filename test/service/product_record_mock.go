package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRecordService struct {
	mock.Mock
}

func (m *MockProductRecordService) GetAllProductRecords() ([]models.ProductRecord, error) {
	args := m.Called()
	return args.Get(0).([]models.ProductRecord), args.Error(1)
}

func (m *MockProductRecordService) CreateProductRecord(productRecord models.ProductRecordRequest) (models.ProductRecord, error) {
	args := m.Called(productRecord)
	return args.Get(0).(models.ProductRecord), args.Error(1)
}
