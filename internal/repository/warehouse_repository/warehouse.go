package warehouse_repository

import "app/pkg/models"

type WarehouseRepository interface {
	CreateWarehouse(wh models.Warehouse) (models.Warehouse, error)
	FindWarehouse() ([]models.Warehouse, error)
	FindWarehouseById(id int) (models.Warehouse, error)
	UpdateWarehouse(id int, w models.Warehouse) (models.Warehouse, error)
	DeleteWarehouse(id int) error
}
