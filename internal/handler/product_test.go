package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAllProducts(t *testing.T) {
	expectedResponse := `{"data": [{"id": 1, "product_code": "1234567890", "description": "Product 1", "product_type_id": 1, "width": 10, "height": 20, "length": 30, "net_weight": 100, "expiration_rate": 1, "recommended_freezing_temperature": 10, "freezing_rate": 1, "seller_id": null}, {"id": 2, "product_code": "1234567891", "description": "Product 2", "product_type_id": 2, "width": 10, "height": 20, "length": 30, "net_weight": 100, "expiration_rate": 1, "recommended_freezing_temperature": 10, "freezing_rate": 1, "seller_id": null}]}`
	t.Run("should return all products successfully", func(t *testing.T) {

		// Arrange
		expectedProducts := []models.Product{
			{
				ID:                             1,
				ProductCode:                    "1234567890",
				Description:                    "Product 1",
				Width:                          10.0,
				Height:                         20.0,
				Length:                         30.0,
				NetWeight:                      100.0,
				ExpirationRate:                 1,
				RecommendedFreezingTemperature: 10.0,
				FreezingRate:                   1,
				ProductTypeId:                  1,
				SellerId:                       nil,
			},
			{
				ID:                             2,
				ProductCode:                    "1234567891",
				Description:                    "Product 2",
				Width:                          10.0,
				Height:                         20.0,
				Length:                         30.0,
				NetWeight:                      100.0,
				ExpirationRate:                 1,
				RecommendedFreezingTemperature: 10.0,
				FreezingRate:                   1,
				ProductTypeId:                  2,
				SellerId:                       nil,
			},
		}
		mockService := new(service.MockProductService)
		mockService.On("GetAllProducts").Return(expectedProducts, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetAllProducts(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("should return error when products is empty", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetAllProducts").Return([]models.Product{}, serviceError)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetAllProducts(w, req)

		// Assert
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")
		mockService.AssertExpectations(t)
	})
}
