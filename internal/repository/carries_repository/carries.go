package carries_repository

import "app/pkg/models"

type WarehouseRepository interface {
	CreateCarries(c models.Carries) (int64, error)
}
