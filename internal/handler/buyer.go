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

func (h *BuyerDefault) GetAllBuyers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		// process
		b, err := h.sv.FindAllBuyers()
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		// response
		data := b

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BuyerDefault) GetBuyerByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		// process

		b, err := h.sv.FindBuyerByID(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		// response

		data := b
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BuyerDefault) CreateBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var _buyer_request models.BuyerCreateRequest

		buyer_request, err := utils.InstantiateVarFromBody(&r.Body, _buyer_request)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		// process

		_b := buyer_request.ToBuyer()

		b, err := h.sv.CreateBuyer(_b)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		// response
		data := b
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BuyerDefault) PatchBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		var _buyer_patch models.BuyerPatch

		buyer_patch, err := utils.InstantiateVarFromBody(&r.Body, _buyer_patch)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		// process
		b, err := h.sv.UpdateBuyerByID(id, buyer_patch)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		// response
		data := b
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BuyerDefault) DeleteBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		// process

		err = h.sv.DeleteBuyerByID(id)

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		// response

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "success",
			"data":    nil,
		})
	}
}
