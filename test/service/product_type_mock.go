package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockProductTypeService struct {
	mock.Mock
}

func (m *MockProductTypeService) GetProductTypeById(id int) (models.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(models.ProductType), args.Error(1)
}
