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
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		// response
		data := make(map[int]models.Buyer)
		for _, buyer := range b {
			data[buyer.Id] = buyer
		}
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
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}
		// process

		b, err := h.sv.FindBuyerByID(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
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

		// process

		// response
	}
}

func (h *BuyerDefault) PatchBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process

		// response
	}
}

func (h *BuyerDefault) DeleteBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process

		// response
	}
}
