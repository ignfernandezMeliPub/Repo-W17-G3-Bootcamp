package inbound_order_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

func NewInboundOrderDb(db *sql.DB) InboundOrderRepository {

	return &InboundOrderDb{db: db}

}

type InboundOrderDb struct {
	db *sql.DB
}

func (r *InboundOrderDb) CreateInboundOrder(details models.InboundOrderDetails) (inboundOrder models.InboundOrder, err error) {
	sql_utils.LogAudit("CreateInboundOrder", logger.LogStatusInProgress, "Insert inbound order")

	args := []any{details.OrderDate, details.OrderNumber, details.EmployeeId, details.ProductBatchId, details.WarehouseId}

	lastId, err := sql_utils.Insert(r.db, "INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateInboundOrder", "Insert inbound order", err)
		return models.InboundOrder{}, err
	}

	inboundOrder = models.InboundOrder{

		Id:                  int(lastId),
		InboundOrderDetails: details,
	}

	sql_utils.LogAudit("CreateInboundOrder", logger.LogStatusSuccess, "Insert inbound order. Id: "+strconv.Itoa(inboundOrder.Id))
	return inboundOrder, nil

}
