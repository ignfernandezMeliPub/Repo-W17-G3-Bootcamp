package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateInboundOrder(t *testing.T) {

	t.Run("Case 1. create inbound order succesfully", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "77777",
				"employee_id":4,
				"product_batch_id": 1,
				"warehouse_id": 1
				}
			}`

		orderDate := "2021-04-04"
		orderNumber := "77777"
		employeeId := 4
		productBatchId := 1
		warehouseId := 1

		inputInboundOrder := models.InboundOrderRequestBody{
			Data: &models.InboundOrderData{

				OrderDate:      &orderDate,
				OrderNumber:    &orderNumber,
				EmployeeId:     &employeeId,
				ProductBatchId: &productBatchId,
				WarehouseId:    &warehouseId,
			},
		}

		returnedInboundOrder := models.InboundOrder{
			Id: 1,
			InboundOrderDetails: models.InboundOrderDetails{
				OrderDate:      "2021-04-04",
				OrderNumber:    "77777",
				EmployeeId:     4,
				ProductBatchId: 1,
				WarehouseId:    1,
			},
		}

		mockService := new(service.MockInboundOrderService)
		mockService.On("CreateInboundOrder", inputInboundOrder).Return(returnedInboundOrder, nil)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{ "data": { "id": 1, "order_date": "2021-04-04", "order_number": "77777", "employee_id": 4, "product_batch_id": 1, "warehouse_id": 1 } }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error null body - 400 Bad request", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockInboundOrderService)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"error": "Invalid body format.", "message": "Bad request"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. error unique order_number duplicated - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "ORD001",
				"employee_id":4,
				"product_batch_id": 1,
				"warehouse_id": 1
				}
			}`

		orderDate := "2021-04-04"
		orderNumber := "ORD001"
		employeeId := 4
		productBatchId := 1
		warehouseId := 1

		inputInboundOrder := models.InboundOrderRequestBody{
			Data: &models.InboundOrderData{

				OrderDate:      &orderDate,
				OrderNumber:    &orderNumber,
				EmployeeId:     &employeeId,
				ProductBatchId: &productBatchId,
				WarehouseId:    &warehouseId,
			},
		}

		returnedInboundOrder := models.InboundOrder{}
		returnedError := custom_errors.ErrUniqueAttributeViolationError
		returnedError.AttributeName = "order_number"
		returnedError.Value = "ORD001"

		mockService := new(service.MockInboundOrderService)
		mockService.
			On("CreateInboundOrder", inputInboundOrder).
			Return(returnedInboundOrder, returnedError)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Invalid value {ORD001} for unique attribute {order_number}. Value already being used.", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 4. error employee id not found - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "77777",
				"employee_id":99,
				"product_batch_id": 1,
				"warehouse_id": 1
				}
			}`

		orderDate := "2021-04-04"
		orderNumber := "77777"
		employeeId := 99
		productBatchId := 1
		warehouseId := 1

		inputInboundOrder := models.InboundOrderRequestBody{
			Data: &models.InboundOrderData{

				OrderDate:      &orderDate,
				OrderNumber:    &orderNumber,
				EmployeeId:     &employeeId,
				ProductBatchId: &productBatchId,
				WarehouseId:    &warehouseId,
			},
		}

		returnedInboundOrder := models.InboundOrder{}

		returnedError := custom_errors.ErrForeignKeyViolation
		returnedError.ConstraintName = "employee_id"

		mockService := new(service.MockInboundOrderService)
		mockService.
			On("CreateInboundOrder", inputInboundOrder).
			Return(returnedInboundOrder, returnedError)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Unknown entity identifier value: employee_id ", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 5. error product batch id not found - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "77777",
				"employee_id":1,
				"product_batch_id": 99,
				"warehouse_id": 1
				}
			}`

		orderDate := "2021-04-04"
		orderNumber := "77777"
		employeeId := 1
		productBatchId := 99
		warehouseId := 1

		inputInboundOrder := models.InboundOrderRequestBody{
			Data: &models.InboundOrderData{

				OrderDate:      &orderDate,
				OrderNumber:    &orderNumber,
				EmployeeId:     &employeeId,
				ProductBatchId: &productBatchId,
				WarehouseId:    &warehouseId,
			},
		}

		returnedInboundOrder := models.InboundOrder{}

		returnedError := custom_errors.ErrForeignKeyViolation
		returnedError.ConstraintName = "product_batch_id"

		mockService := new(service.MockInboundOrderService)
		mockService.
			On("CreateInboundOrder", inputInboundOrder).
			Return(returnedInboundOrder, returnedError)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Unknown entity identifier value: product_batch_id ", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 6. error warehouse id not found - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "77777",
				"employee_id": 1,
				"product_batch_id": 1,
				"warehouse_id": 99
				}
			}`

		orderDate := "2021-04-04"
		orderNumber := "77777"
		employeeId := 1
		productBatchId := 1
		warehouseId := 99

		inputInboundOrder := models.InboundOrderRequestBody{
			Data: &models.InboundOrderData{

				OrderDate:      &orderDate,
				OrderNumber:    &orderNumber,
				EmployeeId:     &employeeId,
				ProductBatchId: &productBatchId,
				WarehouseId:    &warehouseId,
			},
		}

		returnedInboundOrder := models.InboundOrder{}

		returnedError := custom_errors.ErrForeignKeyViolation
		returnedError.ConstraintName = "warehouse_id"

		mockService := new(service.MockInboundOrderService)
		mockService.
			On("CreateInboundOrder", inputInboundOrder).
			Return(returnedInboundOrder, returnedError)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Unknown entity identifier value: warehouse_id ", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 7. error empty body - 422 UnprocessableEntity", func(t *testing.T) {

		// Arrange
		body := `{}`

		mockService := new(service.MockInboundOrderService)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{"error":"Argument/s {data} is/are mandatory", "message":"Unprocessable Entity"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 8. missing field - 422 UnprocessableEntity", func(t *testing.T) {

		// Arrange
		body := `{
  			"data": {
				"order_date": "2021-04-04",
				"order_number": "77777",
				"product_batch_id": 1,
				"warehouse_id": 99
				}
			}`

		mockService := new(service.MockInboundOrderService)

		handler := NewInboundOrderController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/inboundOrders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateInboundOrder(w, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{"error":"Argument/s {data.employee_id} is/are mandatory", "message":"Unprocessable Entity"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

}
