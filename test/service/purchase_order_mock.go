package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockPurchaseOrderService struct {
	mock.Mock
}

func (m *MockPurchaseOrderService) CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (models.PurchaseOrder, error) {
	args := m.Called(purchaseOrder)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}
