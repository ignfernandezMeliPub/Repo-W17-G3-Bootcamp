package carries_repository

import "app/pkg/models"

type CarriesRepository interface {
	CreateCarries(c models.Carries) (int64, error)
}
