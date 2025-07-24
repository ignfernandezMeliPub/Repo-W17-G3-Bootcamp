package service

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSectionsServiceImpl_UpdateSectionById(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSc := NewSectionsService(dbMock)

	dbMock.On("GetSectionById", 1).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 2, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	dbMock.On("GetSectionById", 2).Return(models.Section{}, &custom_errors.ResourceNotFoundError{})

	dbMock.On("UpdateSectionById", models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1.5, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1.5, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	t.Run("success", func(t *testing.T) {
		currentTemperature := float32(1.5)
		req := models.SectionRequest{
			CurrentTemperature: &currentTemperature,
		}
		sec, err := svSc.UpdateSectionById(1, req)
		exp := models.Section{
			ID: 1, SectionNumber: 1, CurrentTemperature: 1.5, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
		}

		require.NoError(t, err)
		require.Equal(t, exp, sec)
	})

	t.Run("fail", func(t *testing.T) {
		currentTemperature := float32(1.5)
		req := models.SectionRequest{
			CurrentTemperature: &currentTemperature,
		}
		sec, err := svSc.UpdateSectionById(2, req)
		errExp := &custom_errors.ResourceNotFoundError{}

		require.NotNil(t, err)
		require.ErrorIs(t, err, errExp)
		require.Equal(t, models.Section{}, sec)
	})
}
