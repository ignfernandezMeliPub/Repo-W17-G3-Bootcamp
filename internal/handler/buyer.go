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
	b, err := h.sv.GetAllBuyers()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}

func (h *BuyerDefault) GetBuyerById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	b, err := h.sv.GetBuyerById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}

func (h *BuyerDefault) CreateBuyer(w http.ResponseWriter, r *http.Request) {
	var buyerRequest models.BuyerCreateRequest

	buyerRequest, err := utils.InstantiateVarFromBody(&r.Body, buyerRequest)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	buyerFromRequest := buyerRequest.ToBuyer()

	newBuyer, err := h.sv.CreateBuyer(buyerFromRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": newBuyer,
	})
}

func (h *BuyerDefault) PatchBuyer(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	var buyerPatch models.BuyerPatch

	buyerPatch, err = utils.InstantiateVarFromBody(&r.Body, buyerPatch)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	updatedBuyer, err := h.sv.UpdateBuyerById(id, buyerPatch)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": updatedBuyer,
	})
}

func (h *BuyerDefault) DeleteBuyer(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	err = h.sv.DeleteBuyerById(id)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *BuyerDefault) GetBuyersPurchaseOrdersCount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	if id != nil {
		// Purchase orders qty by buyer id
		b, err := h.sv.GetBuyerPurchaseOrdersCount(*id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": b,
		})
		return
	}

	// Purchase orders qty by every buyer
	b, err := h.sv.GetBuyersPurchaseOrdersCount()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": b,
	})
}
