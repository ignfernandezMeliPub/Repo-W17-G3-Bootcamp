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

func TestGetAllBuyers(t *testing.T) {
	t.Run("should return all buyers successfully", func(t *testing.T) {
		// Arrange
		expectedBuyers := []models.Buyer{
			{
				Id:           1,
				CardNumberId: "12345",
				FirstName:    "John",
				LastName:     "Doe",
			},
			{
				Id:           2,
				CardNumberId: "67890",
				FirstName:    "Jane",
				LastName:     "Smith",
			},
		}

		mockService := new(service.MockBuyerService)
		mockService.On("GetAllBuyers").Return(expectedBuyers, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetAllBuyers(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")

		// Verificar que la data contenga los buyers esperados
		data, ok := response["data"].([]any)
		require.True(t, ok)
		require.Len(t, data, 2)

		// Verificar el primer buyer
		firstBuyer := data[0].(map[string]any)
		require.Equal(t, float64(1), firstBuyer["id"])
		require.Equal(t, "12345", firstBuyer["card_number_id"])
		require.Equal(t, "John", firstBuyer["first_name"])
		require.Equal(t, "Doe", firstBuyer["last_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetAllBuyers").Return([]models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetAllBuyers(w, req)

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

func TestGetBuyerById(t *testing.T) {
	t.Run("should return the buyer successfully", func(t *testing.T) {
		// Arrange
		expectedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "John",
			LastName:     "Doe",
		}

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyerById", 1).Return(expectedBuyer, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyerById(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")

		// Verificar que la data contenga los buyers esperados
		data, ok := response["data"].(map[string]any)
		require.True(t, ok)

		// Verificar el primer buyer
		buyer := data
		require.Equal(t, float64(1), buyer["id"])
		require.Equal(t, "12345", buyer["card_number_id"])
		require.Equal(t, "John", buyer["first_name"])
		require.Equal(t, "Doe", buyer["last_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetBuyerById", 1).Return(models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.GetBuyerById(w, req)

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

	t.Run("should handle bad request error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.GetBuyerById(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "GetBuyerById", 0)
	})

}

func TestCreateBuyer(t *testing.T) {
	t.Run("should create the buyer successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": "Doe"
		}`

		inputBuyer := models.Buyer{
			CardNumberId: "1001",
			FirstName:    "Jhon",
			LastName:     "Doe",
		}

		returnedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "1001",
			FirstName:    "Jhon",
			LastName:     "Doe",
		}

		mockService := new(service.MockBuyerService)
		mockService.
			On("CreateBuyer", inputBuyer).
			Return(returnedBuyer, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")
		data, ok := response["data"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, float64(1), data["id"])
		require.Equal(t, "1001", data["card_number_id"])
		require.Equal(t, "Jhon", data["first_name"])
		require.Equal(t, "Doe", data["last_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateBuyer", 0)
	})

	t.Run("should handle bad request error for body typo", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": "Doe",
		}`

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateBuyer", 0)
	})

	t.Run("should handle conflict error for unique constraint violation", func(t *testing.T) {
		// Arrange
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": "Doe"
		}`

		inputBuyer := models.Buyer{
			CardNumberId: "1001",
			FirstName:    "Jhon",
			LastName:     "Doe",
		}

		serviceError := custom_errors.ErrUniqueAttributeViolationError

		mockService := new(service.MockBuyerService)
		mockService.
			On("CreateBuyer", inputBuyer).
			Return(models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.CreateBuyer(w, req)

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

	t.Run("should handle unprocessable content error for missing required fields", func(t *testing.T) {
		// Arrange
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon"
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.CreateBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateBuyer", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		// Arrange
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": ""
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.CreateBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "CreateBuyer", 0)
	})

}

func TestPatchBuyer(t *testing.T) {
	t.Run("should patch the buyer successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"first_name": "New name",
  	 		"last_name": "New lastname"
		}`

		firstName := "New name"
		lastName := "New lastname"

		inputBuyerPatch := models.BuyerPatch{
			FirstName: &firstName,
			LastName:  &lastName,
		}

		returnedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "1001",
			FirstName:    "New name",
			LastName:     "New lastname",
		}

		mockService := new(service.MockBuyerService)
		mockService.
			On("UpdateBuyerById", 1, inputBuyerPatch).
			Return(returnedBuyer, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "data")
		data, ok := response["data"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "New name", data["first_name"])
		require.Equal(t, "New lastname", data["last_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for invalid url param", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange

		body := `{
			"first_name": "New name",
  	 		"last_name": "New lastname"
		}`

		firstName := "New name"
		lastName := "New lastname"

		mockService := new(service.MockBuyerService)
		mockService.
			On("UpdateBuyerById", 1, models.BuyerPatch{
				FirstName: &firstName,
				LastName:  &lastName,
			}).
			Return(models.Buyer{}, custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchBuyer(w, req)

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

	t.Run("should handle bad request error for body typo", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": "Doe",
		}`

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle conflict error for unique constraint violation", func(t *testing.T) {
		// Arrange
		body := `{
			"card_number_id": "1002"
		}`

		cardNumberId := "1002"
		inputBuyerPatch := models.BuyerPatch{
			CardNumberId: &cardNumberId,
		}

		serviceError := custom_errors.ErrUniqueAttributeViolationError

		mockService := new(service.MockBuyerService)
		mockService.
			On("UpdateBuyerById", 1, inputBuyerPatch).
			Return(models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.PatchBuyer(w, req)

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

	t.Run("should handle unprocessable content error for missing required fields (patching nothing)", func(t *testing.T) {
		// Arrange
		body := `{
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		// Arrange
		body := `{
			"first_name": ""
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		// Act
		handler.PatchBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "UpdateBuyerById", 0)
	})

}

func TestDeleteBuyer(t *testing.T) {
	t.Run("should delete and return no content", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)
		mockService.On("DeleteBuyerById", 1).Return(nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Empty(t, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)
		mockService.On("DeleteBuyerById", 1).Return(custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteBuyer(w, req)

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

	t.Run("should handle bad request error", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.DeleteBuyer(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "DeleteBuyerById", 0)
	})

}

func TestGetBuyersPurchaseOrdersCount(t *testing.T) {
	t.Run("should return all buyers purchase orders count successfully", func(t *testing.T) {
		// Arrange
		expectedPurchaseOrdersCount := []models.BuyerPurchaseOrdersCount{
			{
				Id:                  1,
				CardNumberId:        "12345",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 10,
			},
			{
				Id:                  2,
				CardNumberId:        "67890",
				FirstName:           "Jane",
				LastName:            "Smith",
				PurchaseOrdersCount: 20,
			},
		}

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return(expectedPurchaseOrdersCount, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyersPurchaseOrdersCount(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		type Result struct {
			Data []models.BuyerPurchaseOrdersCount `json:"data"`
		}

		var result Result
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		require.Equal(t, expectedPurchaseOrdersCount, result.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error for buyers purchase orders count", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyersPurchaseOrdersCount(w, req)

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

	t.Run("should return buyer purchase orders count by id successfully", func(t *testing.T) {
		// Arrange
		expectedPurchaseOrdersCount := []models.BuyerPurchaseOrdersCount{
			{
				Id:                  1,
				CardNumberId:        "12345",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 10,
			},
		}

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return(expectedPurchaseOrdersCount, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=1", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyersPurchaseOrdersCount(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		type Result struct {
			Data []models.BuyerPurchaseOrdersCount `json:"data"`
		}

		var result Result
		err := json.Unmarshal(w.Body.Bytes(), &result)
		require.NoError(t, err)

		require.Equal(t, expectedPurchaseOrdersCount, result.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error for buyer purchase orders count", func(t *testing.T) {
		// Arrange

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=1", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyersPurchaseOrdersCount(w, req)

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

	t.Run("should handle bad request error for buyer purchase orders count", func(t *testing.T) {
		// Arrange

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.QueryParamDecodeErrorI)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=a", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetBuyersPurchaseOrdersCount(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNotCalled(t, "GetBuyersPurchaseOrdersCount", 0)
	})
}
