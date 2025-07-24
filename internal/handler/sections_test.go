package handler

import (
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSectionsController_Create(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSect := service.NewSectionsService(dbMock)
	hdSect := NewSectionsController(svSect)

	dbMock.On("CreateSection", models.Section{
		ID: 0, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	dbMock.On("CreateSection", models.Section{
		ID: 0, SectionNumber: 2, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}).Return(models.Section{}, &custom_errors.UniqueAttributeViolationErr{
		AttributeName: "section_number", Value: "2",
	})

	t.Run("create_ok", func(t *testing.T) {
		body := strings.NewReader(`{
			"section_number": 1,
			"current_temperature": 1,
			"minimum_temperature": 1,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 1,
			"warehouse_id": 1,
			"product_type_id": 1
		}`)
		req := httptest.NewRequest(http.MethodPost, "/sections", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
			"data": {
				"id": 1,
				"section_number": 1,
				"current_temperature": 1,
				"minimum_temperature": 1,
				"current_capacity": 1,
				"minimum_capacity": 1,
				"maximum_capacity": 1,
				"warehouse_id": 1,
				"product_type_id": 1
			}
		}`

		hdSect.CreateSection(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("create_fail", func(t *testing.T) {
		body := strings.NewReader(`{
			"section_number": 2,
			"current_temperature": 1,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 1,
			"warehouse_id": 1,
			"product_type_id": 1
		}`)
		req := httptest.NewRequest(http.MethodPost, "/sections", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
    		"error": "Argument/s {minimum_temperature} is/are mandatory",
    		"message": "Unprocessable Entity"
		}`

		hdSect.CreateSection(w, req)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("create_conflict", func(t *testing.T) {
		body := strings.NewReader(`{
			"section_number": 2,
			"current_temperature": 1,
			"minimum_temperature": 1,
			"current_capacity": 1,
			"minimum_capacity": 1,
			"maximum_capacity": 1,
			"warehouse_id": 1,
			"product_type_id": 1
		}`)
		req := httptest.NewRequest(http.MethodPost, "/sections", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
    		"error": "Invalid value {2} for unique attribute {section_number}. Value already being used.",
    		"message": "Conflict"
		}`

		hdSect.CreateSection(w, req)
		require.Equal(t, http.StatusConflict, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})
}

func TestSectionsController_Read(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSect := service.NewSectionsService(dbMock)
	hdSect := NewSectionsController(svSect)
	dbMock.On("GetAllSections").Return([]models.Section{
		{
			ID: 1, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
		},
		{
			ID: 2, SectionNumber: 2, CurrentTemperature: 2, MinimumTemperature: 2, CurrentCapacity: 2, MinimumCapacity: 2, MaximumCapacity: 2, WarehouseId: 2, ProductTypeId: 2,
		},
	}, nil)

	dbMock.On("GetSectionById", 1).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	dbMock.On("GetSectionById", 3).Return(models.Section{}, custom_errors.ErrNotFound)

	t.Run("find_all", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sections", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
			"data":[{
				"id": 1,
				"section_number": 1,
				"current_temperature": 1,
				"minimum_temperature": 1,
				"current_capacity": 1,
				"minimum_capacity": 1,
				"maximum_capacity": 1,
				"warehouse_id": 1,
				"product_type_id": 1
			},
			{
				"id": 2,
				"section_number": 2,
				"current_temperature": 2,
				"minimum_temperature": 2,
				"current_capacity": 2,
				"minimum_capacity": 2,
				"maximum_capacity": 2,
				"warehouse_id": 2,
				"product_type_id": 2
			}
			]
		}`

		hdSect.GetSections(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sections/3", nil)
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "3")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		expected := `{
			"error":"Resource not found.", 
			"message":"Not found"
		}`

		hdSect.GetSectionById(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("find_by_id_existent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sections/1", nil)
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		expected := `{
			"data":{
				"id": 1,
				"section_number": 1,
				"current_temperature": 1,
				"minimum_temperature": 1,
				"current_capacity": 1,
				"minimum_capacity": 1,
				"maximum_capacity": 1,
				"warehouse_id": 1,
				"product_type_id": 1
			}
		}`

		hdSect.GetSectionById(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})
}

func TestSectionsController_Update(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSect := service.NewSectionsService(dbMock)
	hdSect := NewSectionsController(svSect)

	dbMock.On("GetSectionById", 1).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 2, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	dbMock.On("UpdateSectionById", models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}).Return(models.Section{
		ID: 1, SectionNumber: 1, CurrentTemperature: 1, MinimumTemperature: 1, CurrentCapacity: 1, MinimumCapacity: 1, MaximumCapacity: 1, WarehouseId: 1, ProductTypeId: 1,
	}, nil)

	dbMock.On("GetSectionById", 2).Return(models.Section{}, custom_errors.ErrNotFound)

	t.Run("update_ok", func(t *testing.T) {
		body := strings.NewReader(`{
			"current_temperature": 1
		}`)
		req := httptest.NewRequest(http.MethodPatch, "/sections/1", body)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{
			"data": {
				"id": 1,
				"section_number": 1,
				"current_temperature": 1,
				"minimum_temperature": 1,
				"current_capacity": 1,
				"minimum_capacity": 1,
				"maximum_capacity": 1,
				"warehouse_id": 1,
				"product_type_id": 1
			}
		}`

		hdSect.PatchSection(w, req)
		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})

	t.Run("update_non_existent", func(t *testing.T) {
		body := strings.NewReader(`{
			"current_temperature": 1
		}`)
		req := httptest.NewRequest(http.MethodPatch, "/sections/2", body)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		expected := `{"error":"Resource not found.", "message":"Not found"}`

		hdSect.PatchSection(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
		require.JSONEq(t, expected, string(w.Body.Bytes()))
	})
}

func TestSectionsController_Delete(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSect := service.NewSectionsService(dbMock)
	hdSect := NewSectionsController(svSect)
	dbMock.On("DeleteSectionById", 1).Return(nil)
	dbMock.On("DeleteSectionById", 2).Return(custom_errors.ErrNotFound)

	t.Run("delete_ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sections/1", nil)
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		hdSect.DeleteSection(w, req)
		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "", string(w.Body.Bytes()))
	})

	t.Run("delete_non_existent", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sections/2", nil)
		req.Header.Set("Content-Type", "application/json")
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "2")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
		w := httptest.NewRecorder()

		expected := `{"error":"Resource not found.","message":"Not found"}`

		hdSect.DeleteSection(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, expected, string(w.Body.Bytes()))
	})
}
