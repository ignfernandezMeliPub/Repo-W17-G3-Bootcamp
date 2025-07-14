package carries_repository

import "app/pkg/models"

type WarehouseRepository interface {
	Store(c models.Carries) (int64, error)
}
