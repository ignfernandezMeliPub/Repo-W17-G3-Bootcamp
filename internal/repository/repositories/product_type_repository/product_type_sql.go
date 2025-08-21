package product_type_repository

import (
	"database/sql"
	"strconv"

	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
)

type ProductTypeRepositoryMySQL struct {
	db *sql.DB
}

func NewProductTypeRepositoryMySQL(db *sql.DB) *ProductTypeRepositoryMySQL {
	return &ProductTypeRepositoryMySQL{db: db}
}

func (r *ProductTypeRepositoryMySQL) GetProductTypeById(id int) (models.ProductType, error) {
	sql_utils.Log("GetProductTypeById", logger.LogStatusInProgress, "Select product type by id "+strconv.Itoa(id))

	productType, err := sql_utils.QueryRow[models.ProductType](r.db,
		`SELECT id, name FROM product_types WHERE id = ?`, []any{id})
	if err != nil {
		sql_utils.LogError("GetProductTypeById", "Select product type by id "+strconv.Itoa(id), err)
		return models.ProductType{}, err
	}

	sql_utils.Log("GetProductTypeById", logger.LogStatusSuccess, "Select product type by id "+strconv.Itoa(id))
	return productType, nil
}
