package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetAllBuyers(t *testing.T) {
	t.Run("should return all buyers successfully", func(t *testing.T) {
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

		expectedResponse := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345",
					"first_name": "John",
					"last_name": "Doe"
				},
				{
					"id": 2,
					"card_number_id": "67890",
					"first_name": "Jane",
					"last_name": "Smith"
				}
			]
		}`

		mockService := new(service.MockBuyerService)
		mockService.On("GetAllBuyers").Return(expectedBuyers, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers", nil)
		w := httptest.NewRecorder()

		handler.GetAllBuyers(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
		}`

		mockService := new(service.MockBuyerService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetAllBuyers").Return([]models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers", nil)
		w := httptest.NewRecorder()

		handler.GetAllBuyers(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})
}

func TestGetBuyerById(t *testing.T) {
	t.Run("should return the buyer successfully", func(t *testing.T) {
		expectedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "John",
			LastName:     "Doe",
		}

		expectedResponse := `{
			"data": {
				"id": 1,
				"card_number_id": "12345",
				"first_name": "John",
				"last_name": "Doe"
			}
		}`

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyerById", 1).Return(expectedBuyer, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		handler.GetBuyerById(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
		}`

		mockService := new(service.MockBuyerService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetBuyerById", 1).Return(models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		handler.GetBuyerById(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Failed to decode url param {id}. Please verify format."
		}`
		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		handler.GetBuyerById(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "GetBuyerById", 0)
	})

}

