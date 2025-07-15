package product_repository

import (
	"database/sql"

	"app/internal/repository/sql_utils"
	"app/pkg/models"
)

type ProductRepositoryMySQL struct {
	db *sql.DB
}

func NewProductRepositoryMySQL(db *sql.DB) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{db: db}
}

func (r *ProductRepositoryMySQL) GetAllProducts() ([]models.Product, error) {
	products, err := sql_utils.Query[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products`, []any{})
	return products, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) GetProductById(id int) (models.Product, error) {
	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE id = ?`, []any{id})
	return product, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) GetProductByCode(code string) (models.Product, error) {
	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE product_code = ?`, []any{code})
	return product, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) CreateProduct(p models.Product) (models.Product, error) {
	id, err := sql_utils.Insert(r.db,
		`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId})
	if err != nil {
		return models.Product{}, sql_utils.HandleSqlError(err)
	}
	p.ID = int(id)
	return p, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) UpdateProductById(p models.Product) (models.Product, error) {
	_, err := sql_utils.Update(r.db,
		`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId, p.ID})
	return p, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) DeleteProductById(id int) error {
	_, err := sql_utils.Delete(r.db,
		`DELETE FROM products WHERE id = ?`, []any{id})
	return sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) GetReportRecords(id int) ([]models.ProductRecordReport, error) {
	ProductRecordReport, err := sql_utils.Query[models.ProductRecordReport](r.db,
		`SELECT p.id product_id, p.description, COUNT(r.id) as records_count 
		FROM products p 
		LEFT JOIN product_records r ON r.product_id = p.id
		WHERE p.id = ? GROUP BY p.id, p.description`, []any{id})
	return ProductRecordReport, sql_utils.HandleSqlError(err)
}

func (r *ProductRepositoryMySQL) GetAllReportRecords() ([]models.ProductRecordReport, error) {
	ProductRecordReport, err := sql_utils.Query[models.ProductRecordReport](r.db,
		`SELECT p.id product_id, p.description, COUNT(r.id) as records_count 
		FROM products p 
		LEFT JOIN product_records r ON r.product_id = p.id
		GROUP BY p.id, p.description`, []any{})
	return ProductRecordReport, sql_utils.HandleSqlError(err)
}
