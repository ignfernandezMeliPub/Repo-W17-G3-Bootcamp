package handler

import (
	"app/internal/handler/dto"
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
)

type SellerHandler struct {
	service service.SellerService
}

func NewSellerHandler(service service.SellerService) SellerHandler {
	return SellerHandler{service: service}
}

func (h *SellerHandler) GetAllSellers(w http.ResponseWriter, _ *http.Request) {
	all, err := h.service.GetAllSellers()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
}

func (h *SellerHandler) GetSellerById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	seller, err := h.service.GetSellerById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
}

func (h *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) {
	var createSellerDto dto.CreateSellerDto
	createSellerDto, err := utils.InstantiateVarFromBody(&r.Body, createSellerDto)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	newSeller, err := h.service.CreateSeller(*createSellerDto.CompanyId, *createSellerDto.CompanyName, *createSellerDto.Address, *createSellerDto.Telephone)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": []models.Seller{newSeller},
	})
}

func (h *SellerHandler) PatchSeller(w http.ResponseWriter, r *http.Request) {
	var patchSellerDto dto.PatchSellerDto
	patchSellerDto, err := utils.InstantiateVarFromBody(&r.Body, patchSellerDto)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	seller, err := h.service.UpdateSellerById(*patchSellerDto.Id, patchSellerDto.CompanyId, patchSellerDto.CompanyName, patchSellerDto.Address, patchSellerDto.Telephone)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
}

func (h *SellerHandler) DeleteSeller(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	err = h.service.DeleteSellerById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
