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

func TestGetAllEmployees(t *testing.T) {

	t.Run("Case 1. success - return all employees succesfully", func(t *testing.T) {

		//Arrange

		employees := []models.Employee{
			{
				Id: 1,
				EmployeeAttributes: models.EmployeeAttributes{
					CardNumberId: "EMP001",
					FirstName:    "Raul",
					LastName:     "García",
					WarehouseId:  1,
				},
			},
			{
				Id: 2,
				EmployeeAttributes: models.EmployeeAttributes{
					CardNumberId: "EMP002",
					FirstName:    "Sandra",
					LastName:     "Rojas",
					WarehouseId:  2,
				},
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetAllEmployees").Return(employees, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees", nil)
		w := httptest.NewRecorder()

		//act

		handler.GetAllEmployees(w, req)

		//assert

		expectedCode := http.StatusOK
		expectedBody := `{"data": [{ "id": 1, "card_number_id": "EMP001", "first_name": "Raul", "last_name": "García", "warehouse_id": 1 }, { "id": 2, "card_number_id": "EMP002", "first_name": "Sandra", "last_name": "Rojas", "warehouse_id": 2 }]}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error - return not found when employee table is empty", func(t *testing.T) {

		//Arrange

		employees := []models.Employee{}
		serviceError := custom_errors.ErrNotFound

		mockService := new(service.MockEmployeeService)
		mockService.On("GetAllEmployees").Return(employees, serviceError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees", nil)
		w := httptest.NewRecorder()

		//act

		handler.GetAllEmployees(w, req)

		//assert

		expectedCode := http.StatusNotFound
		expectedBody := `{"error": "Resource not found.","message": "Not found"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

}

func TestGetEmployeeById(t *testing.T) {

	t.Run("Case 1. success - return employee requested", func(t *testing.T) {

		//Arrange

		employee := models.Employee{
			Id: 1,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberId: "EMP001",
				FirstName:    "Raul",
				LastName:     "García",
				WarehouseId:  1,
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetEmployeeById", 1).Return(employee, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		//act

		handler.GetEmployeeById(w, req)

		//assert

		expectedCode := http.StatusOK
		expectedBody := `{ "data": { "id": 1, "card_number_id": "EMP001", "first_name": "Raul", "last_name": "García", "warehouse_id": 1 } }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error - not found employee", func(t *testing.T) {

		//Arrange

		employee := models.Employee{}
		serviceError := custom_errors.ErrNotFound

		mockService := new(service.MockEmployeeService)
		mockService.On("GetEmployeeById", 1000).Return(employee, serviceError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/1000", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1000")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		//act

		handler.GetEmployeeById(w, req)

		//assert

		expectedCode := http.StatusNotFound
		expectedBody := `{"error": "Resource not found.","message": "Not found"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. error - Bad request - id has a wrong formart", func(t *testing.T) {

		//Arrange

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/1.0", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1.0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		//act

		handler.GetEmployeeById(w, req)

		//assert

		expectedCode := http.StatusBadRequest
		expectedBody := `{ "error": "Failed to decode url param {id}. Please verify format.", "message": "Bad request" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

	})

}

func TestCreateEmployee(t *testing.T) {

	t.Run("Case 1. create employee succesfully", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "123456",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 1
		}`

		cardNumberId := "123456"
		firstName := "Juan"
		lastName := "Pérez"
		warehouseId := 1

		inputEmployeeAttributes := models.EmployeePostRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{
			Id: 1,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberId: "123456",
				FirstName:    "Juan",
				LastName:     "Pérez",
				WarehouseId:  1,
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.
			On("CreateEmployee", inputEmployeeAttributes).
			Return(returnedEmployee, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusCreated
		expectedBody := `{ "data": { "id": 1, "card_number_id": "123456", "first_name": "Juan", "last_name": "Pérez", "warehouse_id": 1 } }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error null body - 400 Bad request", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"error": "Invalid body format.", "message": "Bad request"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error unique card id duplicated - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "123456",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 1
		}`

		cardNumberId := "123456"
		firstName := "Juan"
		lastName := "Pérez"
		warehouseId := 1

		inputEmployeeAttributes := models.EmployeePostRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{}
		returnedError := custom_errors.ErrUniqueAttributeViolationError
		returnedError.AttributeName = "card_number_id"
		returnedError.Value = 123456

		mockService := new(service.MockEmployeeService)
		mockService.
			On("CreateEmployee", inputEmployeeAttributes).
			Return(returnedEmployee, returnedError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Invalid value {123456} for unique attribute {card_number_id}. Value already being used.", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. error warehouse id not found - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "55555",
			"first_name": "Juan",
			"last_name": "Pérez",
			"warehouse_id": 99
		}`

		cardNumberId := "55555"
		firstName := "Juan"
		lastName := "Pérez"
		warehouseId := 99

		inputEmployeeAttributes := models.EmployeePostRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{}
		returnedError := custom_errors.ErrForeignKeyViolation
		returnedError.ConstraintName = "warehouse_id"

		mockService := new(service.MockEmployeeService)
		mockService.
			On("CreateEmployee", inputEmployeeAttributes).
			Return(returnedEmployee, returnedError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Unknown entity identifier value: warehouse_id ", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. error empty body - 422 UnprocessableEntity", func(t *testing.T) {

		// Arrange
		body := `{}`

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{"error":"Argument/s {card_number_id} is/are mandatory", "message":"Unprocessable Entity"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 4. missing field - 422 UnprocessableEntity", func(t *testing.T) {

		// Arrange
		body := `{"card_number_id": "55555",
			"first_name": "Juan",
			"warehouse_id": 99}`

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateEmployee(w, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{"error":"Argument/s {last_name} is/are mandatory", "message":"Unprocessable Entity"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

}

func TestPatchEmployee(t *testing.T) {

	t.Run("Case 1. Update employee succesfully - all attributes", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "99999",
			"first_name": "Juana",
			"last_name": "Peralta",
			"warehouse_id": 2
		}`

		cardNumberId := "99999"
		firstName := "Juana"
		lastName := "Peralta"
		warehouseId := 2

		inputEmployeeAttributes := models.EmployeePatchRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{
			Id: 1,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberId: "99999",
				FirstName:    "Juana",
				LastName:     "Peralta",
				WarehouseId:  2,
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.
			On("UpdateEmployeeById", 1, inputEmployeeAttributes).
			Return(returnedEmployee, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{ "data": { "id": 1, "card_number_id": "99999", "first_name": "Juana", "last_name": "Peralta", "warehouse_id": 2 } }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. wrong types in body - 400 Bad Request", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": 12345,
            "warehouse_id": 1
		}`

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"error": "Field {card_number_id} is expected to be of type {string}, but found {number}","message": "Bad request"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. url param with wrong format - 400 Bad Request", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "12345",
            "warehouse_id": 1
		}`

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1.0", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1.0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"error": "Failed to decode url param {id}. Please verify format.","message": "Bad request"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 4. body is null - 400 Bad Request", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{"error": "Invalid body format.", "message": "Bad request"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 5. error employee doesn't exists - 404 Not found", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "99999",
			"first_name": "Juana",
			"last_name": "Peralta",
			"warehouse_id": 2
		}`

		cardNumberId := "99999"
		firstName := "Juana"
		lastName := "Peralta"
		warehouseId := 2

		inputEmployeeAttributes := models.EmployeePatchRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{}
		returnedError := custom_errors.ErrNotFound

		mockService := new(service.MockEmployeeService)
		mockService.
			On("UpdateEmployeeById", 99, inputEmployeeAttributes).
			Return(returnedEmployee, returnedError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/99", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "99")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{"error": "Resource not found.","message": "Not found"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 6. error unique card number id - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "EMP002",
			"first_name": "Juana",
			"last_name": "Peralta",
			"warehouse_id": 2
		}`

		cardNumberId := "EMP002"
		firstName := "Juana"
		lastName := "Peralta"
		warehouseId := 2

		inputEmployeeAttributes := models.EmployeePatchRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{}
		returnedError := custom_errors.ErrUniqueAttributeViolationError
		returnedError.AttributeName = "card_number_id"
		returnedError.Value = "EMP002"

		mockService := new(service.MockEmployeeService)
		mockService.
			On("UpdateEmployeeById", 1, inputEmployeeAttributes).
			Return(returnedEmployee, returnedError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Invalid value {EMP002} for unique attribute {card_number_id}. Value already being used.", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 7. error warehouse doesn't exists - 409 Conflict", func(t *testing.T) {

		// Arrange
		body := `{
			"card_number_id": "99999",
			"first_name": "Juana",
			"last_name": "Peralta",
			"warehouse_id": 99
		}`

		cardNumberId := "99999"
		firstName := "Juana"
		lastName := "Peralta"
		warehouseId := 99

		inputEmployeeAttributes := models.EmployeePatchRequestBody{
			CardNumberId: &cardNumberId,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseId:  &warehouseId,
		}

		returnedEmployee := models.Employee{}
		returnedError := custom_errors.ErrForeignKeyViolation
		returnedError.ConstraintName = "warehouse_id"

		mockService := new(service.MockEmployeeService)
		mockService.
			On("UpdateEmployeeById", 1, inputEmployeeAttributes).
			Return(returnedEmployee, returnedError)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusConflict
		expectedBody := `{ "error": "Unknown entity identifier value: warehouse_id ", "message": "Conflict" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 8. typo in body - 422 UnprocessableEntity", func(t *testing.T) {

		// Arrange
		body := `{"cardNumberId": "55555",
			"firstName": "Juan",
			"lastName": "Peralta",
			"warehouseId": 99}`

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/employees/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchEmployee(w, req)

		// Assert
		expectedCode := http.StatusUnprocessableEntity
		expectedBody := `{"error":"Argument/s {card_number_id or first_name or last_name or warehouse_id} is/are mandatory", "message":"Unprocessable Entity"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

}

func TestDeleteEmployee(t *testing.T) {

	t.Run("Case 1. Success - 204 no content", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockEmployeeService)
		mockService.On("DeleteEmployee", 1).Return(nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/employees/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteEmployee(w, req)

		// Assert
		expectedCode := http.StatusNoContent
		expectedBody := `{"message":"Deleted succesfully"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 2. error - 404 Not Found", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockEmployeeService)
		mockService.On("DeleteEmployee", 1).Return(custom_errors.ErrNotFound)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/employees/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteEmployee(w, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{"error": "Resource not found.","message": "Not found"}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

	t.Run("Case 3. error - wrong id format - 400 Bad Request", func(t *testing.T) {

		// Arrange
		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/employees/1.0", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1.0")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		w := httptest.NewRecorder()

		// Act
		handler.DeleteEmployee(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{ "error": "Failed to decode url param {id}. Please verify format.", "message": "Bad request" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)

	})

}

func TestGetReportInboundOrders(t *testing.T) {

	t.Run("Case 1. return all employees inbound orders count successfully", func(t *testing.T) {
		// Arrange
		inboundOrders := []models.InboundOrderEmployee{
			{
				Employee: models.Employee{

					Id: 1,
					EmployeeAttributes: models.EmployeeAttributes{
						CardNumberId: "EMP001",
						FirstName:    "Raul",
						LastName:     "García",
						WarehouseId:  1,
					},
				},
				InboundOrdersCount: 3,
			},
			{
				Employee: models.Employee{

					Id: 2,
					EmployeeAttributes: models.EmployeeAttributes{
						CardNumberId: "EMP002",
						FirstName:    "Pepe",
						LastName:     "Sierra",
						WarehouseId:  2,
					},
				},
				InboundOrdersCount: 2,
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetReportInboundOrders", (*int)(nil)).Return(inboundOrders, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/reportInboundOrders", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetReportInboundOrders(w, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{ "data": [ { "id": 1, "card_number_id": "EMP001", "first_name": "Raul", "last_name": "García", "warehouse_id": 1, "inbound_orders_count": 3 }, { "id": 2, "card_number_id": "EMP002", "first_name": "Pepe", "last_name": "Sierra", "warehouse_id": 2, "inbound_orders_count": 2 }]}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Case 2. return employee inbound orders count successfully", func(t *testing.T) {
		// Arrange

		employeeId := 1

		inboundOrders := []models.InboundOrderEmployee{
			{
				Employee: models.Employee{

					Id: employeeId,
					EmployeeAttributes: models.EmployeeAttributes{
						CardNumberId: "EMP001",
						FirstName:    "Raul",
						LastName:     "García",
						WarehouseId:  1,
					},
				},
				InboundOrdersCount: 3,
			},
		}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetReportInboundOrders", &employeeId).Return(inboundOrders, nil)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/reportInboundOrders?id=1", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetReportInboundOrders(w, req)

		// Assert
		expectedCode := http.StatusOK
		expectedBody := `{ "data": [ { "id": 1, "card_number_id": "EMP001", "first_name": "Raul", "last_name": "García", "warehouse_id": 1, "inbound_orders_count": 3 }]}`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Case 3. error wrong id format - 400 Bad Request", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockEmployeeService)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/reportInboundOrders?id=1.0", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetReportInboundOrders(w, req)

		// Assert
		expectedCode := http.StatusBadRequest
		expectedBody := `{ "error": "Failed to decode query param {id}. Please verify format.", "message": "Bad request" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Case 4. error employees table empty - no query - 404 Not Found", func(t *testing.T) {

		// Arrange
		inboundOrders := []models.InboundOrderEmployee{}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetReportInboundOrders", (*int)(nil)).Return(inboundOrders, custom_errors.ErrNotFound)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/reportInboundOrders", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetReportInboundOrders(w, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{ "error": "Resource not found.", "message": "Not found" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("Case 5. error employee doesn't exist - 404 Not Found", func(t *testing.T) {

		// Arrange
		employeeId := 5
		inboundOrders := []models.InboundOrderEmployee{}

		mockService := new(service.MockEmployeeService)
		mockService.On("GetReportInboundOrders", &employeeId).Return(inboundOrders, custom_errors.ErrNotFound)

		handler := NewEmployeeController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/employees/reportInboundOrders?id=5", nil)

		w := httptest.NewRecorder()

		// Act
		handler.GetReportInboundOrders(w, req)

		// Assert
		expectedCode := http.StatusNotFound
		expectedBody := `{ "error": "Resource not found.", "message": "Not found" }`

		require.Equal(t, expectedCode, w.Code)
		require.JSONEq(t, expectedBody, w.Body.String())

		mockService.AssertExpectations(t)
	})

}
