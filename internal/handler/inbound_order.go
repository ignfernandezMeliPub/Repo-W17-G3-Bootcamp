package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type InboundOrderController struct {
	svInboundOrder service.InboundOrderServiceInterface
}

func NewInboundOrderController(svInboundOrder service.InboundOrderServiceInterface) *InboundOrderController {
	return &InboundOrderController{svInboundOrder: svInboundOrder}
}

func (c *InboundOrderController) CreateInboundOrder(w http.ResponseWriter, r *http.Request) {
	utils.Log(r, "CreateInboundOrder", logger.LogStatusInProgress)

	var newInboundOrder models.InboundOrderRequestBody

	newInboundOrder, err := utils.InstantiateVarFromBody(&r.Body, newInboundOrder)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateInboundOrder", err, httpStatus)
		return
	}

	inboundOrder, err := c.svInboundOrder.CreateInboundOrder(newInboundOrder)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateInboundOrder", err, httpStatus)
		return
	}

	utils.Log(r, "CreateInboundOrder", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": inboundOrder,
	})

}
