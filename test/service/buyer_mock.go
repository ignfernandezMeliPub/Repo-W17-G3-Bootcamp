package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockBuyerService struct {
	mock.Mock
}

func (m *MockBuyerService) GetAllBuyers() ([]models.Buyer, error) {
	args := m.Called()
	return args.Get(0).([]models.Buyer), args.Error(1)
}

func (m *MockBuyerService) GetBuyerById(id int) (models.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) CreateBuyer(buyer models.Buyer) (models.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) UpdateBuyerById(id int, buyerPatch models.BuyerPatch) (models.Buyer, error) {
	args := m.Called(id, buyerPatch)
	return args.Get(0).(models.Buyer), args.Error(1)
}

func (m *MockBuyerService) DeleteBuyerById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBuyerService) GetBuyersPurchaseOrdersCount(buyerId *int) ([]models.BuyerPurchaseOrdersCount, error) {
	args := m.Called(buyerId)
	return args.Get(0).([]models.BuyerPurchaseOrdersCount), args.Error(1)
}
