package purchase_order_repository

import "app/pkg/models"

type PurchaseOrderRepository interface {
	CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error)
}
