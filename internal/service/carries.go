package service

import (
	"app/internal/repository/repositories/carries_repository"
	"app/pkg/models"
)

type CarriesService interface {
	CreateCarrie(c models.Carries) (models.Carries, error)
}

type CarriesServiceDefault struct {
	rp carries_repository.CarriesRepository
}

func NewCarriesServiceDefault(rp carries_repository.CarriesRepository) *CarriesServiceDefault {
	return &CarriesServiceDefault{rp: rp}
}

func (s *CarriesServiceDefault) CreateCarrie(c models.Carries) (models.Carries, error) {
	return s.rp.CreateCarrie(c)
}
