package handler

import (
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProductBatchController_CreateProductBatch(t *testing.T) {
	dbMock := repository.NewProductBatchMock()
	svPb := service.NewProductBatchService(dbMock)
	hdPb := NewProductBatchController(svPb)

	dbMock.On("CreateProductBatch", models.ProductBatch{
		ID: 0, BatchNumber: 1, CurrentQuantity: 80, CurrentTemperature: 4, DueDate: "2024-07-31", InitialQuantity: 100, ManufacturingDate: "2024-06-20", ManufacturingHour: 9, MinimumTemperature: 2, ProductId: 25, SectionId: 3,
	}).Return(models.ProductBatch{
		ID: 1, BatchNumber: 1, CurrentQuantity: 80, CurrentTemperature: 4, DueDate: "2024-07-31", InitialQuantity: 100, ManufacturingDate: "2024-06-20", ManufacturingHour: 9, MinimumTemperature: 2, ProductId: 25, SectionId: 3,
	}, nil)

	dbMock.On("CreateProductBatch", models.ProductBatch{
		ID: 0, BatchNumber: 2, CurrentQuantity: 80, CurrentTemperature: 4, DueDate: "2024-07-31", InitialQuantity: 100, ManufacturingDate: "2024-06-20", ManufacturingHour: 9, MinimumTemperature: 2, ProductId: 25, SectionId: 3,
	}).Return(models.ProductBatch{}, &custom_errors.UniqueAttributeViolationErr{
		AttributeName: "batch_number", Value: "2",
	})

	t.Run("create_ok", func(t *testing.T) {
		body := strings.NewReader(`{
			"data": {
				"batch_number": 1,
				"current_quantity": 80,
				"current_temperature": 4,
				"due_date": "2024-07-31",
				"initial_quantity": 100,
				"manufacturing_date": "2024-06-20",
				"manufacturing_hour": 9,
				"minimum_temperature":2,
				"product_id": 25,
				"section_id": 3
			  }
			}`)
		req := httptest.NewRequest(http.MethodPost, "/productBatches", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
			"data": {
				"id": 1,
				"batch_number": 1,
				"current_quantity": 80,
				"current_temperature": 4,
				"due_date": "2024-07-31",
				"initial_quantity": 100,
				"manufacturing_date": "2024-06-20",
				"manufacturing_hour": 9,
				"minimum_temperature": 2,
				"product_id": 25,
				"section_id": 3
			}
		}`

		hdPb.CreateProductBatch(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("create_fail", func(t *testing.T) {
		body := strings.NewReader(`{
			"data": {
				"batch_number": 1,
				"current_quantity": 80,
				"current_temperature": 4,
				"due_date": "2024-07-31",
				"initial_quantity": 100,
				"manufacturing_date": "2024-06-20",
				"manufacturing_hour": 9,
				"minimum_temperature":2,
				"product_id": 25
			  }
			}`)
		req := httptest.NewRequest(http.MethodPost, "/productBatches", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{"error":"Argument/s {section_id} is/are mandatory", "message":"Unprocessable Entity"}`

		hdPb.CreateProductBatch(w, req)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("create_conflict", func(t *testing.T) {
		body := strings.NewReader(`{
			"data": {
				"batch_number": 2,
				"current_quantity": 80,
				"current_temperature": 4,
				"due_date": "2024-07-31",
				"initial_quantity": 100,
				"manufacturing_date": "2024-06-20",
				"manufacturing_hour": 9,
				"minimum_temperature":2,
				"product_id": 25,
				"section_id": 3
			  }
			}`)
		req := httptest.NewRequest(http.MethodPost, "/productBatches", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{"error":"Invalid value {2} for unique attribute {batch_number}. Value already being used.", "message":"Conflict"}`

		hdPb.CreateProductBatch(w, req)
		require.Equal(t, http.StatusConflict, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})
}
