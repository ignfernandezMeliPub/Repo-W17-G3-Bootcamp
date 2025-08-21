package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type ProductRecordHandler struct {
	service service.ProductRecordService
}

func NewProductRecordHandler(service service.ProductRecordService) ProductRecordHandler {
	return ProductRecordHandler{service: service}
}

func (h *ProductRecordHandler) GetAllProductRecords(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetAllProductRecords", logger.LogStatusInProgress)

	productRecords, err := h.service.GetAllProductRecords()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllProductRecords", err, httpStatus)
		return
	}

	utils.Log(r, "GetAllProductRecords", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": productRecords,
	})
	return
}

func (h *ProductRecordHandler) CreateProductRecord(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateProductRecord", logger.LogStatusInProgress)

	var productRecordRequest models.ProductRecordRequest

	productRecordRequest, err = utils.InstantiateVarFromBody(&r.Body, productRecordRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProductRecord", err, httpStatus)
		return
	}

	productRecordModel, err := h.service.CreateProductRecord(productRecordRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProductRecord", err, httpStatus)
		return
	}

	utils.Log(r, "CreateProductRecord", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": productRecordModel,
	})
	return
}
