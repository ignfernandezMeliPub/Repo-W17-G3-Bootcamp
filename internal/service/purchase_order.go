package service

import (
	"app/internal/repository/purchase_order_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"strings"
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

	if err = validatePurchaseOrderAttributes(purchaseOrder); err != nil {
		return
	}

	return s.rp.CreatePurchaseOrder(purchaseOrder)
}

func validatePurchaseOrderAttributes(purchaseOrder models.PurchaseOrder) error {

	if strings.TrimSpace(purchaseOrder.OrderNumber) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "order_number",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if purchaseOrder.OrderDate.IsZero() {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "order_date",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if strings.TrimSpace(purchaseOrder.TrackingCode) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "tracking_code",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if purchaseOrder.BuyerId == 0 {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "buyer_id",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if len(purchaseOrder.PurchaseOrderDetails) == 0 {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "purchase_order_details",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	for _, detail := range purchaseOrder.PurchaseOrderDetails {
		if detail.ProductRecordId <= 0 {
			return &custom_errors.InvalidArgValueErr{
				Argument:  "product_record_id",
				Value:     detail.ProductRecordId,
				ExtraInfo: "Value must be greater than 0",
			}
		}
		if detail.Quantity <= 0 {
			return &custom_errors.InvalidArgValueErr{
				Argument:  "quantity",
				Value:     detail.Quantity,
				ExtraInfo: "Value must be greater than 0",
			}
		}
	}

	return nil
}
