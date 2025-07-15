package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"github.com/bootcamp-go/web/request"
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
	if err := request.JSON(r, &prodBatch); err != nil {
		utils.ResponseHttpError(w, &custom_errors.InvalidBodyError{})
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
