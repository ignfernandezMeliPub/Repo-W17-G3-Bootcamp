package product_repository

import (
	"database/sql"
	"strconv"

	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type ProductRepositoryMySQL struct {
	db *sql.DB
}

func NewProductRepositoryMySQL(db *sql.DB) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{db: db}
}

func (r *ProductRepositoryMySQL) GetAllProducts() ([]models.Product, error) {
	sql_utils.Log("GetAllProducts", logger.LogStatusInProgress, "Select all products")

	products, err := sql_utils.Query[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products`, []any{})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetAllProducts", "Select all products", err)
		return products, err
	}

	sql_utils.Log("GetAllProducts", logger.LogStatusSuccess, "Select all products")
	return products, nil
}

func (r *ProductRepositoryMySQL) GetProductById(id int) (models.Product, error) {
	sql_utils.Log("GetProductById", logger.LogStatusInProgress, "Select product by id "+strconv.Itoa(id))

	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE id = ?`, []any{id})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetProductById", "Select product by id "+strconv.Itoa(id), err)
		return product, err
	}

	sql_utils.Log("GetProductById", logger.LogStatusSuccess, "Select product by id "+strconv.Itoa(id))
	return product, nil
}

func (r *ProductRepositoryMySQL) GetProductByCode(code string) (models.Product, error) {
	sql_utils.Log("GetProductByCode", logger.LogStatusInProgress, "Select product by code "+code)

	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE product_code = ?`, []any{code})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetProductByCode", "Select product by code "+code, err)
		return product, err
	}

	sql_utils.Log("GetProductByCode", logger.LogStatusSuccess, "Select product by code "+code)
	return product, nil
}

func (r *ProductRepositoryMySQL) CreateProduct(p models.Product) (models.Product, error) {
	sql_utils.LogAudit("CreateProduct", logger.LogStatusInProgress, "Insert product")

	id, err := sql_utils.Insert(r.db,
		`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateProduct", "Insert product", err)
		return models.Product{}, err
	}
	p.ID = int(id)

	sql_utils.LogAudit("CreateProduct", logger.LogStatusSuccess, "Insert product. Id: "+strconv.Itoa(p.ID))
	return p, nil
}

func (r *ProductRepositoryMySQL) UpdateProductById(p models.Product) (models.Product, error) {
	sql_utils.LogAudit("UpdateProductById", logger.LogStatusInProgress, "Update product by id: "+strconv.Itoa(p.ID))

	_, err := sql_utils.Update(r.db,
		`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId, p.ID})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogAuditError("UpdateProductById", "Update product by id: "+strconv.Itoa(p.ID), err)
		return p, err
	}

	sql_utils.LogAudit("UpdateProductById", logger.LogStatusSuccess, "Update product by id: "+strconv.Itoa(p.ID))
	return p, nil
}

func (r *ProductRepositoryMySQL) DeleteProductById(id int) error {
	sql_utils.LogAudit("DeleteProductById", logger.LogStatusInProgress, "Delete product by id: "+strconv.Itoa(id))

	rowsAffected, err := sql_utils.Delete(r.db,
		`DELETE FROM products WHERE id = ?`, []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteProductById", "Delete product by id: "+strconv.Itoa(id), err)
		return err
	}
	if rowsAffected == 0 {
		err = custom_errors.ErrNotFound
		sql_utils.LogAuditError("DeleteProductById", "Delete product by id: "+strconv.Itoa(id), err)
		return err
	}

	sql_utils.LogAudit("DeleteProductById", logger.LogStatusSuccess, "Delete product by id: "+strconv.Itoa(id))
	return nil
}

func (r *ProductRepositoryMySQL) GetReportRecords(id *int) ([]models.ProductRecordReport, error) {
	if id != nil {
		sql_utils.Log("GetReportRecords", logger.LogStatusInProgress, "Select product record report by id: "+strconv.Itoa(*id))
	} else {
		sql_utils.Log("GetReportRecords", logger.LogStatusInProgress, "Select product record report")
	}

	query := `SELECT p.id product_id, p.description, COUNT(r.id) as records_count 
		FROM products p 
		LEFT JOIN product_records r ON r.product_id = p.id`
	var args []any

	if id != nil {
		query += " WHERE p.id = ?"
		args = append(args, *id)
	}

	query += " GROUP BY p.id, p.description"

	ProductRecordReport, err := sql_utils.Query[models.ProductRecordReport](r.db, query, args)
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetReportRecords", "Select product record report", err)
		return ProductRecordReport, err
	}

	if id != nil {
		sql_utils.Log("GetReportRecords", logger.LogStatusSuccess, "Select product record report by id: "+strconv.Itoa(*id))
	} else {
		sql_utils.Log("GetReportRecords", logger.LogStatusSuccess, "Select product record report")
	}
	return ProductRecordReport, nil
}
