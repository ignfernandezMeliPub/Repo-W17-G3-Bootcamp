package service

import (
	warehouserepository "app/internal/repository/repositories/warehouse_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type IWarehouseService interface {
	GetAllWarehouses() ([]models.Warehouse, error)
	GetWarehouseById(id int) (models.Warehouse, error)
	CreateWarehouse(vh models.Warehouse) (models.Warehouse, error)
	UpdateWarehouseById(id int, w models.Warehouse) (models.Warehouse, error)
	DeleteWarehouse(id int) error
}

func NewWarehouseDefault(rp warehouserepository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

type WarehouseDefault struct {
	rp warehouserepository.WarehouseRepository
}

func (s *WarehouseDefault) CreateWarehouse(vh models.Warehouse) (models.Warehouse, error) {

	return s.rp.CreateWarehouse(vh)
}
func (s *WarehouseDefault) GetAllWarehouses() ([]models.Warehouse, error) {
	return s.rp.GetAllWarehouses()

}

func (s *WarehouseDefault) GetWarehouseById(id int) (models.Warehouse, error) {
	return s.rp.GetWarehouseById(id)
}

func (s *WarehouseDefault) UpdateWarehouseById(id int, w models.Warehouse) (models.Warehouse, error) {
	_, err := s.rp.GetWarehouseById(id)
	if err != nil {
		return models.Warehouse{}, &custom_errors.ResourceNotFoundError{}
	}

	if w.MinimumCapacity < 0 {
		return models.Warehouse{}, &custom_errors.InvalidArgValueErr{Argument: "minimun_capacity", Value: w.MinimumCapacity, ExtraInfo: "minimun_capacity cannot be less than zero"}
	}
	w.Id = id

	return s.rp.UpdateWarehouseById(id, w)
}

func (s *WarehouseDefault) DeleteWarehouse(id int) error {
	_, err := s.rp.GetWarehouseById(id)
	if err != nil {
		return &custom_errors.ResourceNotFoundError{}
	}
	return s.rp.DeleteWarehouseById(id)
}
