package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

type ProductBatchController struct {
	sv service.ProductBatchService
}

func NewProductBatchController(sv service.ProductBatchService) *ProductBatchController {
	return &ProductBatchController{sv}
}

func (p *ProductBatchController) CreateProductBatch(w http.ResponseWriter, r *http.Request) {
	var prodBatch models.ProductBatchRequest
	prodBatch, err := utils.InstantiateVarFromBody(&r.Body, prodBatch)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	res, err := p.sv.CreateProductBatch(prodBatch)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"data": res})
	return
}
