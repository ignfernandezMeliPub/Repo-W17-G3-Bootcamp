package product_record_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
)

type ProductRecordRepositorySQL struct {
	db *sql.DB
}

func NewProductRecordRepositorySQL(db *sql.DB) *ProductRecordRepositorySQL {
	return &ProductRecordRepositorySQL{db: db}
}

func (r *ProductRecordRepositorySQL) GetAllProductRecords() ([]models.ProductRecord, error) {
	productRecords, err := sql_utils.Query[models.ProductRecord](r.db, "SELECT id, last_update_date, purchase_price, sale_price, product_id FROM product_records", []any{})
	return productRecords, sql_utils.HandleSqlError(err)
}

func (r *ProductRecordRepositorySQL) CreateProductRecord(productRecord models.ProductRecord) (models.ProductRecord, error) {
	query := "INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)"
	args := []any{productRecord.LastUpdateDate, productRecord.PurchasePrice, productRecord.SalePrice, productRecord.ProductID}

	id, err := sql_utils.Insert(r.db, query, args)
	if err != nil {
		return models.ProductRecord{}, sql_utils.HandleSqlError(err)
	}

	productRecord.ID = int(id)
	return productRecord, nil
}