func TestCreateBuyer(t *testing.T) {
	t.Run("should create the buyer successfully", func(t *testing.T) {
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

		expectedResponse := `{
			"data": {
				"id": 1,
				"card_number_id": "1001",
				"first_name": "Jhon",
				"last_name": "Doe"
			}
		}`

		mockService := new(service.MockBuyerService)
		mockService.
			On("CreateBuyer", inputBuyer).
			Return(returnedBuyer, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Invalid body format."
		}`
		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "CreateBuyer", 0)
	})

	t.Run("should handle bad request error for body typo", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Invalid body format."
		}`
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

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "CreateBuyer", 0)
	})

	t.Run("should handle conflict error for unique constraint violation", func(t *testing.T) {
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

		serviceError := &custom_errors.UniqueAttributeViolationErr{
			AttributeName: "card_number_id",
			Value:         "1001",
		}

		expectedResponse := `{
			"message": "Conflict",
			"error": "Invalid value {1001} for unique attribute {card_number_id}. Value already being used."
		}`

		mockService := new(service.MockBuyerService)
		mockService.
			On("CreateBuyer", inputBuyer).
			Return(models.Buyer{}, serviceError)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle unprocessable content error for missing required fields", func(t *testing.T) {
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon"
		}`

		expectedResponse := `{
			"message": "Unprocessable Entity",
			"error": "Argument/s {last_name} is/are mandatory"
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "CreateBuyer", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		body := `{
			"card_number_id": "1001",
			"first_name": "Jhon",
			"last_name": ""
		}`

		expectedResponse := `{
			"message": "Unprocessable Entity",
			"error": "Invalid Value {} for arg {last_name}. Value must be non-empty."
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/buyers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateBuyer(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "CreateBuyer", 0)
	})

}

func TestPatchBuyer(t *testing.T) {
	t.Run("should patch the buyer successfully", func(t *testing.T) {
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

		expectedResponse := `{
			"data": {
				"id": 1,
				"card_number_id": "1001",
				"first_name": "New name",
				"last_name": "New lastname"
			}
		}`

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

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for invalid url param", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Failed to decode url param {id}. Please verify format."
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Invalid body format."
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		body := `{
			"first_name": "New name",
  	 		"last_name": "New lastname"
		}`

		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
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

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for body typo", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Invalid body format."
		}`

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

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle conflict error for unique constraint violation", func(t *testing.T) {
		body := `{
			"card_number_id": "1002"
		}`

		cardNumberId := "1002"
		inputBuyerPatch := models.BuyerPatch{
			CardNumberId: &cardNumberId,
		}

		serviceError := &custom_errors.UniqueAttributeViolationErr{
			AttributeName: "card_number_id",
			Value:         "1002",
		}

		expectedResponse := `{
			"message": "Conflict",
			"error": "Invalid value {1002} for unique attribute {card_number_id}. Value already being used."
		}`

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

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle unprocessable content error for missing required fields (patching nothing)", func(t *testing.T) {
		body := `{
		}`

		expectedResponse := `{
			"message": "Unprocessable Entity",
			"error": "Argument/s {card_number_id or first_name or last_name} is/are mandatory"
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "UpdateBuyerById", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		body := `{
			"first_name": ""
		}`

		expectedResponse := `{
			"message": "Unprocessable Entity",
			"error": "Invalid Value {} for arg {first_name}. Value must be non-empty."
		}`

		mockService := new(service.MockBuyerService)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/buyers/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.PatchBuyer(w, req)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "UpdateBuyerById", 0)
	})

}

func TestDeleteBuyer(t *testing.T) {
	t.Run("should delete and return no content", func(t *testing.T) {
		mockService := new(service.MockBuyerService)
		mockService.On("DeleteBuyerById", 1).Return(nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		handler.DeleteBuyer(w, req)

		require.Equal(t, http.StatusNoContent, w.Code)
		require.Empty(t, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
		}`

		mockService := new(service.MockBuyerService)
		mockService.On("DeleteBuyerById", 1).Return(custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		handler.DeleteBuyer(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Failed to decode url param {id}. Please verify format."
			}`

		mockService := new(service.MockBuyerService)
		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/buyers/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		handler.DeleteBuyer(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "DeleteBuyerById", 0)
	})

}

func TestGetBuyersPurchaseOrdersCount(t *testing.T) {
	t.Run("should return all buyers purchase orders count successfully", func(t *testing.T) {
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

		expectedResponse := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345",
					"first_name": "John",
					"last_name": "Doe",
					"purchase_orders_count": 10
				},
				{
					"id": 2,
					"card_number_id": "67890",
					"first_name": "Jane",
					"last_name": "Smith",
					"purchase_orders_count": 20
				}
			]
		}`

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return(expectedPurchaseOrdersCount, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders", nil)

		w := httptest.NewRecorder()

		handler.GetBuyersPurchaseOrdersCount(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error for buyers purchase orders count", func(t *testing.T) {
		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
		}`

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders", nil)

		w := httptest.NewRecorder()

		handler.GetBuyersPurchaseOrdersCount(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should return buyer purchase orders count by id successfully", func(t *testing.T) {
		expectedPurchaseOrdersCount := []models.BuyerPurchaseOrdersCount{
			{
				Id:                  1,
				CardNumberId:        "12345",
				FirstName:           "John",
				LastName:            "Doe",
				PurchaseOrdersCount: 10,
			},
		}
		expectedResponse := `{
			"data": [
				{
					"id": 1,
					"card_number_id": "12345",
					"first_name": "John",
					"last_name": "Doe",
					"purchase_orders_count": 10
				}
			]
		}`

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return(expectedPurchaseOrdersCount, nil)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=1", nil)

		w := httptest.NewRecorder()

		handler.GetBuyersPurchaseOrdersCount(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error for buyer purchase orders count", func(t *testing.T) {
		expectedResponse := `{
			"message": "Not found",
			"error": "Resource not found."
		}`

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=1", nil)

		w := httptest.NewRecorder()

		handler.GetBuyersPurchaseOrdersCount(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for buyer purchase orders count", func(t *testing.T) {
		expectedResponse := `{
			"message": "Bad request",
			"error": "Failed to decode query param {id}. Please verify format."
		}`

		buyerId := 1

		mockService := new(service.MockBuyerService)
		mockService.On("GetBuyersPurchaseOrdersCount", &buyerId).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.QueryParamDecodeErrorI)

		handler := NewBuyerDefault(mockService)

		req := httptest.NewRequest(http.MethodGet, "/buyers/reportPurchaseOrders?id=a", nil)

		w := httptest.NewRecorder()

		handler.GetBuyersPurchaseOrdersCount(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		require.JSONEq(t, expectedResponse, w.Body.String())

		mockService.AssertNumberOfCalls(t, "GetBuyersPurchaseOrdersCount", 0)
	})
}
