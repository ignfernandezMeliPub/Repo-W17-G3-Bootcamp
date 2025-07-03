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
	UpdateProductById(models.ProductPatchRequest) (models.Product, error)
	DeleteProductById(int) error
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
	return p.ProductRepo.GetAllProducts()
}

func (p *ProductService) GetProductById(id int) (models.Product, error) {

	return p.ProductRepo.GetProductById(id)
}

func (p *ProductService) DeleteProductById(id int) error {
	return p.ProductRepo.DeleteProductById(id)
}

func (p *ProductService) CreateProduct(product models.ProductRequest) (models.Product, error) {
	// validar que productType exista
	if !p.isValidateProductType(*product.ProductTypeId) {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "product_type_id", Value: fmt.Sprint(product.SellerId), ExtraInfo: "Invalid Product Type ID"}
	}

	// validar que el producto con productCode no exista
	if !p.isValidateProductCode(*product.ProductCode) {
		return models.Product{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "product_code", Value: *product.ProductCode}
	}

	// validar que el seller exista
	if !p.isValidSeller(product.SellerId) {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "seller_id", Value: fmt.Sprint(product.SellerId), ExtraInfo: "Invalid Seller ID"}
	}

	// Validaciones de valores positivos y mayores a 0
	if *product.Width <= 0.0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "width", Value: *product.Width, ExtraInfo: "Width must be greater than 0"}
	}
	if *product.Height <= 0.0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "height", Value: *product.Height, ExtraInfo: "Height must be greater than 0"}
	}
	if *product.Length <= 0.0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "length", Value: *product.Length, ExtraInfo: "Length must be greater than 0"}
	}
	if *product.NetWeight <= 0.0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "net_weight", Value: *product.NetWeight, ExtraInfo: "Net weight must be greater than 0"}
	}
	if *product.ExpirationRate <= 0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "expiration_rate", Value: *product.ExpirationRate, ExtraInfo: "Expiration rate must be greater than 0"}
	}
	if *product.FreezingRate <= 0 {
		return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "freezing_rate", Value: *product.FreezingRate, ExtraInfo: "Freezing rate must be greater than 0"}
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

	newProduct, err := p.ProductRepo.CreateProduct(newProduct)

	if err != nil {
		return newProduct, err
	}
	return newProduct, nil

}

func (p *ProductService) UpdateProductById(updateProduct models.ProductPatchRequest) (models.Product, error) {
	product, err := p.ProductRepo.GetProductById(updateProduct.Id)

	if err != nil {
		return models.Product{}, err
	}

	updatedProduct, err := p.patchProduct(product, updateProduct)
	if err != nil {
		return models.Product{}, err
	}

	return p.ProductRepo.UpdateProductById(updatedProduct)
}

func (p *ProductService) patchProduct(product models.Product, updateProduct models.ProductPatchRequest) (models.Product, error) {
	if updateProduct.ProductCode != nil {
		if *updateProduct.ProductCode == "" {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "product_code", Value: *updateProduct.ProductCode, ExtraInfo: "Product code cannot be empty"}
		}

		if product.ProductCode != *updateProduct.ProductCode && !p.isValidateProductCode(*updateProduct.ProductCode) {
			return models.Product{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "product_code", Value: *updateProduct.ProductCode}
		}
		product.ProductCode = *updateProduct.ProductCode
	}

	if updateProduct.Description != nil {
		if *updateProduct.Description == "" {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "description", Value: *updateProduct.Description, ExtraInfo: "Description cannot be empty"}
		}
		product.Description = *updateProduct.Description
	}
	if updateProduct.Width != nil {
		if *updateProduct.Width <= 0.0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "width", Value: *updateProduct.Width, ExtraInfo: "Width must be greater than 0"}
		}
		product.Width = *updateProduct.Width
	}
	if updateProduct.Height != nil {
		if *updateProduct.Height <= 0.0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "height", Value: *updateProduct.Height, ExtraInfo: "Height must be greater than 0"}
		}
		product.Height = *updateProduct.Height
	}
	if updateProduct.Length != nil {
		if *updateProduct.Length <= 0.0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "length", Value: *updateProduct.Length, ExtraInfo: "Length must be greater than 0"}
		}
		product.Length = *updateProduct.Length
	}
	if updateProduct.NetWeight != nil {
		if *updateProduct.NetWeight <= 0.0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "net_weight", Value: *updateProduct.NetWeight, ExtraInfo: "Net weight must be greater than 0"}
		}
		product.NetWeight = *updateProduct.NetWeight
	}
	if updateProduct.ExpirationRate != nil {
		if *updateProduct.ExpirationRate <= 0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "expiration_rate", Value: *updateProduct.ExpirationRate, ExtraInfo: "Expiration rate must be greater than 0"}
		}
		product.ExpirationRate = *updateProduct.ExpirationRate
	}
	if updateProduct.RecommendedFreezingTemperature != nil {
		product.RecommendedFreezingTemperature = *updateProduct.RecommendedFreezingTemperature
	}
	if updateProduct.FreezingRate != nil {
		if *updateProduct.FreezingRate <= 0 {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "freezing_rate", Value: *updateProduct.FreezingRate, ExtraInfo: "Freezing rate must be greater than 0"}
		}
		product.FreezingRate = *updateProduct.FreezingRate
	}
	if updateProduct.ProductTypeId != nil {
		if !p.isValidateProductType(*updateProduct.ProductTypeId) {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "product_type_id", Value: fmt.Sprint(product.SellerId), ExtraInfo: "Invalid Product Type ID"}
		}
		product.ProductTypeId = *updateProduct.ProductTypeId
	}
	if updateProduct.SellerId != nil {
		if !p.isValidSeller(*updateProduct.SellerId) {
			return models.Product{}, &custom_errors.InvalidArgValueErr{Argument: "seller_id", Value: fmt.Sprint(product.SellerId), ExtraInfo: "Invalid Seller ID"}
		}
		product.SellerId = *updateProduct.SellerId
	}

	return product, nil

}

func (p *ProductService) isValidateProductType(id int) bool {
	return p.ProductTypeService.IsValidProductType(id)
}

func (p *ProductService) isValidateProductCode(code string) bool {
	_, err := p.ProductRepo.GetProductByCode(code)

	return err != nil
}

func (p *ProductService) isValidSeller(id int) bool {
	if id == 0 {
		return true
	}
	_, err := p.SellerServices.GetSellerById(id)

	return err == nil
}
