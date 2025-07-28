package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) CreateWarehouse(wh models.Warehouse) (models.Warehouse, error) {
	args := m.Called(wh)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) GetAllWarehouses() ([]models.Warehouse, error) {
	args := m.Called()
	return args.Get(0).([]models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) GetWarehouseById(id int) (models.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) UpdateWarehouseById(id int, w models.Warehouse) (models.Warehouse, error) {
	args := m.Called(id, w)
	return args.Get(0).(models.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) DeleteWarehouseById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
