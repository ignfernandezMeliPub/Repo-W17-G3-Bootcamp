package dto

import (
	"app/pkg/custom_errors"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func TestCreateLocalityDto_Verify(t *testing.T) {
	t.Run("valid dto", func(t *testing.T) {
		dto := CreateLocalityDto{
			Data: &CreateLocalityData{
				Id:           stringPtr("123"),
				LocalityName: stringPtr("Test Locality"),
				ProvinceName: stringPtr("Test Province"),
				CountryName:  stringPtr("Test Country"),
			},
		}

		err := dto.Verify()
		require.NoError(t, err)
	})

	t.Run("nil data", func(t *testing.T) {
		dto := CreateLocalityDto{Data: nil}

		err := dto.Verify()
		require.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		require.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		require.Equal(t, "data", mandatoryErr.Argument)
	})

	t.Run("nil id", func(t *testing.T) {
		dto := CreateLocalityDto{
			Data: &CreateLocalityData{
				Id:           nil,
				LocalityName: stringPtr("Test Locality"),
				ProvinceName: stringPtr("Test Province"),
				CountryName:  stringPtr("Test Country"),
			},
		}

		err := dto.Verify()
		require.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		require.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		require.Equal(t, "data.id", mandatoryErr.Argument)
	})

	t.Run("empty strings are valid", func(t *testing.T) {
		dto := CreateLocalityDto{
			Data: &CreateLocalityData{
				Id:           stringPtr(""),
				LocalityName: stringPtr(""),
				ProvinceName: stringPtr(""),
				CountryName:  stringPtr(""),
			},
		}

		err := dto.Verify()
		require.NoError(t, err)
	})
}
