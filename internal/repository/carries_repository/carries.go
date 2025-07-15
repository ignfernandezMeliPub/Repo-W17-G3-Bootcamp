package carries_repository

import "app/pkg/models"

type CarriesRepository interface {
	CreateCarrie(c models.Carries) (models.Carries, error)
	GetCarriesReport(localityId *string) ([]models.CarriesReport, error)
}
