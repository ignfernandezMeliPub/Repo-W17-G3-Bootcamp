package handler

import (
	//"app/internal/handler/custom_errors"
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
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
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *WarehouseDefault) FindWarehouseById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		data, err := h.sv.FindWarehouseById(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
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
			"data": data,
		})

	}
}

func (h *WarehouseDefault) UpdateWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		var wh models.Warehouse
		wh, err = utils.InstantiateVarFromBody(&r.Body, wh)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		data, err := h.sv.UpdateWarehouse(id, wh)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

func (h *WarehouseDefault) DeleteWarehouse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		err = h.sv.DeleteWarehouse(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusNoContent, map[string]any{})
	}
}
