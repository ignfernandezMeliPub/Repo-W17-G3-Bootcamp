package models_test

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/require"
)

var minimumTemperature = float64(10.2)

func TestWarehouse_Verify(t *testing.T) {
	t.Run("should pass verification with all valid fields", func(t *testing.T) {
		temp := 10.0
		w := models.Warehouse{
			WarehouseCode:      "WH-123",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &temp,
		}

		err := w.Verify()
		require.NoError(t, err)
	})
	t.Run("should return error if warehouse_code is empty", func(t *testing.T) {
		w := models.Warehouse{
			WarehouseCode:      "",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		}

		err := w.Verify()
		require.Error(t, err)

		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.Contains(t, err.Error(), "warehouse_code")
	})

	t.Run("should return error if address is empty", func(t *testing.T) {
		w := models.Warehouse{
			WarehouseCode:      "WH-123",
			Address:            "",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		}
		err := w.Verify()
		require.Error(t, err)

		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.Contains(t, err.Error(), "address")
	})
	t.Run("should return error if telephone is empty", func(t *testing.T) {
		w := models.Warehouse{
			WarehouseCode:      "WH-123",
			Address:            "Street 1",
			Telephone:          "",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		}
		err := w.Verify()
		require.Error(t, err)

		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.Contains(t, err.Error(), "telephone")
	})
	t.Run("should return error if minimum_capacity is zero", func(t *testing.T) {
		w := models.Warehouse{
			WarehouseCode:      "WH-123",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    0,
			MinimumTemperature: &minimumTemperature,
		}
		err := w.Verify()
		require.Error(t, err)

		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.Contains(t, err.Error(), "minimum_capacity")
	})

	t.Run("should return error if minimum temperature is nil", func(t *testing.T) {
		w := models.Warehouse{
			WarehouseCode:      "WH-123",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: nil,
		}
		err := w.Verify()
		require.Error(t, err)

		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.Contains(t, err.Error(), "minimum_temperature")
	})
}
