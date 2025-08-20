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

func NewBuyerDefault(sv service.BuyerService) *BuyerDefault {
	return &BuyerDefault{sv: sv}
}

type BuyerDefault struct {
	sv service.BuyerService
}

func (h *BuyerDefault) GetAllBuyers(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "GetAllBuyers", logger.LogStatusInProgress)

	b, err := h.sv.GetAllBuyers()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllBuyers", err, httpStatus)
		return
	}

	utils.Log(r, "GetAllBuyers", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}

func (h *BuyerDefault) GetBuyerById(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "GetBuyerById", logger.LogStatusInProgress)
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetBuyerById", err, httpStatus)
		return
	}

	b, err := h.sv.GetBuyerById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetBuyerById", err, httpStatus)
		return
	}

	utils.Log(r, "GetBuyerById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}

func (h *BuyerDefault) CreateBuyer(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "CreateBuyer", logger.LogStatusInProgress)
	var buyerRequest models.BuyerCreateRequest

	buyerRequest, err := utils.InstantiateVarFromBody(&r.Body, buyerRequest)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateBuyer", err, httpStatus)
		return
	}

	buyerFromRequest := buyerRequest.ToBuyer()

	newBuyer, err := h.sv.CreateBuyer(buyerFromRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateBuyer", err, httpStatus)
		return
	}

	utils.Log(r, "CreateBuyer", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": newBuyer,
	})
}

func (h *BuyerDefault) PatchBuyer(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "PatchBuyer", logger.LogStatusInProgress)
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchBuyer", err, httpStatus)
		return
	}

	var buyerPatch models.BuyerPatch

	buyerPatch, err = utils.InstantiateVarFromBody(&r.Body, buyerPatch)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchBuyer", err, httpStatus)
		return
	}

	updatedBuyer, err := h.sv.UpdateBuyerById(id, buyerPatch)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchBuyer", err, httpStatus)
		return
	}

	utils.Log(r, "PatchBuyer", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": updatedBuyer,
	})
}

func (h *BuyerDefault) DeleteBuyer(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "DeleteBuyer", logger.LogStatusInProgress)
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteBuyer", err, httpStatus)
		return
	}

	err = h.sv.DeleteBuyerById(id)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteBuyer", err, httpStatus)
		return
	}

	utils.Log(r, "DeleteBuyer", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *BuyerDefault) GetBuyersPurchaseOrdersCount(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "GetBuyersPurchaseOrdersCount", logger.LogStatusInProgress)
	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetBuyersPurchaseOrdersCount", err, httpStatus)
		return
	}

	b, err := h.sv.GetBuyersPurchaseOrdersCount(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetBuyersPurchaseOrdersCount", err, httpStatus)
		return
	}

	utils.Log(r, "GetBuyersPurchaseOrdersCount", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}
