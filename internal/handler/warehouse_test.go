package handler

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetWarehouse(t *testing.T) {
	t.Run("should return all warehouses successfully", func(t *testing.T) {

		expectedWarehouses := []models.Warehouse{
			{
				Id:                 1,
				WarehouseCode:      "Wh-1",
				Address:            "Street 1",
				Telephone:          "1234567890",
				MinimumCapacity:    100,
				MinimumTemperature: &[]float64{10}[0],
			},
			{
				Id:                 2,
				WarehouseCode:      "Wh-2",
				Address:            "Street 2",
				Telephone:          "1494945693",
				MinimumCapacity:    200,
				MinimumTemperature: &[]float64{20}[0],
			},
		}
		mockService := new(service.MockWarehouseService)
		mockService.On("GetAllWarehouses").Return(expectedWarehouses, nil)

		hd := NewWarehouseDefault(mockService)
		request := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		response := httptest.NewRecorder()

		hd.GetWarehouse(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res struct {
			Data []models.Warehouse `json:"data"`
		}
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Equal(t, expectedWarehouses, res.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetAllWarehouses").Return([]models.Warehouse{}, serviceError)

		hd := NewWarehouseDefault(mockService)
		request := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		response := httptest.NewRecorder()

		hd.GetWarehouse(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})
}

func Test_GetWarehouseById(t *testing.T) {
	t.Run("should return the warehouse successfully", func(t *testing.T) {
		expectedWarehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "Wh-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &[]float64{10}[0],
		}
		mockService := new(service.MockWarehouseService)
		mockService.On("GetWarehouseById", 1).Return(expectedWarehouse, nil)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

		response := httptest.NewRecorder()

		hd.GetWarehouseById(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res struct {
			Data models.Warehouse `json:"data"`
		}
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Equal(t, expectedWarehouse, res.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		errorService := custom_errors.ErrNotFound
		mockService.On("GetWarehouseById", 1).Return(models.Warehouse{}, errorService)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.GetWarehouseById(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodGet, "/warehouses/w", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "w")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.GetWarehouseById(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "GetWarehouseById", 0)

	})

	t.Run("should handle internal server error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		serviceError := errors.New("unexpected error")
		mockService.On("GetWarehouseById", 1).Return(models.Warehouse{}, serviceError)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodGet, "/warehouse/", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.GetWarehouseById(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})
}

func Test_CreateWarehouse(t *testing.T) {
	t.Run("should create the warehouse successfully", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "1234567890",
			"minimum_capacity": 100,
			"minimum_temperature": 10
		}`

		input := models.Warehouse{
			WarehouseCode:      "Wh-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &[]float64{10}[0],
		}

		returned := input
		returned.Id = 1

		mockService := new(service.MockWarehouseService)
		mockService.On("CreateWarehouse", input).Return(returned, nil)
		hd := NewWarehouseDefault(mockService)
		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusCreated, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res struct {
			Data models.Warehouse `json:"data"`
		}
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, returned, res.Data)

		mockService.AssertExpectations(t)

	})

	t.Run("should handle bad request error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPost, "/warehouses", nil)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "CreateWarehouse", 0)

	})
	t.Run("should handle bad request error for malformed body", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "1234567890",
			"minimum_capacity": 100,
			"minimum_temperature": 10,
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)
		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "CreateWarehouse", 0)

	})
	t.Run("should handle conflict error for unique constraint violation", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "1234567890",
			"minimum_capacity": 100,
			"minimum_temperature": 10
		}`
		input := models.Warehouse{
			WarehouseCode:      "Wh-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &[]float64{10}[0],
		}
		mockService := new(service.MockWarehouseService)
		mockService.On("CreateWarehouse", input).Return(models.Warehouse{}, custom_errors.ErrUniqueAttributeViolationError)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))

		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusConflict, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})

	t.Run("should handle unprocessable entity error for missing required fields", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1"
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "CreateWarehouse", 0)
	})

	t.Run("should handle unprocessable entity error for invalid fields values", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "",
			"minimum_capacity": 100,
			"minimum_temperature": 10
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "CreateWarehouse", 0)
	})

	t.Run("should return error if minimum_capacity is negative", func(t *testing.T) {
		body := `{
		"warehouse_code": "Wh-1",
		"address": "Street 1",
		"telephone": "1234567890",
		"minimum_capacity": -1,
		"minimum_temperature": 10
	}`

		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPost, "/warehouses", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateWarehouse(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "CreateWarehouse", 0)
	})

}

