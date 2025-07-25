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
	"time"

	"github.com/stretchr/testify/require"
)

var orderDate, _ = time.Parse("2006-01-02", "2024-01-15")

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("should create the purchase order successfully with single detail", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		inputPurchaseOrder := models.PurchaseOrder{
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			BuyerId:      1,
			OrderDate:    orderDate,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 5},
			},
		}

		returnedPurchaseOrder := models.PurchaseOrder{
			Id:           1,
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			BuyerId:      1,
			OrderDate:    orderDate,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{Id: 1, ProductRecordId: 1, Quantity: 5},
			},
		}

		mockService := new(service.MockPurchaseOrderService)
		mockService.
			On("CreatePurchaseOrder", inputPurchaseOrder).
			Return(returnedPurchaseOrder, nil)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

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
		require.Equal(t, "ORD-12345", data["order_number"])
		require.Equal(t, "TRK-67890", data["tracking_code"])
		require.Equal(t, float64(1), data["buyer_id"])

		mockService.AssertExpectations(t)
	})

	t.Run("should create the purchase order successfully with multiple details", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					},
					{
						"product_record_id": 2,
						"quantity": 3
					}
				]
			}
		}`

		inputPurchaseOrder := models.PurchaseOrder{
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			BuyerId:      1,
			OrderDate:    orderDate,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 5},
				{ProductRecordId: 2, Quantity: 3},
			},
		}

		returnedPurchaseOrder := models.PurchaseOrder{
			Id:           1,
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			BuyerId:      1,
			OrderDate:    orderDate,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{Id: 1, ProductRecordId: 1, Quantity: 5},
				{Id: 2, ProductRecordId: 2, Quantity: 3},
			},
		}

		mockService := new(service.MockPurchaseOrderService)
		mockService.
			On("CreatePurchaseOrder", inputPurchaseOrder).
			Return(returnedPurchaseOrder, nil)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

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

		details, ok := data["purchase_order_details"].([]any)
		require.True(t, ok)
		require.Len(t, details, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle bad request error for invalid JSON body", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				],
			}
		}`

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle not found error when buyer_id does not exist", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 999,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		inputPurchaseOrder := models.PurchaseOrder{
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			OrderDate:    orderDate,
			BuyerId:      999,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 5},
			},
		}

		serviceError := custom_errors.ErrNotFound

		mockService := new(service.MockPurchaseOrderService)
		mockService.
			On("CreatePurchaseOrder", inputPurchaseOrder).
			Return(models.PurchaseOrder{}, serviceError)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

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

	t.Run("should handle conflict error when order_number already exists", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		inputPurchaseOrder := models.PurchaseOrder{
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			BuyerId:      1,
			OrderDate:    orderDate,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 5},
			},
		}

		serviceError := custom_errors.ErrUniqueAttributeViolationError

		mockService := new(service.MockPurchaseOrderService)
		mockService.
			On("CreatePurchaseOrder", inputPurchaseOrder).
			Return(models.PurchaseOrder{}, serviceError)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

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

	t.Run("should handle not found error when product_record_id does not exist", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 999,
						"quantity": 5
					}
				]
			}
		}`

		inputPurchaseOrder := models.PurchaseOrder{
			OrderNumber:  "ORD-12345",
			TrackingCode: "TRK-67890",
			OrderDate:    orderDate,
			BuyerId:      1,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 999, Quantity: 5},
			},
		}

		serviceError := custom_errors.ErrForeignKeyViolation

		mockService := new(service.MockPurchaseOrderService)
		mockService.
			On("CreatePurchaseOrder", inputPurchaseOrder).
			Return(models.PurchaseOrder{}, serviceError)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

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

	t.Run("should handle unprocessable content error for empty purchase_order_details", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": []
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle unprocessable content error for missing required fields", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle unprocessable content error for invalid field values", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle unprocessable content error for negative quantity", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": -5
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle bad request error for invalid quantity type", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": "not-a-number"
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle bad request error for invalid product_record_id type", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "2024-01-15",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": "not-a-number",
						"quantity": 5
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})

	t.Run("should handle unprocessable content error for invalid date format", func(t *testing.T) {
		// Arrange
		body := `{
			"data": {
				"order_number": "ORD-12345",
				"order_date": "invalid-date",
				"tracking_code": "TRK-67890",
				"buyer_id": 1,
				"purchase_order_details": [
					{
						"product_record_id": 1,
						"quantity": 5
					}
				]
			}
		}`

		mockService := new(service.MockPurchaseOrderService)

		handler := NewPurchaseOrderDefault(mockService)

		req := httptest.NewRequest(http.MethodPost, "/purchase-orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreatePurchaseOrder(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Contains(t, response, "message")
		require.Contains(t, response, "error")

		mockService.AssertNumberOfCalls(t, "CreatePurchaseOrder", 0)
	})
}
