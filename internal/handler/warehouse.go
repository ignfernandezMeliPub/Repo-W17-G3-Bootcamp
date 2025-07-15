package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/bootcamp-go/web/response"
)

func NewWarehouseDefault(sv service.IWarehouseService) *WarehouseDefault {
	return &WarehouseDefault{sv: sv}
}

type WarehouseDefault struct {
	sv service.IWarehouseService
}

func (h *WarehouseDefault) GetWarehouse(w http.ResponseWriter, r *http.Request) {

	var data []models.Warehouse
	data, err := h.sv.GetAllWarehouses()

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
	return
}

func (h *WarehouseDefault) GetWarehouseById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	data, err := h.sv.GetWarehouseById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
	return
}

func (h *WarehouseDefault) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	// request
	var wh models.Warehouse
	wh, err := utils.InstantiateVarFromBody(&r.Body, wh)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	err = validateWarehouseAttributes(wh)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	// process
	data, err := h.sv.CreateWarehouse(wh)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// response
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
	return
}

func (h *WarehouseDefault) PatchWarehouse(w http.ResponseWriter, r *http.Request) {

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

	data, err := h.sv.UpdateWarehouseById(id, wh)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
	return
}

func (h *WarehouseDefault) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
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
	return
}

func validateWarehouseAttributes(wh models.Warehouse) error {
	if strings.TrimSpace(wh.WarehouseCode) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "warehouse_code"}
	}
	if strings.TrimSpace(wh.Address) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "address"}
	}
	if strings.TrimSpace(wh.Telephone) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "telephone"}
	}
	if wh.MinimumCapacity < 0 {
		return &custom_errors.MandatoryArgMissingErr{Argument: "minimun_capacity"}
	}
	if wh.MinimumTemperature == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "minimun_temperature"}
	}
	return nil
}
