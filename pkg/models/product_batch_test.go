package models

import (
	"app/pkg/custom_errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func StringPtr(v string) *string { return &v }

func TestProductBatchRequest_Verify(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := ProductBatchRequest{
			Data: &ProductBatchData{
				BatchNumber:        IntPtr(1),
				CurrentQuantity:    IntPtr(1),
				CurrentTemperature: IntPtr(1),
				DueDate:            StringPtr("2020-09-20"),
				InitialQuantity:    IntPtr(1),
				ManufacturingDate:  StringPtr("2020-09-20"),
				ManufacturingHour:  IntPtr(1),
				MinimumTemperature: IntPtr(1),
				ProductId:          IntPtr(1),
				SectionId:          IntPtr(1),
			},
		}
		err := req.Verify()
		require.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		req := ProductBatchRequest{
			Data: &ProductBatchData{
				BatchNumber:        IntPtr(1),
				CurrentQuantity:    IntPtr(1),
				CurrentTemperature: nil,
				DueDate:            StringPtr("2020-09-20"),
				InitialQuantity:    IntPtr(1),
				ManufacturingDate:  StringPtr("2020-09-20"),
				ManufacturingHour:  IntPtr(1),
				MinimumTemperature: IntPtr(1),
				ProductId:          IntPtr(1),
				SectionId:          IntPtr(1),
			},
		}
		err := req.Verify()
		exp := &custom_errors.MandatoryArgMissingErr{Argument: "current_temperature"}
		require.Error(t, err)
		require.Equal(t, err, exp)
	})

	t.Run("invalid_due_date", func(t *testing.T) {
		req := ProductBatchRequest{
			Data: &ProductBatchData{
				BatchNumber:        IntPtr(1),
				CurrentQuantity:    IntPtr(1),
				CurrentTemperature: IntPtr(1),
				DueDate:            StringPtr("202a"),
				InitialQuantity:    IntPtr(1),
				ManufacturingDate:  StringPtr("2020-09-20"),
				ManufacturingHour:  IntPtr(1),
				MinimumTemperature: IntPtr(1),
				ProductId:          IntPtr(1),
				SectionId:          IntPtr(1),
			},
		}
		err := req.Verify()
		exp := &custom_errors.InvalidArgValueErr{Argument: "due_date", Value: "202a", ExtraInfo: "Invalid date format."}
		require.Error(t, err)
		require.Equal(t, err, exp)
	})

	t.Run("invalid_manufacturing_date", func(t *testing.T) {
		req := ProductBatchRequest{
			Data: &ProductBatchData{
				BatchNumber:        IntPtr(1),
				CurrentQuantity:    IntPtr(1),
				CurrentTemperature: IntPtr(1),
				DueDate:            StringPtr("2020-09-20"),
				InitialQuantity:    IntPtr(1),
				ManufacturingDate:  StringPtr("202a"),
				ManufacturingHour:  IntPtr(1),
				MinimumTemperature: IntPtr(1),
				ProductId:          IntPtr(1),
				SectionId:          IntPtr(1),
			},
		}
		err := req.Verify()
		exp := &custom_errors.InvalidArgValueErr{Argument: "manufacturing_date", Value: "202a", ExtraInfo: "Invalid date format."}
		require.Error(t, err)
		require.Equal(t, err, exp)
	})
}
