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
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepositoryMySQL) GetProductById(id int) (models.Product, error) {
	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE id = ?`, []any{id})
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepositoryMySQL) GetProductByCode(code string) (models.Product, error) {
	product, err := sql_utils.QueryRow[models.Product](r.db,
		`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id 
	FROM products WHERE product_code = ?`, []any{code})
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepositoryMySQL) CreateProduct(p models.Product) (models.Product, error) {
	id, err := sql_utils.Insert(r.db,
		`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId})
	if err != nil {
		return models.Product{}, err
	}
	p.ID = int(id)
	return p, nil
}

func (r *ProductRepositoryMySQL) UpdateProductById(p models.Product) (models.Product, error) {
	_, err := sql_utils.Update(r.db,
		`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`, []any{p.ProductCode, p.Description, p.Width, p.Height, p.Length, p.NetWeight, p.ExpirationRate, p.RecommendedFreezingTemperature, p.FreezingRate, p.ProductTypeId, p.SellerId, p.ID})
	if err != nil {
		return models.Product{}, err
	}

	return p, nil
}

func (r *ProductRepositoryMySQL) DeleteProductById(id int) error {
	_, err := sql_utils.Delete(r.db,
		`DELETE FROM products WHERE id = ?`, []any{id})
	if err != nil {
		return err
	}
	return nil
}
