package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"net/http"
	"time"

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

	// Validar formato de fecha
	if productRecordRequest.LastUpdateDate != nil {
		_, err := time.Parse("2006-01-02", *productRecordRequest.LastUpdateDate)
		if err != nil {
			utils.ResponseHttpError(w, &custom_errors.InvalidArgValueErr{
				Argument:  "last_update_date",
				Value:     *productRecordRequest.LastUpdateDate,
				ExtraInfo: "Date format must be YYYY-MM-DD",
			})
			return
		}
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
