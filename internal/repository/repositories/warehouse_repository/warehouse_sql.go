package warehouse_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"strconv"
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
		) VALUES (?, ?, ?, ?, ?)`

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
	sql_utils.LogAudit("CreateWarehouse", logger.LogStatusInProgress, "Insert warehouse")

	args := []any{wh.WarehouseCode, wh.Address, wh.Telephone, wh.MinimumCapacity, wh.MinimumTemperature}
	id, err := sql_utils.Insert(r.db, queryCreateWarehouse, args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateWarehouse", "Insert warehouse", err)
		return models.Warehouse{}, err
	}
	wh.Id = int(id)

	sql_utils.LogAudit("CreateWarehouse", logger.LogStatusSuccess, "Insert warehouse. Id: "+strconv.Itoa(wh.Id))
	return wh, nil
}

func (r *WarehouseSql) GetAllWarehouses() ([]models.Warehouse, error) {
	sql_utils.Log("GetAllWarehouses", logger.LogStatusInProgress, "Select all warehouses")

	whs, err := sql_utils.Query[models.Warehouse](r.db, queryGetAllWarehouses, nil)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetAllWarehouses", "Select all warehouses", err)
		return nil, err
	}

	sql_utils.Log("GetAllWarehouses", logger.LogStatusSuccess, "Select all warehouses")
	return whs, nil
}

func (r *WarehouseSql) GetWarehouseById(id int) (models.Warehouse, error) {
	sql_utils.Log("GetWarehouseById", logger.LogStatusInProgress, "Select warehouse by id "+strconv.Itoa(id))

	wh, err := sql_utils.QueryRow[models.Warehouse](r.db, queryGetWarehouseById, []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetWarehouseById", "Select warehouse by id "+strconv.Itoa(id), err)
		return models.Warehouse{}, err
	}

	sql_utils.Log("GetWarehouseById", logger.LogStatusSuccess, "Select warehouse by id "+strconv.Itoa(id))
	return wh, nil
}

func (r *WarehouseSql) DeleteWarehouseById(id int) error {
	sql_utils.LogAudit("DeleteWarehouseById", logger.LogStatusInProgress, "Delete warehouse by id: "+strconv.Itoa(id))

	row, err := sql_utils.Delete(r.db, queryDeleteWarehouseById, []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteWarehouseById", "Delete warehouse by id: "+strconv.Itoa(id), err)
		return err
	}
	if row == 0 {
		err = custom_errors.ErrNotFound
		sql_utils.LogAuditError("DeleteWarehouseById", "Delete warehouse by id: "+strconv.Itoa(id), err)
		return err
	}

	sql_utils.LogAudit("DeleteWarehouseById", logger.LogStatusSuccess, "Delete warehouse by id: "+strconv.Itoa(id))
	return nil
}

func (r *WarehouseSql) UpdateWarehouseById(id int, wh models.Warehouse) (models.Warehouse, error) {
	sql_utils.LogAudit("UpdateWarehouseById", logger.LogStatusInProgress, "Update warehouse by id: "+strconv.Itoa(id))

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
		err := &custom_errors.MandatoryArgMissingErr{Argument: "warehouse_code or address or telephone or minimum_capacity or minimum_temperature "}
		sql_utils.LogAuditError("UpdateWarehouseById", "Update warehouse by id: "+strconv.Itoa(id), err)
		return models.Warehouse{}, err
	}
	args = append(args, id)
	_, err := sql_utils.Update(r.db, queryUpdateWarehouseById, args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("UpdateWarehouseById", "Update warehouse by id: "+strconv.Itoa(id), err)
		return models.Warehouse{}, err
	}

	warehouse, err := r.GetWarehouseById(id)

	if err == nil {
		sql_utils.LogAudit("UpdateWarehouseById", logger.LogStatusSuccess, "Update warehouse by id: "+strconv.Itoa(id))
	}

	return warehouse, err
}
