package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
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
	utils.Log(r, "GetWarehouse", logger.LogStatusInProgress)

	var data []models.Warehouse
	data, err := h.sv.GetAllWarehouses()

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetWarehouse", err, httpStatus)
		return
	}

	utils.Log(r, "GetWarehouse", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *WarehouseDefault) GetWarehouseById(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "GetWarehouseById", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetWarehouseById", err, httpStatus)
		return
	}

	data, err := h.sv.GetWarehouseById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetWarehouseById", err, httpStatus)
		return
	}

	utils.Log(r, "GetWarehouseById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *WarehouseDefault) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "CreateWarehouse", logger.LogStatusInProgress)

	// request
	var wh models.Warehouse
	wh, err := utils.InstantiateVarFromBody(&r.Body, wh)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateWarehouse", err, httpStatus)
		return
	}
	err = validateWarehouseAttributes(wh)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateWarehouse", err, httpStatus)
		return
	}
	// process
	data, err := h.sv.CreateWarehouse(wh)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateWarehouse", err, httpStatus)
		return
	}

	utils.Log(r, "CreateWarehouse", logger.LogStatusSuccess)

	// response
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
}

func (h *WarehouseDefault) PatchWarehouse(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "PatchWarehouse", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchWarehouse", err, httpStatus)
		return
	}

	var wh models.Warehouse
	wh, err = utils.InstantiateVarFromBody(&r.Body, wh)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchWarehouse", err, httpStatus)
		return
	}

	data, err := h.sv.UpdateWarehouseById(id, wh)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchWarehouse", err, httpStatus)
		return
	}

	utils.Log(r, "PatchWarehouse", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *WarehouseDefault) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "DeleteWarehouse", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteWarehouse", err, httpStatus)
		return
	}

	err = h.sv.DeleteWarehouse(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteWarehouse", err, httpStatus)
		return
	}

	utils.Log(r, "DeleteWarehouse", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, map[string]any{})
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
