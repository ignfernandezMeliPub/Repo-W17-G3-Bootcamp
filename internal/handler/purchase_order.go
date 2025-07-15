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
	var purchaseOrderRequest models.PurchaseOrderCreateRequest

	purchaseOrderRequest, err := utils.InstantiateVarFromBody(&r.Body, purchaseOrderRequest)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	purchaseOrderFromRequest := purchaseOrderRequest.ToPurchaseOrder()

	newPurchaseOrder, err := h.sv.CreatePurchaseOrder(purchaseOrderFromRequest)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": newPurchaseOrder,
	})
}
