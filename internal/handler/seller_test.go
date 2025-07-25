package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// CREATE TESTS
func TestCreateSeller(t *testing.T) {
	t.Run("create_ok - should create seller successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"cid": 1001,
			"company_name": "Test Company",
			"address": "123 Test St",
			"telephone": "555-1234",
			"locality_id": "LOC001"
		}`

		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mockService := new(service.MockSellerService)
		mockService.On("CreateSeller", 1001, "Test Company", "123 Test St", "555-1234", "LOC001").Return(expectedSeller, nil)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/sellers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateSeller(w, req)

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

		seller := data[0].(map[string]any)
		require.Equal(t, float64(1), seller["id"])
		require.Equal(t, float64(1001), seller["cid"])
		require.Equal(t, "Test Company", seller["company_name"])
		require.Equal(t, "123 Test St", seller["address"])
		require.Equal(t, "555-1234", seller["telephone"])
		require.Equal(t, "LOC001", seller["locality_id"])

		mockService.AssertExpectations(t)
	})

	t.Run("create_bad_request - should handle bad request for malformed JSON", func(t *testing.T) {
		// Arrange
		body := `{
			"cid": 1001,
			"company_name": "Test Company",
			"address": "123 Test St",
			"telephone": "555-1234",
			"locality_id": "LOC001",
		}`

		mockService := new(service.MockSellerService)
		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/sellers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateSeller(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateSeller")
	})

	t.Run("create_fail - should handle unprocessable entity for missing required fields", func(t *testing.T) {
		// Arrange
		body := `{
			"cid": 1001,
			"company_name": "Test Company",
			"address": "123 Test St"
		}`

		mockService := new(service.MockSellerService)
		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/sellers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateSeller(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateSeller")
	})

	t.Run("create_conflict - should handle conflict when cid already exists", func(t *testing.T) {
		// Arrange
		body := `{
			"cid": 1001,
			"company_name": "Test Company",
			"address": "123 Test St",
			"telephone": "555-1234",
			"locality_id": "LOC001"
		}`

		mockService := new(service.MockSellerService)
		mockService.On("CreateSeller", 1001, "Test Company", "123 Test St", "555-1234", "LOC001").Return(models.Seller{}, custom_errors.ErrUniqueAttributeViolationError)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/sellers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateSeller(w, req)

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
}

// READ TESTS
func TestGetAllSellers(t *testing.T) {
	t.Run("find_all - should return all sellers successfully", func(t *testing.T) {
		// Arrange
		expectedSellers := []models.Seller{
			{
				Id:          1,
				CompanyId:   1001,
				CompanyName: "Company A",
				Address:     "123 Test St",
				Telephone:   "555-1111",
				LocalityId:  "LOC001",
			},
			{
				Id:          2,
				CompanyId:   1002,
				CompanyName: "Company B",
				Address:     "456 Test Ave",
				Telephone:   "555-2222",
				LocalityId:  "LOC002",
			},
		}

		mockService := new(service.MockSellerService)
		mockService.On("GetAllSellers").Return(expectedSellers, nil)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/sellers", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetAllSellers(w, req)

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

		// Verify first seller
		firstSeller := data[0].(map[string]any)
		require.Equal(t, float64(1), firstSeller["id"])
		require.Equal(t, float64(1001), firstSeller["cid"])
		require.Equal(t, "Company A", firstSeller["company_name"])

		mockService.AssertExpectations(t)
	})
}

func TestGetSellerById(t *testing.T) {
	t.Run("find_by_id_existent - should return seller successfully", func(t *testing.T) {
		// Arrange
		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mockService := new(service.MockSellerService)
		mockService.On("GetSellerById", 1).Return(expectedSeller, nil)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/sellers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.GetSellerById(w, req)

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

		seller := data[0].(map[string]any)
		require.Equal(t, float64(1), seller["id"])
		require.Equal(t, float64(1001), seller["cid"])
		require.Equal(t, "Test Company", seller["company_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("find_by_id_non_existent - should handle not found error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockSellerService)
		mockService.On("GetSellerById", 999).Return(models.Seller{}, custom_errors.ErrNotFound)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/sellers/999", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.GetSellerById(w, req)

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

// UPDATE TESTS
func TestPatchSeller(t *testing.T) {
	t.Run("update_ok - should update seller successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"company_name": "Updated Company",
			"address": "456 Updated St"
		}`

		companyName := "Updated Company"
		address := "456 Updated St"

		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Updated Company",
			Address:     "456 Updated St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mockService := new(service.MockSellerService)
		mockService.On("UpdateSellerById", 1, (*int)(nil), &companyName, &address, (*string)(nil)).Return(expectedSeller, nil)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/sellers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchSeller(w, req)

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

		seller := data[0].(map[string]any)
		require.Equal(t, "Updated Company", seller["company_name"])
		require.Equal(t, "456 Updated St", seller["address"])

		mockService.AssertExpectations(t)
	})

	t.Run("update_non_existent - should handle not found error", func(t *testing.T) {
		// Arrange
		body := `{
			"company_name": "Updated Company"
		}`

		companyName := "Updated Company"

		mockService := new(service.MockSellerService)
		mockService.On("UpdateSellerById", 999, (*int)(nil), &companyName, (*string)(nil), (*string)(nil)).Return(models.Seller{}, custom_errors.ErrNotFound)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/sellers/999", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchSeller(w, req)

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

// DELETE TESTS
func TestDeleteSeller(t *testing.T) {
	t.Run("delete_ok - should delete seller successfully", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockSellerService)
		mockService.On("DeleteSellerById", 1).Return(nil)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/sellers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteSeller(w, req)

		// Assert
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Empty(t, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("delete_non_existent - should handle not found error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockSellerService)
		mockService.On("DeleteSellerById", 999).Return(custom_errors.ErrNotFound)

		handler := NewSellerHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/sellers/999", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteSeller(w, req)

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
