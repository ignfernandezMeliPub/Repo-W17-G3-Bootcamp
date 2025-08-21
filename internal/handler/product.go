package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
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

func (h *ProductController) GetAllProducts(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetAllProducts", logger.LogStatusInProgress)

	products, err := h.sv.GetAllProducts()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllProducts", err, httpStatus)
		return
	}

	utils.Log(r, "GetAllProducts", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
	return
}

func (h *ProductController) GetProductById(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetProductById", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetProductById", err, httpStatus)
		return
	}

	product, err := h.sv.GetProductById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetProductById", err, httpStatus)
		return
	}

	utils.Log(r, "GetProductById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
	return
}

func (h *ProductController) CreateProduct(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateProduct", logger.LogStatusInProgress)

	var productRequest models.ProductRequest

	productRequest, err = utils.InstantiateVarFromBody(&r.Body, productRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProduct", err, httpStatus)
		return
	}

	product, err := h.sv.CreateProduct(productRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateProduct", err, httpStatus)
		return
	}

	utils.Log(r, "CreateProduct", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": product,
	})
	return
}

func (h *ProductController) PatchProduct(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "PatchProduct", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchProduct", err, httpStatus)
		return
	}

	var productPatchRequest models.ProductPatchRequest
	productPatchRequest, err = utils.InstantiateVarFromBody(&r.Body, productPatchRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchProduct", err, httpStatus)
		return
	}

	productPatchRequest.Id = id

	product, err := h.sv.UpdateProductById(productPatchRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchProduct", err, httpStatus)
		return
	}

	utils.Log(r, "PatchProduct", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
	return
}

func (h *ProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "DeleteProduct", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteProduct", err, httpStatus)
		return
	}

	err = h.sv.DeleteProductById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteProduct", err, httpStatus)
		return
	}

	utils.Log(r, "DeleteProduct", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, map[string]any{})
	return
}

func (h *ProductController) GetReportRecords(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetReportRecords", logger.LogStatusInProgress)

	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetReportRecords", err, httpStatus)
		return
	}

	var reportRecords []models.ProductRecordReport

	reportRecords, err = h.sv.GetReportRecords(id)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetReportRecords", err, httpStatus)
		return
	}

	utils.Log(r, "GetReportRecords", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": reportRecords,
	})
	return
}
