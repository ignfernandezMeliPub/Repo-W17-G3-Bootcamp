package warehouse_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setupWarehouseRepository(t *testing.T) (*WarehouseSql, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewWarehouseSql(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

var temp = 1.0

func TestWarehouseSQL_CreateWarehouse(t *testing.T) {
	repo, mock, cleanup := setupWarehouseRepository(t)
	defer cleanup()
	const query = `
		INSERT INTO warehouses ( warehouse_code, address, telephone, minimum_capacity, minimum_temperature ) VALUES (?, ?, ?, ?, ?)`
	warehouse := models.Warehouse{
		WarehouseCode:      "WH001",
		Address:            "Street 1",
		Telephone:          "1234567890",
		MinimumCapacity:    100,
		MinimumTemperature: &temp,
	}
	t.Run("create warehouse success", func(t *testing.T) {
		expectedWarehouse := warehouse
		expectedWarehouse.Id = 1

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature).
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.CreateWarehouse(warehouse)
		require.NoError(t, err)
		require.Equal(t, expectedWarehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate warehouse code", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'WH001' for key 'warehouses.warehouse_code'",
			})

		_, err := repo.CreateWarehouse(warehouse)
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.CreateWarehouse(warehouse)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestWarehouseSQL_GetAllWarehouses(t *testing.T) {
	repo, mock, cleanup := setupWarehouseRepository(t)
	defer cleanup()
	const query = `
		SELECT id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature 
		FROM warehouses`

	t.Run("should return all warehouses successfully", func(t *testing.T) {
		expectedWarehouses := []models.Warehouse{
			{
				Id:                 1,
				WarehouseCode:      "WH001",
				Address:            "Street 1",
				Telephone:          "1234567890",
				MinimumCapacity:    100,
				MinimumTemperature: &temp,
			},
			{
				Id:                 2,
				WarehouseCode:      "WH002",
				Address:            "Street 2",
				Telephone:          "1234567891",
				MinimumCapacity:    150,
				MinimumTemperature: &temp,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "WH001", "Street 1", "1234567890", 100, 1.0).
			AddRow(2, "WH002", "Street 2", "1234567891", 150, 1.0)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		result, err := repo.GetAllWarehouses()
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, expectedWarehouses, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no rows exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"})

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnRows(rows)

		result, err := repo.GetAllWarehouses()
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.GetAllWarehouses()
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSQL_GetWarehouseById(t *testing.T) {
	repo, mock, cleanup := setupWarehouseRepository(t)
	defer cleanup()
	const query = `
		SELECT id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature 
		FROM warehouses WHERE id = ?`

	t.Run("should return warehouse by id successfully", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &temp,
		}

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "WH001", "Street 1", "1234567890", 100, 1.0)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnRows(rows)

		result, err := repo.GetWarehouseById(1)
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no rows exist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"})

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnRows(rows)

		result, err := repo.GetWarehouseById(1)
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		result, err := repo.GetWarehouseById(1)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWarehouseSQL_UpdateWarehouseById(t *testing.T) {
	repo, mock, cleanup := setupWarehouseRepository(t)
	defer cleanup()
	const updateQuery = `
		UPDATE warehouses SET warehouse_code = ?, address = ?, telephone = ?, minimum_capacity = ?, minimum_temperature = ? WHERE id = ?`

	const selectQuery = `
		SELECT 
			id,
			warehouse_code,
			address,
			telephone,
			minimum_capacity,
			minimum_temperature 
		FROM warehouses 
		WHERE id = ?`

	t.Run("should update warehouse by id successfully", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &temp,
		}

		mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature, warehouse.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		selectQuery := `
		SELECT 
	    id,
		warehouse_code,
		address,
		telephone,
		minimum_capacity,
		minimum_temperature 
		FROM warehouses 
		WHERE id = ?`

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "WH001", "Street 1", "1234567890", 100, 1.0)

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(warehouse.Id).
			WillReturnRows(rows)

		result, err := repo.UpdateWarehouseById(warehouse.Id, warehouse)
		require.NoError(t, err)
		require.Equal(t, warehouse, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return not found error when no rows exist", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &temp,
		}
		mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature, warehouse.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		rows := sqlmock.NewRows([]string{"id", "warehouse_code", "address", "telephone", "minimum_capacity", "minimum_temperature"})

		mock.ExpectQuery(regexp.QuoteMeta(selectQuery)).
			WithArgs(warehouse.Id).
			WillReturnRows(rows)

		result, err := repo.UpdateWarehouseById(warehouse.Id, warehouse)
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when query fails", func(t *testing.T) {
		warehouse := models.Warehouse{
			Id:                 1,
			WarehouseCode:      "WH001",
			Address:            "Street 1",
			Telephone:          "1234567890",
			MinimumCapacity:    100,
			MinimumTemperature: &temp,
		}
		mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
			WithArgs(warehouse.WarehouseCode, warehouse.Address, warehouse.Telephone, warehouse.MinimumCapacity, warehouse.MinimumTemperature, warehouse.Id).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.UpdateWarehouseById(warehouse.Id, warehouse)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when no fields are provided to update", func(t *testing.T) {
		emptyWarehouse := models.Warehouse{
			Id: 1,
		}

		_, err := repo.UpdateWarehouseById(emptyWarehouse.Id, emptyWarehouse)

		require.Error(t, err)
		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestWarehouseSQL_DeleteWarehouseById(t *testing.T) {
	repo, mock, cleanup := setupWarehouseRepository(t)
	defer cleanup()
	const query = `
		DELETE FROM warehouses WHERE id = ?`

	t.Run("should delete warehouse by id successfully", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteWarehouseById(1)
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return not found error when no rows exist", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteWarehouseById(1)
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when query fails", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		err := repo.DeleteWarehouseById(1)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
