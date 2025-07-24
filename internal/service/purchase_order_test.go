package service

import (
	"app/pkg/models"
	"app/test/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var purchaseOrder = models.PurchaseOrder{
	Id:           1,
	OrderNumber:  "1234567890",
	OrderDate:    time.Now(),
	TrackingCode: "1234567890",
	BuyerId:      1,
}

func TestNewPurchaseOrderDefault(t *testing.T) {
	mockRepository := new(repository.MockPurchaseOrderRepository)

	service := NewPurchaseOrderDefault(mockRepository)

	require.NotNil(t, service)
}

func TestPurchaseOrderServiceDefault_CreatePurchaseOrder(t *testing.T) {
	mockRepository := new(repository.MockPurchaseOrderRepository)
	mockRepository.On("CreatePurchaseOrder", purchaseOrder).Return(purchaseOrder, nil)
	service := NewPurchaseOrderDefault(mockRepository)

	service.CreatePurchaseOrder(purchaseOrder)

	mockRepository.AssertExpectations(t)
}
