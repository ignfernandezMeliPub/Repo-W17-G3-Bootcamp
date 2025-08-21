package purchase_order_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

type PurchaseOrderRepositorySQL struct {
	db *sql.DB
}

func NewPurchaseOrderRepositorySQL(db *sql.DB) *PurchaseOrderRepositorySQL {
	return &PurchaseOrderRepositorySQL{db: db}
}

func (r *PurchaseOrderRepositorySQL) CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error) {
	sql_utils.LogAudit("CreatePurchaseOrder", logger.LogStatusInProgress, "Insert purchase order")

	tx, err := r.db.Begin()
	if err != nil {
		sql_utils.LogAuditError("CreatePurchaseOrder", "Insert purchase order", err)
		return
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	orderId, err := sql_utils.Insert(tx,
		"INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)",
		[]any{purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreatePurchaseOrder", "Insert purchase order", err)
		return
	}

	if len(purchaseOrder.PurchaseOrderDetails) > 0 {
		insertQuery := "INSERT INTO purchase_order_details (order_id, product_record_id, quantity) VALUES "
		values := []any{}
		for _, detail := range purchaseOrder.PurchaseOrderDetails {
			insertQuery += "(?, ?, ?), "
			values = append(values, orderId, detail.ProductRecordId, detail.Quantity)
		}

		insertQuery = insertQuery[:len(insertQuery)-2] + ";"

		_, err = sql_utils.Insert(tx, insertQuery, values)
		if err != nil {
			err = sql_utils.HandleSqlError(err)
			sql_utils.LogAuditError("CreatePurchaseOrder", "Insert purchase order", err)
			return p, err
		}

		orderDetails, err := sql_utils.Query[models.PurchaseOrderDetail](tx, "SELECT id, order_id, product_record_id, quantity FROM purchase_order_details WHERE order_id = ?", []any{orderId})
		if err != nil {
			err = sql_utils.HandleSqlError(err)
			sql_utils.LogAuditError("CreatePurchaseOrder", "Insert purchase order", err)
			return p, err
		}

		purchaseOrder.PurchaseOrderDetails = orderDetails
	}

	err = tx.Commit()
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreatePurchaseOrder", "Insert purchase order", err)
		return
	}
	committed = true

	// Set the returned purchase order with the inserted ID
	p = purchaseOrder
	p.Id = int(orderId)

	sql_utils.LogAudit("CreatePurchaseOrder", logger.LogStatusSuccess, "Insert purchase order. Id: "+strconv.Itoa(p.Id))
	return
}
