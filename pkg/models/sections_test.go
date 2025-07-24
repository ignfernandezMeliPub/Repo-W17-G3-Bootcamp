package models

import (
	"app/pkg/custom_errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func IntPtr(v int) *int             { return &v }
func Float32Ptr(v float32) *float32 { return &v }

func TestSectionRequest_Verify(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := SectionRequest{
			SectionNumber:      IntPtr(1),
			CurrentTemperature: Float32Ptr(3.5),
			MinimumTemperature: Float32Ptr(2.0),
			CurrentCapacity:    IntPtr(100),
			MinimumCapacity:    IntPtr(10),
			MaximumCapacity:    IntPtr(200),
			WarehouseId:        IntPtr(5),
			ProductTypeId:      IntPtr(7),
		}
		err := req.Verify()
		require.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		req := SectionRequest{
			SectionNumber:      IntPtr(1),
			CurrentTemperature: nil,
			MinimumTemperature: Float32Ptr(2.0),
			CurrentCapacity:    IntPtr(100),
			MinimumCapacity:    IntPtr(10),
			MaximumCapacity:    IntPtr(200),
			WarehouseId:        IntPtr(5),
			ProductTypeId:      IntPtr(7),
		}
		err := req.Verify()
		exp := &custom_errors.MandatoryArgMissingErr{Argument: "current_temperature"}
		require.Error(t, err)
		require.Equal(t, err, exp)
	})
}
