package locality_repository

import "app/pkg/models"

type LocalityRepository interface {
	CreateLocality(locality models.Locality) (models.Locality, error)
	GetCarriesReport(localityId string) ([]models.CarriesReport, error)
}
