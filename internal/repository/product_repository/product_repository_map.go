package product_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type ProductRepositoryMap struct {
	database map[int]models.Product
	seed     int
}

func NewProductRepositoryMap(database map[int]models.Product) *ProductRepositoryMap {

	if database == nil {
		database = map[int]models.Product{}
	}

	maxId := 0
	for key := range database {
		if key > maxId {
			maxId = key
		}
	}
	return &ProductRepositoryMap{database: database, seed: maxId + 1}

}

func (r *ProductRepositoryMap) FindAllProducts() ([]models.Product, error) { //cambia a array
	products := make([]models.Product, 0)
	for _, prod := range r.database {
		products = append(products, prod)
	}
	return products, nil
}

func (r *ProductRepositoryMap) FindProductById(id int) (models.Product, error) {
	product, found := r.database[id]
	if !found {
		return models.Product{}, &custom_errors.ResourceNotFoundError{}
	}
	return product, nil
}

func (r *ProductRepositoryMap) DeleteProduct(id int) error {
	_, found := r.database[id]
	if !found {
		return &custom_errors.ResourceNotFoundError{}
	}
	delete(r.database, id)
	return nil
}

func (r *ProductRepositoryMap) SaveProduct(p models.Product) (models.Product, error) {
	p.ID = r.seed
	r.database[p.ID] = p
	r.seed++
	return p, nil
}

func (r *ProductRepositoryMap) FindProductByCode(code string) (models.Product, error) {
	for _, product := range r.database {
		if product.ProductCode == code {
			return product, nil
		}
	}
	return models.Product{}, &custom_errors.ResourceNotFoundError{}
}

func (r *ProductRepositoryMap) UpdateProduct(p models.Product) (models.Product, error) {
	r.database[p.ID] = p
	return p, nil
}
