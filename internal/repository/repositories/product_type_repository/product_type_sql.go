package product_type_repository

import (
	"database/sql"

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
	productType, err := sql_utils.QueryRow[models.ProductType](r.db,
		`SELECT id, name FROM product_types WHERE id = ?`, []any{id})
	if err != nil {
		return models.ProductType{}, err
	}
	return productType, nil
}

func (r *ProductTypeRepositoryMySQL) IsValidProductType(id int) bool {
	_, err := sql_utils.QueryRow[models.ProductType](r.db,
		`SELECT id FROM product_types WHERE id = ?`, []any{id})
	if err != nil {
		return false
	}
	return true
}
