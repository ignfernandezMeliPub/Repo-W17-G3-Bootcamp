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

func (h *ProductController) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := h.sv.GetAllProducts()
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": products,
		})
	}
}

func (h *ProductController) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (h *ProductController) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (h *ProductController) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		product, err := h.sv.UpdateProduct(productPatchRequest)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": product,
		})
	}
}

func (h *ProductController) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		err = h.sv.DeleteProduct(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{})
	}
}
