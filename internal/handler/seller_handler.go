package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"github.com/bootcamp-go/web/response"
	"net/http"
	"strconv"
)

type SellerHandler struct {
	service service.SellerService
}

func (h *SellerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	all, err := h.service.GetAll()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "ok",
		"data":    all,
	})
}

func (h *SellerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err) // TODO Manejar este error
		return
	}

	seller, err := h.service.GetById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "ok",
		"data":    seller,
	})
}

func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {

}

func (h *SellerHandler) Patch(w http.ResponseWriter, r *http.Request) {

}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err) // TODO Manejar este error
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
