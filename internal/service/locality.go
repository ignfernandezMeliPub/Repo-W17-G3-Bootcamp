package service

import (
	"app/internal/repository/locality_repository"
	"app/pkg/models"
)

type LocalityService interface {
	CreateLocality(id string, localityName string, provinceName string, countryName string) (models.Locality, error)
	GetLocalitySellerCount(localityId string) (models.LocalitySellerCount, error)
	GetLocalitiesSellerCount() ([]models.LocalitySellerCount, error)
}

type LocalityServiceImpl struct {
	repository locality_repository.LocalityRepository
}

func NewLocalityServiceImpl(repository locality_repository.LocalityRepository) LocalityServiceImpl {
	return LocalityServiceImpl{repository: repository}
}

// CreateLocality Creates a new locality
func (s *LocalityServiceImpl) CreateLocality(id string, localityName string, provinceName string, countryName string) (models.Locality, error) {
	return s.repository.CreateLocality(models.Locality{Id: id, LocalityName: localityName, ProvinceName: provinceName, CountryName: countryName})
}

// GetLocalitySellerCount Returns LocalitySellerCount for a specific locality
func (s *LocalityServiceImpl) GetLocalitySellerCount(localityId string) (models.LocalitySellerCount, error) {
	return s.repository.GetLocalitySellerCount(localityId)
}

// GetLocalitiesSellerCount Returns LocalitySellerCount for all localities
func (s *LocalityServiceImpl) GetLocalitiesSellerCount() ([]models.LocalitySellerCount, error) {
	return s.repository.GetLocalitiesSellerCount()
}
