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

// CREATE LOCALITY TESTS

func TestCreateLocality(t *testing.T) {
	t.Run("should create locality successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"id": "LOC001",
				"locality_name": "Buenos Aires",
				"province_name": "Buenos Aires",
				"country_name": "Argentina"
			}
		}`

		expectedLocality := models.Locality{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		}

		mockService := new(service.MockLocalityService)
		mockService.On("CreateLocality", "LOC001", "Buenos Aires", "Buenos Aires", "Argentina").Return(expectedLocality, nil)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateLocality(w, req)

		// Assert
		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")
		data, ok := response["data"].([]any)
		require.True(t, ok)
		require.Len(t, data, 1)

		locality := data[0].(map[string]any)
		require.Equal(t, "LOC001", locality["id"])
		require.Equal(t, "Buenos Aires", locality["locality_name"])
		require.Equal(t, "Buenos Aires", locality["province_name"])
		require.Equal(t, "Argentina", locality["country_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request for malformed JSON", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"id": "LOC001",
				"locality_name": "Buenos Aires",
				"province_name": "Buenos Aires",
				"country_name": "Argentina",
			}
		}`

		mockService := new(service.MockLocalityService)
		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateLocality(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateLocality")
	})

	t.Run("should handle unprocessable entity for missing required fields", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"id": "LOC001",
				"locality_name": "Buenos Aires"
			}
		}`

		mockService := new(service.MockLocalityService)
		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateLocality(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateLocality")
	})

	t.Run("should handle conflict when locality id already exists", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"id": "LOC001",
				"locality_name": "Buenos Aires",
				"province_name": "Buenos Aires",
				"country_name": "Argentina"
			}
		}`

		mockService := new(service.MockLocalityService)
		mockService.On("CreateLocality", "LOC001", "Buenos Aires", "Buenos Aires", "Argentina").Return(models.Locality{}, custom_errors.ErrUniqueAttributeViolationError)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateLocality(w, req)

		// Assert
		require.Equal(t, http.StatusConflict, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertExpectations(t)
	})

	t.Run("should handle missing data object", func(t *testing.T) {
		// Arrange
		body := `{}`

		mockService := new(service.MockLocalityService)
		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/localities", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateLocality(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateLocality")
	})
}

// GET LOCALITY SELLER COUNT TESTS

func TestGetLocalitySellerCount(t *testing.T) {
	t.Run("should return seller count for specific locality", func(t *testing.T) {
		// Arrange
		expectedResult := models.LocalitySellerCount{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			SellersCount: 5,
		}

		mockService := new(service.MockLocalityService)
		mockService.On("GetLocalitySellerCount", "LOC001").Return(expectedResult, nil)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=LOC001", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetLocalitySellerCount(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")
		data, ok := response["data"].([]any)
		require.True(t, ok)
		require.Len(t, data, 1)

		locality := data[0].(map[string]any)
		require.Equal(t, "LOC001", locality["id"])
		require.Equal(t, "Buenos Aires", locality["locality_name"])
		require.Equal(t, float64(5), locality["sellers_count"])

		mockService.AssertExpectations(t)
	})

	t.Run("should return seller count for all localities", func(t *testing.T) {
		// Arrange
		expectedResults := []models.LocalitySellerCount{
			{
				Id:           "LOC001",
				LocalityName: "Buenos Aires",
				SellersCount: 5,
			},
			{
				Id:           "LOC002",
				LocalityName: "Córdoba",
				SellersCount: 3,
			},
		}

		mockService := new(service.MockLocalityService)
		mockService.On("GetLocalitiesSellerCount").Return(expectedResults, nil)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/localities/reportSellers", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetLocalitySellerCount(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")
		data, ok := response["data"].([]any)
		require.True(t, ok)
		require.Len(t, data, 2)

		// Verify first locality
		firstLocality := data[0].(map[string]any)
		require.Equal(t, "LOC001", firstLocality["id"])
		require.Equal(t, "Buenos Aires", firstLocality["locality_name"])
		require.Equal(t, float64(5), firstLocality["sellers_count"])

		// Verify second locality
		secondLocality := data[1].(map[string]any)
		require.Equal(t, "LOC002", secondLocality["id"])
		require.Equal(t, "Córdoba", secondLocality["locality_name"])
		require.Equal(t, float64(3), secondLocality["sellers_count"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error for specific locality", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockLocalityService)
		mockService.On("GetLocalitySellerCount", "NONEXISTENT").Return(models.LocalitySellerCount{}, custom_errors.ErrNotFound)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/localities/reportSellers?id=NONEXISTENT", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetLocalitySellerCount(w, req)

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

	t.Run("should handle not found error for all localities", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockLocalityService)
		mockService.On("GetLocalitiesSellerCount").Return([]models.LocalitySellerCount{}, custom_errors.ErrNotFound)

		handler := NewLocalityHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/localities/reportSellers", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetLocalitySellerCount(w, req)

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
