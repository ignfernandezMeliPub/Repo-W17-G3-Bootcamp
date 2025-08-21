package handler

import (
	"app/internal/handler/dto"
	"app/internal/handler/utils"
	"app/internal/logger"
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

func (h *SellerHandler) GetAllSellers(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetAllSellers", logger.LogStatusInProgress)

	all, err := h.service.GetAllSellers()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllSellers", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetAllSellers", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
	return nil
}

func (h *SellerHandler) GetSellerById(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetSellerById", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetSellerById", err, httpStatus)
		return nil
	}

	seller, err := h.service.GetSellerById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetSellerById", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetSellerById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
	return nil
}

func (h *SellerHandler) CreateSeller(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateSeller", logger.LogStatusInProgress)

	var createSellerDto dto.CreateSellerDto
	createSellerDto, err = utils.InstantiateVarFromBody(&r.Body, createSellerDto)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateSeller", err, httpStatus)
		return nil
	}

	newSeller, err := h.service.CreateSeller(*createSellerDto.CompanyId, *createSellerDto.CompanyName, *createSellerDto.Address, *createSellerDto.Telephone, *createSellerDto.LocalityId)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateSeller", err, httpStatus)
		return nil
	}

	utils.Log(r, "CreateSeller", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": []models.Seller{newSeller},
	})
	return nil
}

func (h *SellerHandler) PatchSeller(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "PatchSeller", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSeller", err, httpStatus)
		return nil
	}

	var patchSellerDto dto.PatchSellerDto
	patchSellerDto, err = utils.InstantiateVarFromBody(&r.Body, patchSellerDto)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSeller", err, httpStatus)
		return nil
	}

	seller, err := h.service.UpdateSellerById(id, patchSellerDto.CompanyId, patchSellerDto.CompanyName, patchSellerDto.Address, patchSellerDto.Telephone)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSeller", err, httpStatus)
		return nil
	}

	utils.Log(r, "PatchSeller", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": []models.Seller{seller},
	})
	return nil
}

func (h *SellerHandler) DeleteSeller(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "DeleteSeller", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteSeller", err, httpStatus)
		return nil
	}

	err = h.service.DeleteSellerById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteSeller", err, httpStatus)
		return nil
	}

	utils.Log(r, "DeleteSeller", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, nil)
	return nil
}
