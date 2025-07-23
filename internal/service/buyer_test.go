package service

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"app/test/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

var buyer1 = models.Buyer{
	Id:           1,
	CardNumberId: "12345",
	FirstName:    "John",
	LastName:     "Doe",
}

var buyer2 = models.Buyer{
	Id:           2,
	CardNumberId: "67890",
	FirstName:    "Jane",
	LastName:     "Smith",
}

var buyer3 = models.Buyer{
	Id:           3,
	CardNumberId: "123456",
	FirstName:    "John",
	LastName:     "Doe",
}

var buyerPatch = models.BuyerPatch{
	FirstName: &buyer1.FirstName,
	LastName:  &buyer1.LastName,
}

var buyer1PurchaseOrdersCount = models.BuyerPurchaseOrdersCount{
	Id:                  1,
	CardNumberId:        "12345",
	FirstName:           "John",
	LastName:            "Doe",
	PurchaseOrdersCount: 1,
}

var buyer2PurchaseOrdersCount = models.BuyerPurchaseOrdersCount{
	Id:                  2,
	CardNumberId:        "67890",
	FirstName:           "Jane",
	LastName:            "Smith",
	PurchaseOrdersCount: 2,
}

var buyer3PurchaseOrdersCount = models.BuyerPurchaseOrdersCount{
	Id:                  3,
	CardNumberId:        "123456",
	FirstName:           "John",
	LastName:            "Doe",
	PurchaseOrdersCount: 1,
}

func TestGetAllBuyers(t *testing.T) {
	t.Run("should return all buyers successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetAllBuyers").Return([]models.Buyer{buyer1, buyer2, buyer3}, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyers, err := service.GetAllBuyers()

		// Assert
		require.NoError(t, err)
		require.Equal(t, []models.Buyer{buyer1, buyer2, buyer3}, buyers)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetAllBuyers").Return([]models.Buyer{}, custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyers, err := service.GetAllBuyers()

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		require.Empty(t, buyers)
		mockRepository.AssertExpectations(t)
	})
}

func TestGetBuyerById(t *testing.T) {
	t.Run("should return buyer by id successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyerById", 1).Return(buyer1, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.GetBuyerById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, buyer1, buyer)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyerById", 1).Return(models.Buyer{}, custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.GetBuyerById(1)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		require.Empty(t, buyer)
		mockRepository.AssertExpectations(t)
	})
}

func TestCreateBuyer(t *testing.T) {
	t.Run("should create buyer successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("CreateBuyer", buyer1).Return(buyer1, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.CreateBuyer(buyer1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, buyer1, buyer)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle conflict error on duplicate card number id", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("CreateBuyer", buyer1).Return(models.Buyer{}, custom_errors.ErrUniqueAttributeViolationError)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.CreateBuyer(buyer1)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrUniqueAttributeViolationError)
		require.Empty(t, buyer)
		mockRepository.AssertExpectations(t)
	})
}

func TestUpdateBuyer(t *testing.T) {
	t.Run("should update buyer successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyerById", 1).Return(buyer1, nil)
		mockRepository.On("UpdateBuyer", buyer1).Return(buyer1, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.UpdateBuyerById(1, buyerPatch)

		// Assert
		require.NoError(t, err)
		require.Equal(t, buyer1, buyer)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle conflict error on duplicate card number id", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyerById", 1).Return(buyer1, nil)
		mockRepository.On("UpdateBuyer", buyer1).Return(models.Buyer{}, custom_errors.ErrUniqueAttributeViolationError)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.UpdateBuyerById(1, buyerPatch)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrUniqueAttributeViolationError)
		require.Empty(t, buyer)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyerById", 1).Return(models.Buyer{}, custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyer, err := service.UpdateBuyerById(1, buyerPatch)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		require.Empty(t, buyer)
		mockRepository.AssertExpectations(t)
	})
}

func TestDeleteBuyerById(t *testing.T) {
	t.Run("should delete buyer by id successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("DeleteBuyerById", 1).Return(nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		err := service.DeleteBuyerById(1)

		// Assert
		require.NoError(t, err)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("DeleteBuyerById", 1).Return(custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		err := service.DeleteBuyerById(1)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		mockRepository.AssertExpectations(t)
	})
}

func TestGetBuyersPurchaseOrdersCount(t *testing.T) {
	t.Run("should get all buyers purchase orders count successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return([]models.BuyerPurchaseOrdersCount{buyer1PurchaseOrdersCount, buyer2PurchaseOrdersCount, buyer3PurchaseOrdersCount}, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyersPurchaseOrdersCount, err := service.GetBuyersPurchaseOrdersCount(nil)

		// Assert
		require.NoError(t, err)
		require.Equal(t, []models.BuyerPurchaseOrdersCount{buyer1PurchaseOrdersCount, buyer2PurchaseOrdersCount, buyer3PurchaseOrdersCount}, buyersPurchaseOrdersCount)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error for all buyers purchase orders count", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyersPurchaseOrdersCount", (*int)(nil)).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyerPurchaseOrdersCount, err := service.GetBuyersPurchaseOrdersCount(nil)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		require.Empty(t, buyerPurchaseOrdersCount)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should get buyer purchase orders count successfully", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyersPurchaseOrdersCount", &buyer1.Id).Return([]models.BuyerPurchaseOrdersCount{buyer1PurchaseOrdersCount}, nil)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyerPurchaseOrdersCount, err := service.GetBuyersPurchaseOrdersCount(&buyer1.Id)

		// Assert
		require.NoError(t, err)
		require.Equal(t, []models.BuyerPurchaseOrdersCount{buyer1PurchaseOrdersCount}, buyerPurchaseOrdersCount)
		mockRepository.AssertExpectations(t)
	})

	t.Run("should handle not found error for buyer purchase orders count", func(t *testing.T) {
		// Arrange
		mockRepository := new(repository.MockBuyerRepository)
		mockRepository.On("GetBuyersPurchaseOrdersCount", &buyer1.Id).Return([]models.BuyerPurchaseOrdersCount{}, custom_errors.ErrNotFound)

		service := NewBuyerDefault(mockRepository)

		// Act
		buyerPurchaseOrdersCount, err := service.GetBuyersPurchaseOrdersCount(&buyer1.Id)

		// Assert
		require.Error(t, err)
		require.ErrorAs(t, err, &custom_errors.ErrNotFound)
		require.Empty(t, buyerPurchaseOrdersCount)
		mockRepository.AssertExpectations(t)
	})
}
