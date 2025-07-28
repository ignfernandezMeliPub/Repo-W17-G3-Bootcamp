package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockCarriesService struct {
	mock.Mock
}

func (m *MockCarriesService) CreateCarrie(c models.Carries) (models.Carries, error) {
	args := m.Called(c)
	return args.Get(0).(models.Carries), args.Error(1)
}
