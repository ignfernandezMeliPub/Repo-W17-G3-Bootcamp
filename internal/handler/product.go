package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
)

type ProductController struct {
	sv service.ProductServiceI
}

func NewProductController(service service.ProductServiceI) *ProductController {
	return &ProductController{sv: service}
}

func (h *ProductController) GetAllProducts(w http.ResponseWriter, _ *http.Request) {
	products, err := h.sv.GetAllProducts()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
}

func (h *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	product, err := h.sv.GetProductById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (h *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productRequest models.ProductRequest

	productRequest, err := utils.InstantiateVarFromBody(&r.Body, productRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	product, err := h.sv.CreateProduct(productRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": product,
	})
}

func (h *ProductController) PatchProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	var productPatchRequest models.ProductPatchRequest
	productPatchRequest, err = utils.InstantiateVarFromBody(&r.Body, productPatchRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	productPatchRequest.Id = id

	product, err := h.sv.UpdateProductById(productPatchRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (h *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	err = h.sv.DeleteProductById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, map[string]any{})
}

func (h *ProductController) GetReportRecords(w http.ResponseWriter, r *http.Request) {

	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	var reportRecords []models.ProductRecordReport

	if id == nil {
		reportRecords, err = h.sv.GetAllReportRecords()
	} else {
		reportRecords, err = h.sv.GetReportRecords(*id)
	}

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": reportRecords,
	})
}
