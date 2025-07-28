package models

import (
	"app/pkg/custom_errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var productRecord = ProductRecord{
	ID:             1,
	LastUpdateDate: "2021-01-01",
	PurchasePrice:  100,
	SalePrice:      150,
	ProductID:      1,
}

func TestProductRecordVerify(t *testing.T) {
	mandatoryFields := []string{
		"LastUpdateDate", "PurchasePrice", "SalePrice", "ProductID",
	}

	for _, fieldName := range mandatoryFields {
		t.Run("should return error if "+fieldName+" is nil", func(t *testing.T) {
			validRequest := createValidProductRecordRequest()
			setFieldValueProductRecord(&validRequest, fieldName, nil)

			err := validRequest.Verify()
			require.Error(t, err)
			require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		})
	}

	t.Run("should return error if data is nil", func(t *testing.T) {
		validRequest := createValidProductRecordRequest()
		validRequest.Data = nil
		err := validRequest.Verify()
		require.Error(t, err)
		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
	})
	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		validRequest := createValidProductRecordRequest()
		err := validRequest.Verify()
		require.NoError(t, err)
	})
}

func createValidProductRecordRequest() ProductRecordRequest {
	return ProductRecordRequest{
		Data: &ProductRecordData{
			LastUpdateDate: &productRecord.LastUpdateDate,
			PurchasePrice:  &productRecord.PurchasePrice,
			SalePrice:      &productRecord.SalePrice,
			ProductID:      &productRecord.ProductID,
		},
	}
}

func setFieldValueProductRecord(request *ProductRecordRequest, fieldName string, value interface{}) {
	switch fieldName {
	case "LastUpdateDate":
		if value == nil {
			request.Data.LastUpdateDate = nil
		} else {
			v := value.(*string)
			request.Data.LastUpdateDate = v
		}
	case "PurchasePrice":
		if value == nil {
			request.Data.PurchasePrice = nil
		} else {
			v := value.(*float64)
			request.Data.PurchasePrice = v
		}
	case "SalePrice":
		if value == nil {
			request.Data.SalePrice = nil
		} else {
			v := value.(*float64)
			request.Data.SalePrice = v
		}
	case "ProductID":
		if value == nil {
			request.Data.ProductID = nil
		} else {
			v := value.(*int)
			request.Data.ProductID = v
		}
	}
}
