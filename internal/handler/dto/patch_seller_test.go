package dto

import (
	"app/pkg/custom_errors"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPatchSellerDto_Verify(t *testing.T) {
	t.Run("valid dto with all fields", func(t *testing.T) {
		dto := PatchSellerDto{
			CompanyId:   intPtr(123),
			CompanyName: stringPtr("Test Company"),
			Address:     stringPtr("123 Test Street"),
			Telephone:   stringPtr("555-1234"),
		}

		err := dto.Verify()
		assert.NoError(t, err)
	})

	t.Run("valid dto with only company_id", func(t *testing.T) {
		dto := PatchSellerDto{
			CompanyId:   intPtr(123),
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
		}

		err := dto.Verify()
		assert.NoError(t, err)
	})

	t.Run("all fields nil", func(t *testing.T) {
		dto := PatchSellerDto{
			CompanyId:   nil,
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
		}

		err := dto.Verify()
		assert.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		assert.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		assert.Equal(t, "cid or company_name or address or telephone", mandatoryErr.Argument)
	})
}
