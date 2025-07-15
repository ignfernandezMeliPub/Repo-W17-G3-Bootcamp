package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

func NewPurchaseOrderDefault(sv service.PurchaseOrderService) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{sv: sv}
}

type PurchaseOrderDefault struct {
	sv service.PurchaseOrderService
}

func (h *PurchaseOrderDefault) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	// request
	var _PurchaseOrderRequest models.PurchaseOrderCreateRequest

	PurchaseOrderRequest, err := utils.InstantiateVarFromBody(&r.Body, _PurchaseOrderRequest)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// process
	_p := PurchaseOrderRequest.ToPurchaseOrder()

	p, err := h.sv.CreatePurchaseOrder(_p)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	// response
	data := p
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
}
