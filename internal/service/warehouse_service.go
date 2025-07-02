package service

import (
	warehouserepository "app/internal/repository/warehouse_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type IWarehouseService interface {
	CreateWarehouse(vh models.Warehouse) (models.Warehouse, error)
	FindWarehouse() ([]models.Warehouse, error)
	FindWarehouseById(id int) (models.Warehouse, error)
	UpdateWarehouse(id int, w models.Warehouse) (models.Warehouse, error)
	DeleteWarehouse(id int) error
}

func NewWarehouseDefault(rp warehouserepository.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{rp: rp}
}

type WarehouseDefault struct {
	rp warehouserepository.WarehouseRepository
}

func (s *WarehouseDefault) CreateWarehouse(vh models.Warehouse) (models.Warehouse, error) {
	if vh.MinimumCapacity < 0 {
		return models.Warehouse{}, &custom_errors.InvalidArgValueErr{Argument: "minimun_capacity", Value: vh.MinimumCapacity, ExtraInfo: "minimun_capacity cannot be less than zero"}
	}

	return s.rp.CreateWarehouse(vh)
}

func (s *WarehouseDefault) FindWarehouse() ([]models.Warehouse, error) {
	return s.rp.FindWarehouse()
}

func (s *WarehouseDefault) FindWarehouseById(id int) (models.Warehouse, error) {
	return s.rp.FindWarehouseById(id)
}

func (s *WarehouseDefault) UpdateWarehouse(id int, w models.Warehouse) (models.Warehouse, error) {
	_, err := s.rp.FindWarehouseById(id)
	if err != nil {
		return models.Warehouse{}, &custom_errors.ResourceNotFoundError{}
	}

	if w.MinimumCapacity < 0 {
		return models.Warehouse{}, &custom_errors.InvalidArgValueErr{Argument: "minimun_capacity", Value: w.MinimumCapacity, ExtraInfo: "minimun_capacity cannot be less than zero"}
	}
	w.Id = id

	return s.rp.UpdateWarehouse(id, w)
}

func (s *WarehouseDefault) DeleteWarehouse(id int) error {
	_, err := s.rp.FindWarehouseById(id)
	if err != nil {
		return &custom_errors.ResourceNotFoundError{}
	}
	return s.rp.DeleteWarehouse(id)
}
