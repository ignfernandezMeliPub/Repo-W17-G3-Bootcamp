package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockSellerService struct {
	mock.Mock
}

func (m *MockSellerService) GetAllSellers() ([]models.Seller, error) {
	args := m.Called()
	return args.Get(0).([]models.Seller), args.Error(1)
}

func (m *MockSellerService) GetSellerById(id int) (models.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) CreateSeller(companyId int, companyName string, address string, telephone string, localityId string) (models.Seller, error) {
	args := m.Called(companyId, companyName, address, telephone, localityId)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error) {
	args := m.Called(id, companyId, companyName, address, telephone)
	return args.Get(0).(models.Seller), args.Error(1)
}

func (m *MockSellerService) DeleteSellerById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
