package inbound_order_repository

import "app/pkg/models"

type InboundOrderRepository interface {
	CreateInboundOrder(details models.InboundOrderDetails) (inboundOrder models.InboundOrder, err error)
}
