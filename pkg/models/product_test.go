package models

import (
	"app/pkg/custom_errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var sellerId = 1

var product = Product{
	ID:                             1,
	ProductCode:                    "1234567890",
	Description:                    "Product Description",
	Width:                          10.0,
	Height:                         20.0,
	Length:                         30.0,
	NetWeight:                      100.0,
	ExpirationRate:                 10,
	RecommendedFreezingTemperature: 10.0,
	FreezingRate:                   10,
	ProductTypeId:                  1,
	SellerId:                       &sellerId,
}

func TestProductRequestVerify(t *testing.T) {
	mandatoryFields := []string{
		"ProductCode", "Description", "Width", "Height", "Length",
		"NetWeight", "ExpirationRate", "RecommendedFreezingTemperature", "FreezingRate", "ProductTypeId",
	}

	for _, fieldName := range mandatoryFields {
		t.Run("should return error if "+fieldName+" is nil", func(t *testing.T) {
			validRequest := createValidProductRequest()
			setFieldValue(&validRequest, fieldName, nil)

			err := validRequest.Verify()
			require.Error(t, err)
			require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		})
	}

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		validRequest := createValidProductRequest()
		err := validRequest.Verify()
		require.NoError(t, err)
	})
}

func TestProductPatchRequestVerify(t *testing.T) {
	t.Run("should return nil for patch request", func(t *testing.T) {
		patchRequest := ProductPatchRequest{
			Id: 1,
		}

		err := patchRequest.Verify()
		require.NoError(t, err)
	})
}

// Helper function to create a valid ProductRequest
func createValidProductRequest() ProductRequest {
	return ProductRequest{
		ProductCode:                    &product.ProductCode,
		Description:                    &product.Description,
		Width:                          &product.Width,
		Height:                         &product.Height,
		Length:                         &product.Length,
		NetWeight:                      &product.NetWeight,
		ExpirationRate:                 &product.ExpirationRate,
		RecommendedFreezingTemperature: &product.RecommendedFreezingTemperature,
		FreezingRate:                   &product.FreezingRate,
		ProductTypeId:                  &product.ProductTypeId,
		SellerId:                       product.SellerId,
	}
}

// Helper function to set field value
func setFieldValue(request *ProductRequest, fieldName string, value interface{}) {
	switch fieldName {
	case "ProductCode":
		if value == nil {
			request.ProductCode = nil
		} else {
			v := value.(*string)
			request.ProductCode = v
		}
	case "Description":
		if value == nil {
			request.Description = nil
		} else {
			v := value.(*string)
			request.Description = v
		}
	case "Width":
		if value == nil {
			request.Width = nil
		} else {
			v := value.(*float64)
			request.Width = v
		}
	case "Height":
		if value == nil {
			request.Height = nil
		} else {
			v := value.(*float64)
			request.Height = v
		}
	case "Length":
		if value == nil {
			request.Length = nil
		} else {
			v := value.(*float64)
			request.Length = v
		}
	case "NetWeight":
		if value == nil {
			request.NetWeight = nil
		} else {
			v := value.(*float64)
			request.NetWeight = v
		}
	case "ExpirationRate":
		if value == nil {
			request.ExpirationRate = nil
		} else {
			v := value.(*int)
			request.ExpirationRate = v
		}
	case "RecommendedFreezingTemperature":
		if value == nil {
			request.RecommendedFreezingTemperature = nil
		} else {
			v := value.(*float64)
			request.RecommendedFreezingTemperature = v
		}
	case "FreezingRate":
		if value == nil {
			request.FreezingRate = nil
		} else {
			v := value.(*int)
			request.FreezingRate = v
		}
	case "ProductTypeId":
		if value == nil {
			request.ProductTypeId = nil
		} else {
			v := value.(*int)
			request.ProductTypeId = v
		}
	case "SellerId":
		if value == nil {
			request.SellerId = nil
		} else {
			v := value.(*int)
			request.SellerId = v
		}
	}
}
