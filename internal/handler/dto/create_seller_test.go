package dto

import (
	"app/pkg/custom_errors"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func TestCreateSellerDto_Verify(t *testing.T) {
	t.Run("valid dto", func(t *testing.T) {
		dto := CreateSellerDto{
			CompanyId:   intPtr(123),
			CompanyName: stringPtr("Test Company"),
			Address:     stringPtr("123 Test Street"),
			Telephone:   stringPtr("555-1234"),
			LocalityId:  stringPtr("locality-123"),
		}

		err := dto.Verify()
		assert.NoError(t, err)
	})

	t.Run("nil company_id", func(t *testing.T) {
		dto := CreateSellerDto{
			CompanyId:   nil,
			CompanyName: stringPtr("Test Company"),
			Address:     stringPtr("123 Test Street"),
			Telephone:   stringPtr("555-1234"),
			LocalityId:  stringPtr("locality-123"),
		}

		err := dto.Verify()
		assert.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		assert.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		assert.Equal(t, "cid", mandatoryErr.Argument)
	})

	t.Run("nil company_name", func(t *testing.T) {
		dto := CreateSellerDto{
			CompanyId:   intPtr(123),
			CompanyName: nil,
			Address:     stringPtr("123 Test Street"),
			Telephone:   stringPtr("555-1234"),
			LocalityId:  stringPtr("locality-123"),
		}

		err := dto.Verify()
		assert.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		assert.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		assert.Equal(t, "company_name", mandatoryErr.Argument)
	})

	t.Run("nil address", func(t *testing.T) {
		dto := CreateSellerDto{
			CompanyId:   intPtr(123),
			CompanyName: stringPtr("Test Company"),
			Address:     nil,
			Telephone:   stringPtr("555-1234"),
			LocalityId:  stringPtr("locality-123"),
		}

		err := dto.Verify()
		assert.Error(t, err)

		var mandatoryErr *custom_errors.MandatoryArgMissingErr
		ok := errors.As(err, &mandatoryErr)
		assert.True(t, ok, "Expected *custom_errors.MandatoryArgMissingErr")
		assert.Equal(t, "address", mandatoryErr.Argument)
	})

	t.Run("empty strings and zero values are valid", func(t *testing.T) {
		dto := CreateSellerDto{
			CompanyId:   intPtr(0),
			CompanyName: stringPtr(""),
			Address:     stringPtr(""),
			Telephone:   stringPtr(""),
			LocalityId:  stringPtr(""),
		}

		err := dto.Verify()
		assert.NoError(t, err)
	})
}
