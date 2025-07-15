package warehouse_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
)

type WarehouseSql struct {
	db *sql.DB
}

func NewWarehouseSql(db *sql.DB) *WarehouseSql {
	return &WarehouseSql{
		db: db,
	}
}

// ************ const queries ************
const queryCreateWarehouse = `
	INSERT INTO warehouses (
		warehouse_code, 
		address, 
		telephone, 
		minimum_capacity, 
		minimum_temperature
		)VALUES (?, ?, ?, ?, ?)`

const queryUpdateWarehouseById = `
	UPDATE warehouses 
	SET 
		warehouse_code = ?, 
		address = ?, 
		telephone = ?, 
		minimum_capacity = ?, 
		minimum_temperature = ? 
		WHERE id = ?`

const queryGetAllWarehouses = `
	SELECT 
	id,
		warehouse_code,
		address,
		telephone,
		minimum_capacity,
		minimum_temperature
		FROM warehouses`
const queryGetWarehouseById = `
	SELECT 
	id,
		warehouse_code,
		address,
		telephone,
		minimum_capacity,
		minimum_temperature 
		FROM warehouses 
		WHERE id = ?`
const queryDeleteWarehouseById = `DELETE FROM warehouses WHERE id = ?`

func (r *WarehouseSql) CreateWarehouse(wh models.Warehouse) (models.Warehouse, error) {
	args := []any{wh.WarehouseCode, wh.Address, wh.Telephone, wh.MinimumCapacity, wh.MinimumTemperature}
	id, err := sql_utils.Insert(r.db, queryCreateWarehouse, args)
	if err != nil {
		return models.Warehouse{}, err
	}
	wh.Id = int(id)
	return wh, nil
}

func (r *WarehouseSql) GetAllWarehouses() ([]models.Warehouse, error) {

	return sql_utils.Query[models.Warehouse](r.db, queryGetAllWarehouses, nil)
}

func (r *WarehouseSql) GetWarehouseById(id int) (models.Warehouse, error) {
	return sql_utils.QueryRow[models.Warehouse](r.db, queryGetWarehouseById, []any{id})
}

func (r *WarehouseSql) DeleteWarehouseById(id int) error {
	row, err := sql_utils.Delete(r.db, queryDeleteWarehouseById, []any{id})
	if err != nil {
		return err
	}
	if row == 0 {
		return custom_errors.ErrNotFound
	}
	return nil
}

func (r *WarehouseSql) UpdateWarehouseById(id int, wh models.Warehouse) (models.Warehouse, error) {
	var columns []string
	var args []any
	if wh.WarehouseCode != "" {
		columns = append(columns, "warehouse_code = ?")
		args = append(args, wh.WarehouseCode)
	}
	if wh.Address != "" {
		columns = append(columns, "address = ?")
		args = append(args, wh.Address)
	}
	if wh.Telephone != "" {
		columns = append(columns, "telephone = ?")
		args = append(args, wh.Telephone)
	}
	if wh.MinimumCapacity != 0 {
		columns = append(columns, "minimum_capacity = ?")
		args = append(args, wh.MinimumCapacity)
	}
	if wh.MinimumTemperature != nil {
		columns = append(columns, "minimum_temperature = ?")
		args = append(args, wh.MinimumTemperature)
	}
	if len(columns) == 0 {
		return models.Warehouse{}, &custom_errors.MandatoryArgMissingErr{Argument: "warehouse_code or address or telephone or minimum_capacity or minimum_temperature "}
	}
	args = append(args, id)
	_, err := sql_utils.Update(r.db, queryUpdateWarehouseById, args)
	if err != nil {
		return models.Warehouse{}, err
	}
	return r.GetWarehouseById(id)
}
