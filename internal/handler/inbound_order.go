package handler

import (
	"app/internal/handler/utils"
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

	var newInboundOrder models.InboundOrderRequestBody

	newInboundOrder, err := utils.InstantiateVarFromBody(&r.Body, newInboundOrder)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	inboundOrder, err := c.svInboundOrder.CreateInboundOrder(newInboundOrder)

	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	response.JSON(w, http.StatusContinue, map[string]any{
		"data": inboundOrder,
	})

}
