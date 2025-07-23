package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockBuyerRepository struct {
	mock.Mock
}

func (m *MockBuyerRepository) GetAllBuyers() (b []models.Buyer, err error) {
	args := m.Called()
	return args.Get(0).([]models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) GetBuyerById(id int) (b models.Buyer, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) CreateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) UpdateBuyer(buyer models.Buyer) (updatedBuyer models.Buyer, err error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) DeleteBuyerById(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBuyerRepository) GetBuyersPurchaseOrdersCount(buyerId *int) (b []models.BuyerPurchaseOrdersCount, err error) {
	args := m.Called(buyerId)
	return args.Get(0).([]models.BuyerPurchaseOrdersCount), args.Error(1)
}
