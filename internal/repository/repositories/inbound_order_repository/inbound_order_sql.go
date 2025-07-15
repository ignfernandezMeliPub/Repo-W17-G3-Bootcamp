package inbound_order_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
)

func NewInboundOrderDb(db *sql.DB) InboundOrderRepository {

	return &InboundOrderDb{db: db}

}

type InboundOrderDb struct {
	db *sql.DB
}

func (r *InboundOrderDb) CreateInboundOrder(details models.InboundOrderDetails) (inboundOrder models.InboundOrder, err error) {

	args := []any{details.OrderDate, details.OrderNumber, details.EmployeeId, details.ProductBatchId, details.WarehouseId}

	lastId, err := sql_utils.Insert(r.db, "INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)", args)

	if err != nil {

		return models.InboundOrder{}, err

	}

	inboundOrder = models.InboundOrder{

		Id:                  int(lastId),
		InboundOrderDetails: details,
	}

	return

}
