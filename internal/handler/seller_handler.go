package handler

import (
	"app/internal/handler/dto"
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"github.com/bootcamp-go/web/response"
	"net/http"
	"strconv"
)

type SellerHandler struct {
	service service.SellerService
}

func NewSellerHandler(service service.SellerService) SellerHandler {
	return SellerHandler{service: service}
}

func (h *SellerHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	all, err := h.service.GetAll()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
}

func (h *SellerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	seller, err := h.service.GetById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
}

func (h *SellerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createSellerDto dto.CreateSellerDto
	createSellerDto, err := utils.InstantiateVarFromBody(&r.Body, createSellerDto)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	newSeller, err := h.service.Create(*createSellerDto.CompanyId, *createSellerDto.CompanyName, *createSellerDto.Address, *createSellerDto.Telephone)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": []models.Seller{newSeller},
	})
}

func (h *SellerHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var patchSellerDto dto.PatchSellerDto
	patchSellerDto, err := utils.InstantiateVarFromBody(&r.Body, patchSellerDto)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	seller, err := h.service.Patch(*patchSellerDto.Id, patchSellerDto.CompanyId, patchSellerDto.CompanyName, patchSellerDto.Address, patchSellerDto.Telephone)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
}

func (h *SellerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
