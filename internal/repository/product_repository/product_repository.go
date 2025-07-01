package product_repository

import "app/pkg/models"

type ProductRepository interface {
	FindAllProducts() ([]models.Product, error)
	FindProductById(id int) (models.Product, error)
	DeleteProduct(id int) error
	SaveProduct(p models.Product) (models.Product, error)
	FindProductByCode(code string) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
}
