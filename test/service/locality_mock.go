package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) CreateLocality(id string, localityName string, provinceName string, countryName string) (models.Locality, error) {
	args := m.Called(id, localityName, provinceName, countryName)
	return args.Get(0).(models.Locality), args.Error(1)
}

func (m *MockLocalityService) GetCarriesReport(localityId string) ([]models.CarriesReport, error) {
	args := m.Called(localityId)
	return args.Get(0).([]models.CarriesReport), args.Error(1)
}

func (m *MockLocalityService) GetLocalitySellerCount(localityId string) (models.LocalitySellerCount, error) {
	args := m.Called(localityId)
	return args.Get(0).(models.LocalitySellerCount), args.Error(1)
}

func (m *MockLocalityService) GetLocalitiesSellerCount() ([]models.LocalitySellerCount, error) {
	args := m.Called()
	return args.Get(0).([]models.LocalitySellerCount), args.Error(1)
}
