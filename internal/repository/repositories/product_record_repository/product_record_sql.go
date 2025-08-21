package product_record_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

type ProductRecordRepositorySQL struct {
	db *sql.DB
}

func NewProductRecordRepositorySQL(db *sql.DB) *ProductRecordRepositorySQL {
	return &ProductRecordRepositorySQL{db: db}
}

func (r *ProductRecordRepositorySQL) GetAllProductRecords() ([]models.ProductRecord, error) {
	sql_utils.Log("GetAllProductRecords", logger.LogStatusInProgress, "Select all product records")

	productRecords, err := sql_utils.Query[models.ProductRecord](r.db, "SELECT id, last_update_date, purchase_price, sale_price, product_id FROM product_records", []any{})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetAllProductRecords", "Select all product records", err)
		return productRecords, err
	}

	sql_utils.Log("GetAllProductRecords", logger.LogStatusSuccess, "Select all product records")
	return productRecords, nil
}

func (r *ProductRecordRepositorySQL) CreateProductRecord(productRecord models.ProductRecord) (models.ProductRecord, error) {
	sql_utils.LogAudit("CreateProductRecord", logger.LogStatusInProgress, "Insert product record")

	query := "INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)"
	args := []any{productRecord.LastUpdateDate, productRecord.PurchasePrice, productRecord.SalePrice, productRecord.ProductID}

	id, err := sql_utils.Insert(r.db, query, args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateProductRecord", "Insert product record", err)
		return models.ProductRecord{}, err
	}

	productRecord.ID = int(id)

	sql_utils.LogAudit("CreateProductRecord", logger.LogStatusSuccess, "Insert product record. Id: "+strconv.Itoa(productRecord.ID))
	return productRecord, nil
}
