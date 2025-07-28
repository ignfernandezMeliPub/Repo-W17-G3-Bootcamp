package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProductRecordHandler(t *testing.T) {
	t.Run("should return a new product record handler", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)
		require.NotNil(t, handler)
	})
}

func TestGetAllProductRecords(t *testing.T) {
	expectedProductRecords := []models.ProductRecord{
		{
			ID:             1,
			LastUpdateDate: "2021-01-01",
			PurchasePrice:  100,
			SalePrice:      150,
			ProductID:      1,
		},
		{
			ID:             2,
			LastUpdateDate: "2021-01-02",
			PurchasePrice:  200,
			SalePrice:      250,
			ProductID:      2,
		},
		{
			ID:             3,
			LastUpdateDate: "2021-01-03",
			PurchasePrice:  300,
			SalePrice:      350,
			ProductID:      3,
		},
	}

	t.Run("should return all product records successfully", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		mockService.On("GetAllProductRecords").Return(expectedProductRecords, nil)

		request := httptest.NewRequest(http.MethodGet, "/product-records", nil)
		response := httptest.NewRecorder()
		handler.GetAllProductRecords(response, request)

		var res struct {
			Data []models.ProductRecord `json:"data"`
		}

		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, expectedProductRecords, res.Data)
		mockService.AssertExpectations(t)
	})
	t.Run("should return not found if there are no product records", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		mockService.On("GetAllProductRecords").Return([]models.ProductRecord{}, custom_errors.ErrNotFound)

		request := httptest.NewRequest(http.MethodGet, "/product-records", nil)
		response := httptest.NewRecorder()
		handler.GetAllProductRecords(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateProductRecord(t *testing.T) {
	body := `{
		"data": {
			"last_update_date": "2021-01-01",
			"purchase_price": 100,
			"sale_price": 150,
			"product_id": 1
		}
	}`
	bodyInvalid := `{
		"data": {
			"last_update_date": "2021-01-01",
			"purchase_price": 100,
			"sale_price": 150,
			"product_id": 1,
		}
	}`
	expectedProductRecord := models.ProductRecord{
		ID:             1,
		LastUpdateDate: "2021-01-01",
		PurchasePrice:  100,
		SalePrice:      150,
		ProductID:      1,
	}
	inputProductRecord := models.ProductRecordRequest{
		Data: &models.ProductRecordData{
			LastUpdateDate: &expectedProductRecord.LastUpdateDate,
			PurchasePrice:  &expectedProductRecord.PurchasePrice,
			SalePrice:      &expectedProductRecord.SalePrice,
			ProductID:      &expectedProductRecord.ProductID,
		},
	}

	t.Run("should create a product record successfully", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		mockService.On("CreateProductRecord", inputProductRecord).Return(expectedProductRecord, nil)

		request := httptest.NewRequest(http.MethodPost, "/productRecords", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.CreateProductRecord(response, request)

		var res struct {
			Data models.ProductRecord `json:"data"`
		}

		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, expectedProductRecord, res.Data)
		mockService.AssertExpectations(t)
	})
	t.Run("should return bad request if the request is invalid", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/productRecords", bytes.NewBufferString(bodyInvalid))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.CreateProductRecord(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		mockService.AssertNumberOfCalls(t, "CreateProductRecord", 0)
	})

	t.Run("should return an error if the service returns an error", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		mockService.On("CreateProductRecord", inputProductRecord).Return(models.ProductRecord{}, &custom_errors.InvalidArgValueErr{})

		request := httptest.NewRequest(http.MethodPost, "/productRecords", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.CreateProductRecord(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return an error if the date is in the future", func(t *testing.T) {
		mockService := new(service.MockProductRecordService)
		handler := NewProductRecordHandler(mockService)

		mockService.On("CreateProductRecord", inputProductRecord).Return(models.ProductRecord{}, &custom_errors.InvalidArgValueErr{
			Argument:  "last_update_date",
			Value:     "2021-01-01",
			ExtraInfo: "Date cannot be in the future",
		})

		request := httptest.NewRequest(http.MethodPost, "/productRecords", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.CreateProductRecord(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Contains(t, response.Body.String(), "Date cannot be in the future")
		mockService.AssertExpectations(t)
	})
}
