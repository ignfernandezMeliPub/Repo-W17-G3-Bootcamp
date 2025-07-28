package service

import (
	"app/pkg/models"
	"app/test/repository"
	"app/test/service"
	"errors"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProductService(t *testing.T) {
	mockProductRepository := new(repository.MockProductRepository)
	mockProductTypeService := new(service.MockProductTypeService)
	mockSellerService := new(service.MockSellerService)
	productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)
	require.NotNil(t, productService)
}

func TestCreateProduct(t *testing.T) {
	validProductCode := "1234567890"
	validDescription := "Description 1"
	validWidth := 10.0
	validHeight := 10.0
	validLength := 10.0
	validNetWeight := 10.0
	validExpirationRate := 10
	validRecommendedFreezingTemperature := 10.0
	validFreezingRate := 10
	validProductTypeId := 1
	var validSellerId *int = nil

	validProductRequest := models.ProductRequest{
		ProductCode:                    &validProductCode,
		Description:                    &validDescription,
		Width:                          &validWidth,
		Height:                         &validHeight,
		Length:                         &validLength,
		NetWeight:                      &validNetWeight,
		ExpirationRate:                 &validExpirationRate,
		RecommendedFreezingTemperature: &validRecommendedFreezingTemperature,
		FreezingRate:                   &validFreezingRate,
		ProductTypeId:                  &validProductTypeId,
		SellerId:                       validSellerId,
	}

	validProduct := models.Product{
		ID:                             1,
		ProductCode:                    validProductCode,
		Description:                    validDescription,
		Width:                          validWidth,
		Height:                         validHeight,
		Length:                         validLength,
		NetWeight:                      validNetWeight,
		ExpirationRate:                 validExpirationRate,
		RecommendedFreezingTemperature: validRecommendedFreezingTemperature,
		FreezingRate:                   validFreezingRate,
		ProductTypeId:                  validProductTypeId,
		SellerId:                       validSellerId,
	}

	t.Run("should create a product successfully", func(t *testing.T) {
		inputProduct := models.Product{
			ProductCode:                    validProductCode,
			Description:                    validDescription,
			Width:                          validWidth,
			Height:                         validHeight,
			Length:                         validLength,
			NetWeight:                      validNetWeight,
			ExpirationRate:                 validExpirationRate,
			RecommendedFreezingTemperature: validRecommendedFreezingTemperature,
			FreezingRate:                   validFreezingRate,
			ProductTypeId:                  validProductTypeId,
			SellerId:                       validSellerId,
		}
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		mockProductRepository.On("CreateProduct", inputProduct).Return(validProduct, nil)
		createdProduct, err := productService.CreateProduct(validProductRequest)
		require.NoError(t, err)
		require.Equal(t, validProduct, createdProduct)
		mockProductRepository.AssertExpectations(t)
	})

	t.Run("should return error when width is less than or equal to 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidWidth := 0.0
		invalidProductRequest := validProductRequest
		invalidProductRequest.Width = &invalidWidth

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Width must be greater than 0")
	})

	t.Run("should return error when height is less than or equal to 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidHeight := -5.0
		invalidProductRequest := validProductRequest
		invalidProductRequest.Height = &invalidHeight

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Height must be greater than 0")
	})

	t.Run("should return error when length is less than or equal to 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidLength := 0.0
		invalidProductRequest := validProductRequest
		invalidProductRequest.Length = &invalidLength

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Length must be greater than 0")
	})

	t.Run("should return error when net_weight is less than or equal to 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidNetWeight := -1.0
		invalidProductRequest := validProductRequest
		invalidProductRequest.NetWeight = &invalidNetWeight

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Net weight must be greater than 0")
	})

	t.Run("should return error when expiration_rate is less than or equal to 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidExpirationRate := 0
		invalidProductRequest := validProductRequest
		invalidProductRequest.ExpirationRate = &invalidExpirationRate

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Expiration rate must be greater than 0")
	})

	t.Run("should return error when freezing_rate is less than 0", func(t *testing.T) {
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		invalidFreezingRate := -5
		invalidProductRequest := validProductRequest
		invalidProductRequest.FreezingRate = &invalidFreezingRate

		_, err := productService.CreateProduct(invalidProductRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Freezing rate must be greater than 0")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		inputProduct := models.Product{
			ProductCode:                    validProductCode,
			Description:                    validDescription,
			Width:                          validWidth,
			Height:                         validHeight,
			Length:                         validLength,
			NetWeight:                      validNetWeight,
			ExpirationRate:                 validExpirationRate,
			RecommendedFreezingTemperature: validRecommendedFreezingTemperature,
			FreezingRate:                   validFreezingRate,
			ProductTypeId:                  validProductTypeId,
			SellerId:                       validSellerId,
		}
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		expectedError := errors.New("database error")
		mockProductRepository.On("CreateProduct", inputProduct).Return(models.Product{}, expectedError)

		_, err := productService.CreateProduct(validProductRequest)
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		mockProductRepository.AssertExpectations(t)
	})

	t.Run("should create product with seller_id", func(t *testing.T) {
		sellerId := 5

		inputProduct := models.Product{
			ProductCode:                    validProductCode,
			Description:                    validDescription,
			Width:                          validWidth,
			Height:                         validHeight,
			Length:                         validLength,
			NetWeight:                      validNetWeight,
			ExpirationRate:                 validExpirationRate,
			RecommendedFreezingTemperature: validRecommendedFreezingTemperature,
			FreezingRate:                   validFreezingRate,
			ProductTypeId:                  validProductTypeId,
			SellerId:                       &sellerId,
		}
		mockProductRepository := new(repository.MockProductRepository)
		mockProductTypeService := new(service.MockProductTypeService)
		mockSellerService := new(service.MockSellerService)
		productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

		productRequestWithSeller := validProductRequest
		productRequestWithSeller.SellerId = &sellerId

		expectedProduct := validProduct
		expectedProduct.SellerId = &sellerId

		mockProductRepository.On("CreateProduct", inputProduct).Return(expectedProduct, nil)

		_, err := productService.CreateProduct(productRequestWithSeller)
		require.NoError(t, err)
		mockProductRepository.AssertExpectations(t)
	})
}

