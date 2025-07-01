package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
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
			"message": "success",
			"data":    products,
		})
	}
}

func (h *ProductController) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID format",
				"error":   err.Error(),
			})
			return
		}

		product, err := h.sv.GetProductById(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    product,
		})
	}
}

func (h *ProductController) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var productRequest models.ProductRequest
		var err error

		productRequest, err = utils.InstantiateVarFromBody(&r.Body, productRequest)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		product, err := h.sv.CreateProduct(productRequest)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error creating product",
				"error":   err.Error(),
			})
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "Product created successfully",
			"data":    product,
		})
	}
}

func (h *ProductController) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID format",
				"error":   err.Error(),
			})
			return
		}

		var productRequest models.ProductRequest
		productRequest, err = utils.InstantiateVarFromBody(&r.Body, productRequest)
		if err != nil {
			status := http.StatusBadRequest
			message := "invalid request body"

			var decodeErr *custom_errors.DecodeError
			if errors.As(err, &decodeErr) {
				message = "invalid field type in request body"
			}

			response.JSON(w, status, map[string]any{
				"message": message,
				"error":   err.Error(),
			})
			return
		}

		productRequest.Id = id

		product, err := h.sv.UpdateProduct(productRequest)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error updating product",
				"error":   err.Error(),
			})
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product updated successfully",
			"data":    product,
		})
	}
}

func (h *ProductController) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID format",
				"error":   err.Error(),
			})
			return
		}

		err = h.sv.DeleteProduct(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Error deleting product",
				"error":   err.Error(),
			})
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "Product deleted successfully",
		})
	}
}
