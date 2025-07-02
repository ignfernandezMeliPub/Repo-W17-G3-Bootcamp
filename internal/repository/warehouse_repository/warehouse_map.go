package warehouserepository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"strings"
)

type WarehouseRepositoryMap struct {
	db map[int]models.Warehouse
}

func NewWarehouseMap(db map[int]models.Warehouse) *WarehouseRepositoryMap {

	data := make(map[int]models.Warehouse)
	if db != nil {
		data = db
	}
	return &WarehouseRepositoryMap{db: data}
}

func (r *WarehouseRepositoryMap) CreateWarehouse(wh models.Warehouse) (models.Warehouse, error) {

	var newId int
	for id, w := range r.db {
		if id >= newId {
			newId = id + 1
		}
		if strings.EqualFold(w.Warehouse_code, wh.Warehouse_code) {
			return models.Warehouse{}, &custom_errors.ResourceConflictError{
				Value:     w.Warehouse_code,
				Argument:  "warehouse_code",
				ExtraInfo: "warehouse_code must be unique",
			}
		}
	}
	wh.Id = newId
	r.db[newId] = wh
	return wh, nil
}

func (r *WarehouseRepositoryMap) FindWarehouse() ([]models.Warehouse, error) {

	if len(r.db) == 0 {
		return nil, &custom_errors.ResourceNotFoundError{}
	}
	w := make([]models.Warehouse, len(r.db))
	i := 0
	for _, value := range r.db {
		w[i] = value
		i++
	}

	return w, nil
}

func (r *WarehouseRepositoryMap) FindWarehouseById(id int) (models.Warehouse, error) {
	wh, exist := r.db[id]
	if !exist {
		return models.Warehouse{}, &custom_errors.ResourceNotFoundError{}
	}
	return wh, nil
}

func (r *WarehouseRepositoryMap) UpdateWarehouse(id int, w models.Warehouse) (models.Warehouse, error) {
	if r.FindWarehouseByCode(w.Warehouse_code) {
		return models.Warehouse{}, &custom_errors.ResourceConflictError{
			Value:     w.Warehouse_code,
			Argument:  "warehouse_code",
			ExtraInfo: "warehouse_code must be unique"}
	}

	r.db[id] = w
	return w, nil
}

func (r *WarehouseRepositoryMap) DeleteWarehouse(id int) error {
	delete(r.db, id)
	return nil
}

func (r *WarehouseRepositoryMap) FindWarehouseByCode(code string) bool {

	for _, w := range r.db {
		if strings.EqualFold(w.Warehouse_code, code) {
			return true
		}
	}
	return false
}
