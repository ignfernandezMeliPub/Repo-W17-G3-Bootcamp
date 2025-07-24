package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockInboundOrderService struct {
	mock.Mock
}

func (m *MockInboundOrderService) CreateInboundOrder(details models.InboundOrderRequestBody) (inboundOrder models.InboundOrder, err error) {

	args := m.Called(details)
	return args.Get(0).(models.InboundOrder), args.Error(1)

}
