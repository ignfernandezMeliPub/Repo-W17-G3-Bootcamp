package warehouse_repository

import "app/pkg/models"

type WarehouseRepository interface {
	CreateWarehouse(wh models.Warehouse) (models.Warehouse, error)
	GetAllWarehouses() ([]models.Warehouse, error)
	GetWarehouseById(id int) (models.Warehouse, error)
	UpdateWarehouseById(id int, w models.Warehouse) (models.Warehouse, error)
	DeleteWarehouseById(id int) error
}