func Test_PatchWarehouse(t *testing.T) {
	t.Run("should update the warehouse successfully", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "1234567890",
			"minimum_capacity": 900,
			"minimum_temperature": 10
		}`
		address := "Street 1234"
		minimumCapacity := 900

		returned := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "Wh-1",
			Address:            address,
			Telephone:          "1234567890",
			MinimumCapacity:    minimumCapacity,
			MinimumTemperature: &[]float64{10}[0],
		}

		mockService := new(service.MockWarehouseService)
		mockService.On("UpdateWarehouseById", 1, mock.MatchedBy(func(wh models.Warehouse) bool {
			return wh.WarehouseCode == "Wh-1" &&
				wh.Telephone == "1234567890" &&
				wh.Address == "Street 1" &&
				wh.MinimumCapacity == 900 &&
				wh.MinimumTemperature != nil && *wh.MinimumTemperature == 10
		})).Return(returned, nil)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res struct {
			Data models.Warehouse `json:"data"`
		}
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, returned, res.Data)

		mockService.AssertExpectations(t)

	})

	t.Run("should handle bad request error for invalid url param", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/w", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "w")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "UpdateWarehouseById", 0)
	})

	t.Run("should handle bad request error for empty body", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "UpdateWarehouseById", 0)
	})

	t.Run("should handle malformed body", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "UpdateWarehouseById", 0)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		body := `{
			"warehouse_code": "1230",
			"address": "Fake Street 123",
			"telephone": "1234567890",
			"minimum_capacity": 100,
			"minimum_temperature": 10
		}`
		address := "Fake Street 123"
		warehouseCode := "1230"
		telephone := "1234567890"
		minimumCapacity := 100
		minimumTemperature := 10.0

		inputPatch := models.Warehouse{
			WarehouseCode:      warehouseCode,
			Address:            address,
			Telephone:          telephone,
			MinimumCapacity:    minimumCapacity,
			MinimumTemperature: &minimumTemperature,
		}
		mockService := new(service.MockWarehouseService)
		mockService.On("UpdateWarehouseById", 1, inputPatch).Return(models.Warehouse{}, custom_errors.ErrNotFound)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)

		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})
	t.Run("should handle unprocessable entity error for missing required fields", func(t *testing.T) {
		body := `{
		
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))
		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "UpdateWarehouseById", 0)
	})
	t.Run("should handle unprocessable entity error for invalid fields values", func(t *testing.T) {
		body := `{
			"warehouse_code": "Wh-1",
			"address": "Street 1",
			"telephone": "",
			"minimum_capacity": 100,
			"minimum_temperature": 10
		}`
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodPatch, "/warehouses/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.PatchWarehouse(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))
		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertNotCalled(t, "UpdateWarehouseById", 0)
	})
}

func Test_DeleteWarehouse(t *testing.T) {
	t.Run("should delete the warehouse successfully", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		mockService.On("DeleteWarehouse", 1).Return(nil)

		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.DeleteWarehouse(response, request)

		require.Equal(t, http.StatusNoContent, response.Code)
		require.Equal(t, "{}", response.Body.String())

		mockService.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		mockService.On("DeleteWarehouse", 1).Return(custom_errors.ErrNotFound)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodDelete, "/warehouses/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.DeleteWarehouse(response, request)

		require.Equal(t, http.StatusNotFound, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")

		mockService.AssertExpectations(t)
	})

	t.Run("should handle bad request error for invalid url param", func(t *testing.T) {
		mockService := new(service.MockWarehouseService)
		hd := NewWarehouseDefault(mockService)

		request := httptest.NewRequest(http.MethodDelete, "/warehouses/w", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "w")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))
		response := httptest.NewRecorder()

		hd.DeleteWarehouse(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		require.Equal(t, "application/json", response.Header().Get("Content-Type"))

		var res map[string]any
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Contains(t, res, "message")
		require.Contains(t, res, "error")
		mockService.AssertNotCalled(t, "DeleteWarehouse", 0)
	})
}

func Test_validateWarehouseAttributes(t *testing.T) {
	minimumTemperature := 5.0

	t.Run("should return error if warehouse_code is empty", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "warehouse_code")
	})

	t.Run("should return error if address is empty", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "WH-1",
			Address:            "   ",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "address")
	})

	t.Run("should return error if telephone is empty", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "WH-1",
			Address:            "Street 1",
			Telephone:          "",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "telephone")
	})

	t.Run("should return error if minimum capacity is negative", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "WH-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    -100,
			MinimumTemperature: &minimumTemperature,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "minimun_capacity")
	})

	t.Run("should return error if minimum temperature is nil", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "WH-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: nil,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "minimun_temperature")
	})

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		err := validateWarehouseAttributes(models.Warehouse{
			WarehouseCode:      "WH-1",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &minimumTemperature,
		})
		require.NoError(t, err)
	})
}
