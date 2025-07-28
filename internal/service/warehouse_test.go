package service

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UpdateWarehouseById(t *testing.T) {
	t.Run("should update warehouse successfully", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "W001",
			Address:            "Street 123",
			Telephone:          "1234567890",
			MinimumCapacity:    10,
			MinimumTemperature: &[]float64{10.0}[0],
		}

		mockRepository := new(repository.MockWarehouseRepository)
		mockRepository.On("GetWarehouseById", 1).Return(warehouse, nil)
		mockRepository.On("UpdateWarehouseById", 1, warehouse).Return(warehouse, nil)

		service := NewWarehouseDefault(mockRepository)

		result, err := service.UpdateWarehouseById(1, warehouse)

		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should return error if warehouse does not exist", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "W001",
			Address:            "Street 123",
			Telephone:          "1234567890",
			MinimumCapacity:    10,
			MinimumTemperature: &[]float64{10.0}[0],
		}

		mockRepository := new(repository.MockWarehouseRepository)
		mockRepository.On("GetWarehouseById", 1).Return(models.Warehouse{}, custom_errors.ErrNotFound)

		service := NewWarehouseDefault(mockRepository)

		result, err := service.UpdateWarehouseById(1, warehouse)

		require.Error(t, err)
		require.IsType(t, &custom_errors.ResourceNotFoundError{}, err)
		require.Empty(t, result)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should return error if minimum capacity is less than 0", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "W001",
			Address:            "Street 123",
			Telephone:          "1234567890",
			MinimumCapacity:    -1,
			MinimumTemperature: &[]float64{10.0}[0],
		}

		mockRepository := new(repository.MockWarehouseRepository)
		mockRepository.On("GetWarehouseById", 1).Return(warehouse, nil)

		service := NewWarehouseDefault(mockRepository)

		_, err := service.UpdateWarehouseById(1, warehouse)

		require.Error(t, err)
		require.IsType(t, &custom_errors.InvalidArgValueErr{}, err)
		mockRepository.AssertExpectations(t)
	})
}
