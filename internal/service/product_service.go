package service

import (
	"app/internal/repository/product_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"fmt"
)

type ProductServiceI interface {
	GetAllProducts() ([]models.Product, error)
	GetProductById(int) (models.Product, error)
	CreateProduct(models.ProductRequest) (models.Product, error)
	UpdateProduct(models.ProductRequest) (models.Product, error)
	DeleteProduct(int) error
}

type ProductService struct {
	ProductRepo        product_repository.ProductRepository
	ProductTypeService ProductTypeServices
	SellerServices     SellerService
}

func NewProductService(prodRepository product_repository.ProductRepository, typeService ProductTypeServices, sellService SellerService) ProductService {
	return ProductService{ProductRepo: prodRepository, ProductTypeService: typeService, SellerServices: sellService}
}

func (p *ProductService) GetAllProducts() ([]models.Product, error) {
	return p.ProductRepo.FindAllProducts()
}

func (p *ProductService) GetProductById(id int) (models.Product, error) {

	return p.ProductRepo.FindProductById(id)
}

func (p *ProductService) DeleteProduct(id int) error {
	return p.ProductRepo.DeleteProduct(id)
}

func (p *ProductService) CreateProduct(product models.ProductRequest) (models.Product, error) {
	//validar que productType exista
	if !p.isValidateProductType(*product.ProductTypeId) {
		return models.Product{}, &custom_errors.ResourceNotFoundError{}
	}
	//validar que el producto con productCode no exista
	if !p.isValidateProductCode(*product.ProductCode) {
		return models.Product{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "product_code", Value: *product.ProductCode}
	}
	//validar que el seller exista
	if !p.isValidSeller(product.SellerId) {
		return models.Product{}, &custom_errors.ResourceNotFoundError{}
	}

	newProduct := models.Product{
		ProductCode:                    *product.ProductCode,
		Description:                    *product.Description,
		Width:                          *product.Width,
		Height:                         *product.Height,
		Length:                         *product.Length,
		NetWeight:                      *product.NetWeight,
		ExpirationRate:                 *product.ExpirationRate,
		RecommendedFreezingTemperature: *product.RecommendedFreezingTemperature,
		FreezingRate:                   *product.FreezingRate,
		ProductTypeId:                  *product.ProductTypeId,
		SellerId:                       product.SellerId,
	}

	newProduct, err := p.ProductRepo.SaveProduct(newProduct)

	if err != nil {
		return newProduct, err
	}
	return newProduct, nil

}

func (p *ProductService) UpdateProduct(updateProduct models.ProductRequest) (models.Product, error) {
	//validar datos de negocio
	product, err := p.ProductRepo.FindProductById(updateProduct.Id)

	if err != nil {
		return models.Product{}, err
	}

	updatedProduct, err := p.patchProduct(product, updateProduct)
	if err != nil {
		return models.Product{}, err
	}

	return p.ProductRepo.UpdateProduct(updatedProduct)
}

func (p *ProductService) patchProduct(product models.Product, updateProduct models.ProductRequest) (models.Product, error) {
	if updateProduct.ProductCode != nil {
		if *updateProduct.ProductCode == "" {
			return models.Product{}, fmt.Errorf("invalid product code #%s", *updateProduct.ProductCode)
		}

		if !p.isValidateProductCode(*updateProduct.ProductCode) {
			return models.Product{}, fmt.Errorf("product code #%s already exists", *updateProduct.ProductCode)
		}
		product.ProductCode = *updateProduct.ProductCode
	}

	if updateProduct.Description != nil {
		if *updateProduct.Description == "" {
			return models.Product{}, fmt.Errorf("invalid description code #%s", *updateProduct.Description)
		}
		product.Description = *updateProduct.Description
	}
	if updateProduct.Width != nil {
		if *updateProduct.Width <= 0.0 {
			return models.Product{}, fmt.Errorf("invalid width code #%.2f", *updateProduct.Width)
		}
		product.Width = *updateProduct.Width
	}
	if updateProduct.Height != nil {
		if *updateProduct.Height <= 0.0 {
			return models.Product{}, fmt.Errorf("invalid height code #%.2f", *updateProduct.Height)
		}
		product.Height = *updateProduct.Height
	}
	if updateProduct.Length != nil {
		if *updateProduct.Length <= 0.0 {
			return models.Product{}, fmt.Errorf("invalid lenght code #%.2f", *updateProduct.Length)
		}
		product.Length = *updateProduct.Length
	}
	if updateProduct.NetWeight != nil {
		if *updateProduct.NetWeight <= 0.0 {
			return models.Product{}, fmt.Errorf("invalid net_weight code #%.2f", *updateProduct.NetWeight)
		}
		product.NetWeight = *updateProduct.NetWeight
	}
	if updateProduct.ExpirationRate != nil {
		if *updateProduct.ExpirationRate <= 0 {
			return models.Product{}, fmt.Errorf("invalid expiration_rate code #%d", *updateProduct.ExpirationRate)
		}
		product.ExpirationRate = *updateProduct.ExpirationRate
	}
	if updateProduct.RecommendedFreezingTemperature != nil {
		product.RecommendedFreezingTemperature = *updateProduct.RecommendedFreezingTemperature
	}
	if updateProduct.FreezingRate != nil {
		if *updateProduct.FreezingRate <= 0 {
			return models.Product{}, fmt.Errorf("invalid freezing_rate code #%d", *updateProduct.FreezingRate)
		}
		product.FreezingRate = *updateProduct.FreezingRate
	}
	if updateProduct.ProductTypeId != nil {
		if !p.isValidateProductType(*updateProduct.ProductTypeId) {
			return models.Product{}, fmt.Errorf("product type with id %d not exist", *updateProduct.ProductTypeId)
		}
		product.ProductTypeId = *updateProduct.ProductTypeId
	}
	if updateProduct.SellerId != 0 {
		//aca se valida si el seller existe
		product.SellerId = updateProduct.SellerId
	}

	return product, nil

}

func (p *ProductService) isValidateProductType(id int) bool {
	return p.ProductTypeService.IsValidProductType(id)
}

func (p *ProductService) isValidateProductCode(code string) bool {
	_, err := p.ProductRepo.FindProductByCode(code)

	return err != nil
}

func (p *ProductService) isValidSeller(id int) bool {
	if id == 0 {
		return true
	}
	_, err := p.SellerServices.GetById(id)

	return err == nil
}
