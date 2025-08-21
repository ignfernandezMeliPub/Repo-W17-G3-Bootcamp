package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
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

func (h *PurchaseOrderDefault) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreatePurchaseOrder", logger.LogStatusInProgress)

	var purchaseOrderRequest models.PurchaseOrderCreateRequest

	purchaseOrderRequest, err = utils.InstantiateVarFromBody(&r.Body, purchaseOrderRequest)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreatePurchaseOrder", err, httpStatus)
		return nil
	}

	purchaseOrderFromRequest := purchaseOrderRequest.ToPurchaseOrder()

	newPurchaseOrder, err := h.sv.CreatePurchaseOrder(purchaseOrderFromRequest)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreatePurchaseOrder", err, httpStatus)
		return nil
	}

	utils.Log(r, "CreatePurchaseOrder", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": newPurchaseOrder,
	})
	return nil
}
