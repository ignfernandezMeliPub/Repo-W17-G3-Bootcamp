package service

import (
	warehouserepository "app/internal/repository/warehouse_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"strings"
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

	if err := validateAttributes(vh); err != nil {
		return models.Warehouse{}, err
	}

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

func validateAttributes(wh models.Warehouse) error {
	if strings.TrimSpace(wh.WarehouseCode) == "" {
		return &custom_errors.InvalidArgValueErr{Argument: "warehouse_code", Value: wh.WarehouseCode, ExtraInfo: "warehouse_code cannot be empty"}
	}
	if strings.TrimSpace(wh.Address) == "" {
		return &custom_errors.InvalidArgValueErr{Argument: "address", Value: wh.Address, ExtraInfo: "address cannot be empty"}
	}
	if strings.TrimSpace(wh.Telephone) == "" {
		return &custom_errors.InvalidArgValueErr{Argument: "telephone", Value: wh.Telephone, ExtraInfo: "telephone cannot be empty"}
	}
	if wh.MinimumCapacity < 0 {
		return &custom_errors.InvalidArgValueErr{Argument: "minimun_capacity", Value: wh.MinimumCapacity, ExtraInfo: "minimun_capacity cannot be less than zero"}
	}
	if wh.MinimumTemperature == nil {
		return &custom_errors.InvalidArgValueErr{Argument: "minimun_temperature", Value: wh.MinimumTemperature, ExtraInfo: "minimun_temperature cannot be empty"}
	}
	return nil
}
