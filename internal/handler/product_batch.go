package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type ProductBatchController struct {
	sv service.ProductBatchService
}

func NewProductBatchController(sv service.ProductBatchService) *ProductBatchController {
	return &ProductBatchController{sv}
}

func (p *ProductBatchController) CreateProductBatch(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "CreateProductBatch", logger.LogStatusInProgress)

	var prodBatch models.ProductBatchRequest
	prodBatch, err := utils.InstantiateVarFromBody(&r.Body, prodBatch)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProductBatch", err, httpStatus)
		return
	}
	res, err := p.sv.CreateProductBatch(prodBatch)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProductBatch", err, httpStatus)
		return
	}

	utils.Log(r, "CreateProductBatch", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{"data": res})
}
