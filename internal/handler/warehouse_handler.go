package handler

import (
	//"app/internal/handler/custom_errors"
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
)

func NewWarehouseDefault(sv service.IWarehouseService) *WarehouseDefault {
	return &WarehouseDefault{sv: sv}
}

type WarehouseDefault struct {
	sv service.IWarehouseService
}

func (h *WarehouseDefault) FindWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data []models.Warehouse
		data, err := h.sv.FindWarehouse()

		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{
				"message": "not warehouse were detected",
				"data":    nil,
			})
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *WarehouseDefault) FindWarehouseById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		data, err := h.sv.FindWarehouseById(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})

	}
}

func (h *WarehouseDefault) CreateWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request
		var wh models.Warehouse
		wh, err := utils.InstantiateVarFromBody(&r.Body, wh)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		//process
		data, err := h.sv.CreateWarehouse(wh)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		//response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "warehouse created successfully",
			"data":    data,
		})

	}
}

func (h *WarehouseDefault) UpdateWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		var wh models.Warehouse
		wh, err = utils.InstantiateVarFromBody(&r.Body, wh)
		if err != nil {
			status := http.StatusBadRequest
			message := "invalid request body"

			var decodeErr *custom_errors.DecodeError
			if errors.As(err, &decodeErr) {
				message = "invalid field type in request body"
			}

			response.JSON(w, status, map[string]any{
				"message": message,
				"error":   err.Error(),
			})
			return
		}
		//process
		var status int
		var message string

		data, err := h.sv.UpdateWarehouse(id, wh)
		if err != nil {
			if err.Error() == "warehose not found" {
				status = http.StatusNotFound
				message = "warehouse not found"

			} else if err.Error() == "minimun_capacity cannot be less than zero" {
				status = http.StatusBadRequest
				message = err.Error()

			}
			response.JSON(w, status, map[string]any{
				"message": message,
				"error":   err.Error(),
			})
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "warehouse updated successfully",
			"data":    data,
		})
	}
}

func (h *WarehouseDefault) DeleteWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//request
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		//process
		err = h.sv.DeleteWarehouse(id)
		if err != nil {
			status := http.StatusInternalServerError
			message := "warehouse could not be deleted"
			if err.Error() == "warehouse not found" {
				status = http.StatusNotFound
				message = "warehouse not found"
			}
			response.JSON(w, status, map[string]any{
				"message": message,
				"error":   err.Error(),
			})
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "warehouse delete successfully",
		})
	}
}
