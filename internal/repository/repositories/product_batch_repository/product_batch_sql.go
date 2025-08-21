package product_batch_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

type ProductBatchRepositorySQL struct {
	db *sql.DB
}

func NewProductBatchRepositorySQL(db *sql.DB) *ProductBatchRepositorySQL {
	return &ProductBatchRepositorySQL{db}
}

func (p *ProductBatchRepositorySQL) CreateProductBatch(pb models.ProductBatch) (prod models.ProductBatch, err error) {
	sql_utils.LogAudit("CreateProductBatch", logger.LogStatusInProgress, "Insert product batch")

	args := []any{pb.BatchNumber, pb.CurrentQuantity, pb.CurrentTemperature, pb.DueDate, pb.InitialQuantity, pb.ManufacturingDate, pb.ManufacturingHour, pb.MinimumTemperature, pb.ProductId, pb.SectionId}
	lastId, err := sql_utils.Insert(p.db,
		"INSERT INTO `product_batches` (`batch_number`,`current_quantity`,`current_temperature`,`due_date`,`initial_quantity`,`manufacturing_date`,`manufacturing_hour`,`minimum_temperature`,`product_id`,`section_id`) "+
			"VALUES (?,?,?,?,?,?,?,?,?,?)", args,
	)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateProductBatch", "Insert product batch", err)
		return models.ProductBatch{}, err
	}
	pb.ID = int(lastId)

	sql_utils.LogAudit("CreateProductBatch", logger.LogStatusSuccess, "Insert product batch. Id: "+strconv.Itoa(pb.ID))
	return pb, nil
}
