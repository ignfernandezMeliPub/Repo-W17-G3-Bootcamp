package service

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProductRecordService(t *testing.T) {
	mockRepository := new(repository.MockProductRecordRepository)
	service := NewProductRecordService(mockRepository)
	require.NotNil(t, service)
	mockRepository.AssertExpectations(t)
}

func TestProductRecordService_CreateProductRecord(t *testing.T) {

	productRecord := models.ProductRecord{
		ID:             0,
		LastUpdateDate: "2021-01-01",
		PurchasePrice:  100,
		SalePrice:      150,
		ProductID:      1,
	}
	productRecordOutput := models.ProductRecord{
		ID:             1,
		LastUpdateDate: "2021-01-01",
		PurchasePrice:  100,
		SalePrice:      150,
		ProductID:      1,
	}
	t.Run("should create a product record successfully", func(t *testing.T) {
		mockRepository := new(repository.MockProductRecordRepository)
		service := NewProductRecordService(mockRepository)
		productRecordRequest := models.ProductRecordRequest{
			Data: &models.ProductRecordData{
				LastUpdateDate: &productRecord.LastUpdateDate,
				PurchasePrice:  &productRecord.PurchasePrice,
				SalePrice:      &productRecord.SalePrice,
				ProductID:      &productRecord.ProductID,
			},
		}
		mockRepository.On("CreateProductRecord", productRecord).Return(productRecordOutput, nil)
		createdProductRecord, err := service.CreateProductRecord(productRecordRequest)
		require.NoError(t, err)
		require.Equal(t, productRecordOutput, createdProductRecord)
		mockRepository.AssertExpectations(t)
	})
	t.Run("should return error when lastUpdateDate is not a valid date", func(t *testing.T) {
		mockRepository := new(repository.MockProductRecordRepository)
		service := NewProductRecordService(mockRepository)

		lastUpdateDate := "2021-99-01"

		productRecordRequest := models.ProductRecordRequest{
			Data: &models.ProductRecordData{
				LastUpdateDate: &lastUpdateDate,
				PurchasePrice:  &productRecord.PurchasePrice,
				SalePrice:      &productRecord.SalePrice,
				ProductID:      &productRecord.ProductID,
			},
		}
		_, err := service.CreateProductRecord(productRecordRequest)
		require.Error(t, err)
		require.IsType(t, &custom_errors.InvalidArgValueErr{}, err)
		require.Equal(t, "last_update_date", err.(*custom_errors.InvalidArgValueErr).Argument)
		mockRepository.AssertNumberOfCalls(t, "CreateProductRecord", 0)
	})
	t.Run("should return error when purchasePrice is not a valid number", func(t *testing.T) {
		mockRepository := new(repository.MockProductRecordRepository)
		service := NewProductRecordService(mockRepository)

		purchasePrice := 0.0
		productRecordRequest := models.ProductRecordRequest{
			Data: &models.ProductRecordData{
				LastUpdateDate: &productRecord.LastUpdateDate,
				PurchasePrice:  &purchasePrice,
				SalePrice:      &productRecord.SalePrice,
				ProductID:      &productRecord.ProductID,
			},
		}
		_, err := service.CreateProductRecord(productRecordRequest)
		require.Error(t, err)
		require.IsType(t, &custom_errors.InvalidArgValueErr{}, err)
		require.Equal(t, "purchase_price", err.(*custom_errors.InvalidArgValueErr).Argument)
		mockRepository.AssertNumberOfCalls(t, "CreateProductRecord", 0)
	})
	t.Run("should return error when salePrice is not a valid number", func(t *testing.T) {
		mockRepository := new(repository.MockProductRecordRepository)
		service := NewProductRecordService(mockRepository)

		salePrice := 0.0
		productRecordRequest := models.ProductRecordRequest{
			Data: &models.ProductRecordData{
				LastUpdateDate: &productRecord.LastUpdateDate,
				PurchasePrice:  &productRecord.PurchasePrice,
				SalePrice:      &salePrice,
				ProductID:      &productRecord.ProductID,
			},
		}
		_, err := service.CreateProductRecord(productRecordRequest)
		require.Error(t, err)
		require.IsType(t, &custom_errors.InvalidArgValueErr{}, err)
		require.Equal(t, "sale_price", err.(*custom_errors.InvalidArgValueErr).Argument)
		mockRepository.AssertNumberOfCalls(t, "CreateProductRecord", 0)
	})

}