func TestPatchProduct(t *testing.T) {
	mockProductRepository := new(repository.MockProductRepository)
	mockProductTypeService := new(service.MockProductTypeService)
	mockSellerService := new(service.MockSellerService)
	productService := NewProductService(mockProductRepository, mockProductTypeService, mockSellerService)

	baseProduct := models.Product{
		ID:                             1,
		ProductCode:                    "1234567890",
		Description:                    "Description",
		Width:                          10.0,
		Height:                         10.0,
		Length:                         10.0,
		NetWeight:                      10.0,
		ExpirationRate:                 10,
		RecommendedFreezingTemperature: 10.0,
		FreezingRate:                   10,
		ProductTypeId:                  1,
		SellerId:                       nil,
	}

	t.Run("should update product_code successfully", func(t *testing.T) {
		newProductCode := "9876543210"
		patchRequest := models.ProductPatchRequest{
			Id:          1,
			ProductCode: &newProductCode,
		}

		expectedProduct := baseProduct
		expectedProduct.ProductCode = newProductCode

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when product_code is empty", func(t *testing.T) {
		emptyProductCode := ""
		patchRequest := models.ProductPatchRequest{
			Id:          1,
			ProductCode: &emptyProductCode,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Product code cannot be empty")
	})

	t.Run("should update description successfully", func(t *testing.T) {
		newDescription := "Updated Description"
		patchRequest := models.ProductPatchRequest{
			Id:          1,
			Description: &newDescription,
		}

		expectedProduct := baseProduct
		expectedProduct.Description = newDescription

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when description is empty", func(t *testing.T) {
		emptyDescription := ""
		patchRequest := models.ProductPatchRequest{
			Id:          1,
			Description: &emptyDescription,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Description cannot be empty")
	})

	t.Run("should update width successfully", func(t *testing.T) {
		newWidth := 20.0
		patchRequest := models.ProductPatchRequest{
			Id:    1,
			Width: &newWidth,
		}

		expectedProduct := baseProduct
		expectedProduct.Width = newWidth

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when width is less than or equal to 0", func(t *testing.T) {
		invalidWidth := 0.0
		patchRequest := models.ProductPatchRequest{
			Id:    1,
			Width: &invalidWidth,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Width must be greater than 0")
	})

	t.Run("should update height successfully", func(t *testing.T) {
		newHeight := 25.0
		patchRequest := models.ProductPatchRequest{
			Id:     1,
			Height: &newHeight,
		}

		expectedProduct := baseProduct
		expectedProduct.Height = newHeight

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when height is less than or equal to 0", func(t *testing.T) {
		invalidHeight := -5.0
		patchRequest := models.ProductPatchRequest{
			Id:     1,
			Height: &invalidHeight,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Height must be greater than 0")
	})

	t.Run("should update length successfully", func(t *testing.T) {
		newLength := 30.0
		patchRequest := models.ProductPatchRequest{
			Id:     1,
			Length: &newLength,
		}

		expectedProduct := baseProduct
		expectedProduct.Length = newLength

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when length is less than or equal to 0", func(t *testing.T) {
		invalidLength := 0.0
		patchRequest := models.ProductPatchRequest{
			Id:     1,
			Length: &invalidLength,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Length must be greater than 0")
	})

	t.Run("should update net_weight successfully", func(t *testing.T) {
		newNetWeight := 15.5
		patchRequest := models.ProductPatchRequest{
			Id:        1,
			NetWeight: &newNetWeight,
		}

		expectedProduct := baseProduct
		expectedProduct.NetWeight = newNetWeight

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when net_weight is less than or equal to 0", func(t *testing.T) {
		invalidNetWeight := -1.0
		patchRequest := models.ProductPatchRequest{
			Id:        1,
			NetWeight: &invalidNetWeight,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Net weight must be greater than 0")
	})

	t.Run("should update expiration_rate successfully", func(t *testing.T) {
		newExpirationRate := 20
		patchRequest := models.ProductPatchRequest{
			Id:             1,
			ExpirationRate: &newExpirationRate,
		}

		expectedProduct := baseProduct
		expectedProduct.ExpirationRate = newExpirationRate

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when expiration_rate is less than or equal to 0", func(t *testing.T) {
		invalidExpirationRate := 0
		patchRequest := models.ProductPatchRequest{
			Id:             1,
			ExpirationRate: &invalidExpirationRate,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Expiration rate must be greater than 0")
	})

	t.Run("should update recommended_freezing_temperature successfully", func(t *testing.T) {
		newTemperature := -5.0
		patchRequest := models.ProductPatchRequest{
			Id:                             1,
			RecommendedFreezingTemperature: &newTemperature,
		}

		expectedProduct := baseProduct
		expectedProduct.RecommendedFreezingTemperature = newTemperature

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should update freezing_rate successfully", func(t *testing.T) {
		newFreezingRate := 15
		patchRequest := models.ProductPatchRequest{
			Id:           1,
			FreezingRate: &newFreezingRate,
		}

		expectedProduct := baseProduct
		expectedProduct.FreezingRate = newFreezingRate

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return error when freezing_rate is less than 0", func(t *testing.T) {
		invalidFreezingRate := -5
		patchRequest := models.ProductPatchRequest{
			Id:           1,
			FreezingRate: &invalidFreezingRate,
		}

		_, err := productService.patchProduct(baseProduct, patchRequest)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Freezing rate must be greater than 0")
	})

	t.Run("should update product_type_id successfully", func(t *testing.T) {
		newProductTypeId := 5
		patchRequest := models.ProductPatchRequest{
			Id:            1,
			ProductTypeId: &newProductTypeId,
		}

		expectedProduct := baseProduct
		expectedProduct.ProductTypeId = newProductTypeId

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should set seller_id to nil when seller_id is 0", func(t *testing.T) {
		productWithSeller := baseProduct
		sellerId := 5
		productWithSeller.SellerId = &sellerId

		zeroSellerId := 0
		patchRequest := models.ProductPatchRequest{
			Id:       1,
			SellerId: &zeroSellerId,
		}

		expectedProduct := productWithSeller
		expectedProduct.SellerId = nil

		updatedProduct, err := productService.patchProduct(productWithSeller, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
		require.Nil(t, updatedProduct.SellerId)
	})

	t.Run("should update seller_id successfully", func(t *testing.T) {
		newSellerId := 10
		patchRequest := models.ProductPatchRequest{
			Id:       1,
			SellerId: &newSellerId,
		}

		expectedProduct := baseProduct
		expectedProduct.SellerId = &newSellerId

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
		require.Equal(t, &newSellerId, updatedProduct.SellerId)
	})

	t.Run("should update multiple fields successfully", func(t *testing.T) {
		newProductCode := "9999999999"
		newDescription := "Multi-field Update"
		newWidth := 25.0
		newHeight := 30.0
		newSellerId := 15

		patchRequest := models.ProductPatchRequest{
			Id:          1,
			ProductCode: &newProductCode,
			Description: &newDescription,
			Width:       &newWidth,
			Height:      &newHeight,
			SellerId:    &newSellerId,
		}

		expectedProduct := baseProduct
		expectedProduct.ProductCode = newProductCode
		expectedProduct.Description = newDescription
		expectedProduct.Width = newWidth
		expectedProduct.Height = newHeight
		expectedProduct.SellerId = &newSellerId

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, expectedProduct, updatedProduct)
	})

	t.Run("should return original product when no fields are provided", func(t *testing.T) {
		patchRequest := models.ProductPatchRequest{
			Id: 1,
		}

		updatedProduct, err := productService.patchProduct(baseProduct, patchRequest)
		require.NoError(t, err)
		require.Equal(t, baseProduct, updatedProduct)
	})
}
