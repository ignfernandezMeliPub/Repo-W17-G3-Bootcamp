package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockPurchaseOrderRepository struct {
	mock.Mock
}

func (m *MockPurchaseOrderRepository) CreatePurchaseOrder(purchaseOrder models.PurchaseOrder) (p models.PurchaseOrder, err error) {
	args := m.Called(purchaseOrder)
	return args.Get(0).(models.PurchaseOrder), args.Error(1)
}
