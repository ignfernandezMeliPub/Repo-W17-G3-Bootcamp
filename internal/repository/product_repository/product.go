package product_repository

import "app/pkg/models"

type ProductRepository interface {
	CreateProduct(p models.Product) (models.Product, error)
	GetAllProducts() ([]models.Product, error)
	GetProductById(id int) (models.Product, error)
	GetProductByCode(code string) (models.Product, error)
	UpdateProductById(product models.Product) (models.Product, error)
	DeleteProductById(id int) error
}
