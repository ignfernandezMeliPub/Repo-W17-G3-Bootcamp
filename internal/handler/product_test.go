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

func TestGetAllProducts(t *testing.T) {
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
		var res struct {
			Data []models.Product `json:"data"`
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

		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, expectedProducts, res.Data)
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
		//ver si esto es necesario
		require.Contains(t, response, "message")
		require.Contains(t, response, "error")
		mockService.AssertExpectations(t)
	})
}

func TestGetProductById(t *testing.T) {
	t.Run("should return the product successfully", func(t *testing.T) {
		// Arrange
		expectedResponse := `{"data": {"id": 1, "product_code": "1234567890", "description": "Product 1", "product_type_id": 1, "width": 10, "height": 20, "length": 30, "net_weight": 100, "expiration_rate": 1, "recommended_freezing_temperature": 10, "freezing_rate": 1, "seller_id": null}}`
		expectedProduct := models.Product{
			ID:                             1,
			ProductCode:                    "1234567890",
			Description:                    "Product 1",
			ProductTypeId:                  1,
			Width:                          10.0,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      100.0,
			ExpirationRate:                 1,
			RecommendedFreezingTemperature: 10.0,
			FreezingRate:                   1,
			SellerId:                       nil,
		}
		mockService := new(service.MockProductService)
		mockService.On("GetProductById", 1).Return(expectedProduct, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.GetProductById(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.JSONEq(t, expectedResponse, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("should return error when product is not found", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		serviceError := custom_errors.ErrNotFound
		mockService.On("GetProductById", 1).Return(models.Product{}, serviceError)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.GetProductById(w, req)

		// Assert
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("should return bad request when id is not a number", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("GetProductById", 1).Return(models.Product{}, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.GetProductById(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

}

func TestCreateProduct(t *testing.T) {

	body := `{
		"product_code": "1234567890",
		"description": "Product 1",
		"product_type_id": 1,
		"width": 10,
		"height": 20,
		"length": 30,
		"net_weight": 100,
		"expiration_rate": 1,
		"recommended_freezing_temperature": 10,
		"freezing_rate": 1
	}`

	bodyInvalid := `{
		"product_code": "1234567890",
		"freezing_rate": 1,
	}`

	bodyIncomplete := `{
		"product_code": "1234567890",
		"freezing_rate": 1
	}`

	bodyInvalidFields := `{
		"product_code": "1234567890",
		"description": "Product 1",
		"product_type_id": 1,
		"width": -10,
		"height": 20,
		"length": 30,
		"net_weight": 100,
		"expiration_rate": 1,
		"recommended_freezing_temperature": 10,
		"freezing_rate": 1
	}`

	productCode := "1234567890"
	description := "Product 1"
	productTypeId := 1
	width := 10.0
	height := 20.0
	length := 30.0
	netWeight := 100.0
	expirationRate := 1
	recommendedFreezingTemperature := 10.0
	freezingRate := 1

	inputProduct := models.ProductRequest{
		ProductCode:                    &productCode,
		Description:                    &description,
		ProductTypeId:                  &productTypeId,
		Width:                          &width,
		Height:                         &height,
		Length:                         &length,
		NetWeight:                      &netWeight,
		ExpirationRate:                 &expirationRate,
		RecommendedFreezingTemperature: &recommendedFreezingTemperature,
		FreezingRate:                   &freezingRate,
		SellerId:                       nil,
	}

	returnedProduct := models.Product{
		ID:                             1,
		ProductCode:                    productCode,
		Description:                    description,
		ProductTypeId:                  productTypeId,
		Width:                          width,
		Height:                         height,
		Length:                         length,
		NetWeight:                      netWeight,
		ExpirationRate:                 expirationRate,
		RecommendedFreezingTemperature: recommendedFreezingTemperature,
		FreezingRate:                   freezingRate,
		SellerId:                       nil,
	}
	t.Run("should create the product successfully", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockProductService)
		mockService.On("CreateProduct", inputProduct).Return(returnedProduct, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusCreated, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var res struct {
			Data models.Product `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &res)

		require.NoError(t, err)
		require.Equal(t, returnedProduct, res.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is empty", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("CreateProduct", models.ProductRequest{}).Return(models.Product{}, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "message")
		require.Contains(t, response, "error")
		mockService.AssertNotCalled(t, "CreateProduct", 0)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockProductService)
		mockService.On("CreateProduct", inputProduct).Return(models.Product{}, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyInvalid))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "CreateProduct", 0)
	})

	t.Run("should return conflict error when product code already exists", func(t *testing.T) {
		// Arrange

		mockService := new(service.MockProductService)
		mockService.On("CreateProduct", inputProduct).Return(models.Product{}, custom_errors.ErrUniqueAttributeViolationError)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

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

	t.Run("should return unprocessable entity when body is incomplete", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyIncomplete))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "message")
		require.Contains(t, response, "error")
		mockService.AssertNotCalled(t, "CreateProduct", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		// Arrange
		// Create the expected ProductRequest with the invalid width value
		invalidWidth := -10.0
		expectedInvalidProduct := models.ProductRequest{
			ProductCode:                    &productCode,
			Description:                    &description,
			ProductTypeId:                  &productTypeId,
			Width:                          &invalidWidth,
			Height:                         &height,
			Length:                         &length,
			NetWeight:                      &netWeight,
			ExpirationRate:                 &expirationRate,
			RecommendedFreezingTemperature: &recommendedFreezingTemperature,
			FreezingRate:                   &freezingRate,
			SellerId:                       nil,
		}

		mockService := new(service.MockProductService)
		mockService.On("CreateProduct", expectedInvalidProduct).Return(models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "width", Value: -10, ExtraInfo: "Width must be greater than 0"})

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(bodyInvalidFields))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Contains(t, response, "message")
		require.Contains(t, response, "error")
		require.Equal(t, "Unprocessable Entity", response["message"])
		require.Contains(t, response["error"], "width")
		require.Contains(t, response["error"], "Width must be greater than 0")
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request when seller id is not a number", func(t *testing.T) {
		// Arrange
		body := `{
			"product_code": "1234567890",
			"description": "Product 1",
			"product_type_id": 1,
			"width": 10,
			"height": 20,
			"length": 30,
			"net_weight": 100,
			"expiration_rate": 1,
			"recommended_freezing_temperature": 10,
			"freezing_rate": 1,
			"seller_id": "a"
		}`
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.CreateProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "CreateProduct", 0)
	})

}

func TestPatchProduct(t *testing.T) {

	body := `{
		"product_code": "1234567890",
		"description": "Product 1",
		"product_type_id": 1,
		"width": 10,
		"height": 20,
		"length": 30,
		"net_weight": 100,
		"expiration_rate": 1,
		"recommended_freezing_temperature": 10,
		"freezing_rate": 1
	}`

	bodyInvalidFields := `{
		"product_code": "1234567890",
		"description": "Product 1",
		"product_type_id": 1,
		"width": -10
	}`

	productCode := "1234567890"
	description := "Product 1"
	productTypeId := 1
	width := 10.0
	height := 20.0
	length := 30.0
	netWeight := 100.0
	expirationRate := 1
	recommendedFreezingTemperature := 10.0
	freezingRate := 1

	inputProduct := models.ProductPatchRequest{
		Id:                             1,
		ProductCode:                    &productCode,
		Description:                    &description,
		ProductTypeId:                  &productTypeId,
		Width:                          &width,
		Height:                         &height,
		Length:                         &length,
		NetWeight:                      &netWeight,
		ExpirationRate:                 &expirationRate,
		RecommendedFreezingTemperature: &recommendedFreezingTemperature,
		FreezingRate:                   &freezingRate,
		SellerId:                       nil,
	}

	returnedProduct := models.Product{
		ID:                             1,
		ProductCode:                    productCode,
		Description:                    description,
		ProductTypeId:                  productTypeId,
		Width:                          width,
		Height:                         height,
		Length:                         length,
		NetWeight:                      netWeight,
		ExpirationRate:                 expirationRate,
		RecommendedFreezingTemperature: recommendedFreezingTemperature,
		FreezingRate:                   freezingRate,
		SellerId:                       nil,
	}

	t.Run("should patch the product successfully", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("UpdateProductById", inputProduct).Return(returnedProduct, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/products/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchProduct(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("should return bad request when id is not a number", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/products/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "UpdateProductById", 0)
	})

	t.Run("should return bad request when body is empty", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/products/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "UpdateProductById", 0)
	})

	t.Run("should handle unprocessable content error for invalid fields values", func(t *testing.T) {
		width := -10.0
		inputPatchProduct := models.ProductPatchRequest{
			Id:            1,
			ProductCode:   &productCode,
			Description:   &description,
			ProductTypeId: &productTypeId,
			Width:         &width,
		}
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("UpdateProductById", inputPatchProduct).Return(models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "width", Value: -10, ExtraInfo: "Width must be greater than 0"})

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/products/1", bytes.NewBufferString(bodyInvalidFields))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler.PatchProduct(w, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	})

	t.Run("should return bad request when seller id is not a number", func(t *testing.T) {
		// Arrange
		body := `{
			"product_code": "1234567890",
			"description": "Product 1",
			"product_type_id": 1,
			"width": 10,
			"height": 20,
			"length": 30,
			"net_weight": 100,
			"expiration_rate": 1,
			"recommended_freezing_temperature": 10,
			"freezing_rate": 1,
			"seller_id": "a"
		}`
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodPatch, "/products/1", bytes.NewBufferString(body))
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.PatchProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "UpdateProductById", 0)
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("should delete the product successfully", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("DeleteProductById", 1).Return(nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.DeleteProduct(w, req)

		// Assert
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "{}", w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when product is not found", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)
		mockService.On("DeleteProductById", 30).Return(custom_errors.ErrNotFound)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/products/30", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "30")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.DeleteProduct(w, req)

		// Assert
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("should return bad request when id is not a number", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/products/a", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "a")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		// Act
		handler.DeleteProduct(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "DeleteProductById", 0)
	})
}

func TestGetReportRecords(t *testing.T) {

	t.Run("should return the report records successfully", func(t *testing.T) {
		// Arrange

		expectedReportRecords := []models.ProductRecordReport{
			{
				ProductID:    1,
				Description:  "Product 1",
				RecordsCount: 10,
			},
			{
				ProductID:    2,
				Description:  "Product 2",
				RecordsCount: 20,
			},
		}
		mockService := new(service.MockProductService)
		mockService.On("GetReportRecords", (*int)(nil)).Return(expectedReportRecords, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/reportRecords", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetReportRecords(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var res struct {
			Data []models.ProductRecordReport `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, expectedReportRecords, res.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request when id is not a number", func(t *testing.T) {
		// Arrange
		mockService := new(service.MockProductService)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/reportRecords?id=a", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetReportRecords(w, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		mockService.AssertNotCalled(t, "GetReportRecords", 0)
	})

	t.Run("should return the report by product id successfully", func(t *testing.T) {
		// Arrange
		expectedReportRecords := []models.ProductRecordReport{
			{
				ProductID:    1,
				Description:  "Product 1",
				RecordsCount: 10,
			},
		}
		productId := 1
		mockService := new(service.MockProductService)
		mockService.On("GetReportRecords", &productId).Return(expectedReportRecords, nil)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/reportRecords?id=1", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetReportRecords(w, req)

		// Assert
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var res struct {
			Data []models.ProductRecordReport `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, expectedReportRecords, res.Data)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when product is not found", func(t *testing.T) {
		// Arrange
		productId := 1
		mockService := new(service.MockProductService)
		mockService.On("GetReportRecords", &productId).Return([]models.ProductRecordReport{}, custom_errors.ErrNotFound)

		handler := NewProductController(mockService)

		req := httptest.NewRequest(http.MethodGet, "/products/reportRecords?id=1", nil)
		w := httptest.NewRecorder()

		// Act
		handler.GetReportRecords(w, req)

		// Assert
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

}
