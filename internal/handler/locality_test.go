package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestLocalityHandler_GetCarriesReport(t *testing.T) {
	t.Run("should return report by locality id successfully", func(t *testing.T) {
		mockService := new(service.MockLocalityService)

		expected := []models.CarriesReport{
			{LocalityId: "L001", LocalityName: "Localidad 1", CarriesCount: 3},
		}
		mockService.
			On("GetCarriesReport", "L001").
			Return(expected, nil)

		hd := NewLocalityHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/localities/reportCarries?id=L001", nil)
		response := httptest.NewRecorder()

		hd.GetCarriesReport(response, request)

		res := response.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

		body := response.Body.String()
		assert.Contains(t, body, "L001")
		assert.Contains(t, body, "Localidad 1")
		assert.Contains(t, body, "3")

		mockService.AssertExpectations(t)
	})
	t.Run("should return all reports with empty id", func(t *testing.T) {
		mockService := new(service.MockLocalityService)

		expected := []models.CarriesReport{
			{LocalityId: "L001", LocalityName: "Locality 1", CarriesCount: 3},
			{LocalityId: "L002", LocalityName: "Locality 2", CarriesCount: 5},
		}
		mockService.
			On("GetCarriesReport", "").
			Return(expected, nil)

		hd := NewLocalityHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/localities/reportCarries?id=", nil)
		response := httptest.NewRecorder()

		hd.GetCarriesReport(response, request)

		res := response.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

		body := response.Body.String()
		assert.Contains(t, body, "L001")
		assert.Contains(t, body, "Locality 1")
		assert.Contains(t, body, "3")
		assert.Contains(t, body, "L002")
		assert.Contains(t, body, "Locality 2")
		assert.Contains(t, body, "5")

		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 when locality not found", func(t *testing.T) {
		mockService := new(service.MockLocalityService)

		mockService.
			On("GetCarriesReport", "NONEXISTENT").
			Return([]models.CarriesReport{}, custom_errors.ErrNotFound)

		hd := NewLocalityHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/localities/reportCarries?id=NONEXISTENT", nil)
		response := httptest.NewRecorder()

		hd.GetCarriesReport(response, request)

		res := response.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

		body := response.Body.String()
		assert.Contains(t, body, "Not found")

		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 when database error occurs", func(t *testing.T) {
		mockService := new(service.MockLocalityService)

		mockService.
			On("GetCarriesReport", "L001").
			Return([]models.CarriesReport{}, sql.ErrConnDone)

		hd := NewLocalityHandler(mockService)

		request := httptest.NewRequest(http.MethodGet, "/localities/reportCarries?id=L001", nil)
		response := httptest.NewRecorder()

		hd.GetCarriesReport(response, request)

		res := response.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

		body := response.Body.String()
		assert.Contains(t, body, "Internal server error")

		mockService.AssertExpectations(t)
	})

}
