package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockWarehouseService struct {
	mock.Mock
}

func (m *MockWarehouseService) GetAllWarehouses() ([]models.Warehouse, error) {
	args := m.Called()
	return args.Get(0).([]models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) GetWarehouseById(id int) (models.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) CreateWarehouse(vh models.Warehouse) (models.Warehouse, error) {
	args := m.Called(vh)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) UpdateWarehouseById(id int, w models.Warehouse) (models.Warehouse, error) {
	args := m.Called(id, w)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseService) DeleteWarehouse(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
