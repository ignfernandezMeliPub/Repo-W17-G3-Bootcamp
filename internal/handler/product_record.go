package handler

import (
	"app/internal/handler/utils"
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

func (h *ProductRecordHandler) GetAllProductRecords(w http.ResponseWriter, r *http.Request) {
	productRecords, err := h.service.GetAllProductRecords()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": productRecords,
	})
}

func (h *ProductRecordHandler) CreateProductRecord(w http.ResponseWriter, r *http.Request) {
	var productRecordRequest models.ProductRecordRequest

	productRecordRequest, err := utils.InstantiateVarFromBody(&r.Body, productRecordRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	productRecordModel, err := h.service.CreateProductRecord(productRecordRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": productRecordModel,
	})

}
