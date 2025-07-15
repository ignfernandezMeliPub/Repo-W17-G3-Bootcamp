package service

import (
	"app/internal/repository/repositories/purchase_order_repository"
	"app/pkg/models"
)

type PurchaseOrderService interface {
	CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error)
}

type PurchaseOrderServiceDefault struct {
	rp purchase_order_repository.PurchaseOrderRepository
}

func NewPurchaseOrderDefault(rp purchase_order_repository.PurchaseOrderRepository) *PurchaseOrderServiceDefault {
	return &PurchaseOrderServiceDefault{rp: rp}
}

func (s *PurchaseOrderServiceDefault) CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error) {
	return s.rp.CreatePurchaseOrder(purchaseOrder)
}
