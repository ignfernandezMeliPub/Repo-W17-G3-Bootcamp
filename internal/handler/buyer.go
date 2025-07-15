package handler

import (
	"app/internal/handler/utils"
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

func (h *BuyerDefault) GetAllBuyers(w http.ResponseWriter, _ *http.Request) {
	// request

	// process
	// process
	b, err := h.sv.GetAllBuyers()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// response
	data := b

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *BuyerDefault) GetBuyerById(w http.ResponseWriter, r *http.Request) {
	// request
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	// process

	b, err := h.sv.GetBuyerById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	// response

	data := b
	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *BuyerDefault) CreateBuyer(w http.ResponseWriter, r *http.Request) {
	// request
	var _buyerRequest models.BuyerCreateRequest

	buyerRequest, err := utils.InstantiateVarFromBody(&r.Body, _buyerRequest)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// process

	_b := buyerRequest.ToBuyer()

	b, err := h.sv.CreateBuyer(_b)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// response
	data := b
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
}

func (h *BuyerDefault) PatchBuyer(w http.ResponseWriter, r *http.Request) {
	// request
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	var _buyerPatch models.BuyerPatch

	buyerPatch, err := utils.InstantiateVarFromBody(&r.Body, _buyerPatch)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// process
	b, err := h.sv.UpdateBuyerById(id, buyerPatch)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// response
	data := b
	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}

func (h *BuyerDefault) DeleteBuyer(w http.ResponseWriter, r *http.Request) {
	// request
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	// process

	err = h.sv.DeleteBuyerById(id)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	// response

	response.JSON(w, http.StatusNoContent, map[string]any{
		"data": nil,
	})
}

func (h *BuyerDefault) GetBuyersPurchaseOrdersCount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	if id != nil {
		// Purchase orders by buyer id
		data, err := h.sv.GetBuyerPurchaseOrdersCount(*id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
		return
	}

	// Purchase orders qty by buyer
	data, err := h.sv.GetBuyersPurchaseOrdersCount()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
}
