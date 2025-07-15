package purchase_order_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
)

type PurchaseOrderRepositorySQL struct {
	db *sql.DB
}

func NewPurchaseOrderRepositorySQL(db *sql.DB) *PurchaseOrderRepositorySQL {
	return &PurchaseOrderRepositorySQL{db: db}
}

func (r *PurchaseOrderRepositorySQL) CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	orderId, err := sql_utils.Insert(tx,
		"INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)",
		[]any{purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId})
	if err != nil {
		tx.Rollback()
		return
	}

	for i, detail := range purchaseOrder.PurchaseOrderDetails {
		detailId, err := sql_utils.Insert(tx,
			"INSERT INTO purchase_order_details (order_id, product_record_id, quantity) VALUES (?, ?, ?)",
			[]any{orderId, detail.ProductRecordId, detail.Quantity})
		if err != nil {
			tx.Rollback()
			return p, err
		}
		purchaseOrder.PurchaseOrderDetails[i].Id = int(detailId)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	// Set the returned purchase order with the inserted ID
	p = purchaseOrder
	p.Id = int(orderId)

	return
}
