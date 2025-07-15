package service

import (
	"app/internal/repository/inbound_order_repository"
	"app/pkg/models"
)

type InboundOrderServiceInterface interface {
	CreateInboundOrder(details models.InboundOrderRequestBody) (inboundOrder models.InboundOrder, err error)
}

func NewInboundOrderService(repository inbound_order_repository.InboundOrderRepository) *InboundOrderService {
	return &InboundOrderService{repository: repository}
}

type InboundOrderService struct {
	repository inbound_order_repository.InboundOrderRepository
}

func (s *InboundOrderService) CreateInboundOrder(details models.InboundOrderRequestBody) (inboundOrder models.InboundOrder, err error) {

	newInboundOrderDetails := models.InboundOrderDetails{
		OrderDate:      *details.Data.OrderDate,
		OrderNumber:    *details.Data.OrderNumber,
		EmployeeId:     *details.Data.EmployeeId,
		ProductBatchId: *details.Data.ProductBatchId,
		WarehouseId:    *details.Data.WarehouseId,
	}

	inboundOrder, err = s.repository.CreateInboundOrder(newInboundOrderDetails)
	return

}
