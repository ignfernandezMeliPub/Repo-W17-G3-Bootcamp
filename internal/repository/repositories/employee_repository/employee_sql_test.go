package employee_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupEmployeeRepository(t *testing.T) (EmployeeRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewEmployeeDb(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestGetAllEmployees(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	t.Run("success get all employees", func(t *testing.T) {
		// Arrange
		expectedEmployees := []models.Employee{
			{
				Id: 1, EmployeeAttributes: models.EmployeeAttributes{
					CardNumberId: "EMP001",
					FirstName:    "Raul",
					LastName:     "García",
					WarehouseId:  1},
			}, {
				Id: 2, EmployeeAttributes: models.EmployeeAttributes{
					CardNumberId: "EMP002",
					FirstName:    "Juan",
					LastName:     "Sierra",
					WarehouseId:  2},
			},
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).AddRow(
			1, "EMP001", "Raul", "García", 1,
		).AddRow(
			2, "EMP002", "Juan", "Sierra", 2,
		)

		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees`).
			WillReturnRows(rows)

		// Act
		employees, err := repo.GetAllEmployees()

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedEmployees, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error no rows", func(t *testing.T) {

		//Arrange

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})
		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees`).WillReturnRows(rows)

		// Act
		employees, err := repo.GetAllEmployees()

		// Assert

		expectedEmployees := []models.Employee(nil)
		expectedError := custom_errors.ErrNotFound

		assert.Equal(t, expectedEmployees, employees)
		assert.Equal(t, expectedError, err)

	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees`).
			WillReturnError(sql.ErrConnDone)

		// Act
		employees, err := repo.GetAllEmployees()

		// Assert

		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Empty(t, employees)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestGetEmployeeById(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	t.Run("success get employee", func(t *testing.T) {

		// Arrange

		id := 1

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).AddRow(
			1, "EMP001", "Raul", "García", 1,
		)

		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(rows)

		// Act
		employee, err := repo.GetEmployeeById(id)

		// Assert
		expectedEmployee := models.Employee{
			Id: 1,
			EmployeeAttributes: models.EmployeeAttributes{
				CardNumberId: "EMP001",
				FirstName:    "Raul",
				LastName:     "García",
				WarehouseId:  1,
			},
		}

		require.NoError(t, err)
		assert.Equal(t, expectedEmployee, employee)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("error employee not found", func(t *testing.T) {
		// Arrange
		id := 99

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"})

		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(rows)

		// Act
		employee, err := repo.GetEmployeeById(id)

		// Assert

		expectedEmployees := models.Employee{}
		expectedError := custom_errors.ErrNotFound

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, expectedEmployees, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		// Arrange
		id := 1

		mock.ExpectQuery(`SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?`).
			WithArgs(id).
			WillReturnError(sql.ErrConnDone)

		// Act
		employee, err := repo.GetEmployeeById(id)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Equal(t, models.Employee{}, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateEmployee(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	t.Run("Success create employee", func(t *testing.T) {
		// Arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP002", FirstName: "Jhon", LastName: "Doe", WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)")).
			WithArgs("EMP002", "Jhon", "Doe", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.CreateEmployee(attributes)

		// Assert
		expectedEmployee := models.Employee{Id: 1, EmployeeAttributes: attributes}

		require.NoError(t, err)
		assert.Equal(t, expectedEmployee, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create employee error duplicated card number id", func(t *testing.T) {
		// Arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP001", FirstName: "Jhon", LastName: "Doe", WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)")).
			WithArgs("EMP001", "Jhon", "Doe", 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Error Code: 1062. Duplicate entry 'EMP001' for key 'employees.card_number_id'",
			})

		// Act
		employee, err := repo.CreateEmployee(attributes)

		// Assert

		expectedEmployees := models.Employee{}
		expectedError := &custom_errors.UniqueAttributeViolationErr{AttributeName: "card_number_id", Value: "EMP001"}
		require.Error(t, err)
		require.Equal(t, expectedEmployees, employee)
		require.Equal(t, err, expectedError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create employee error warehouse id not found", func(t *testing.T) {
		// Arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP002", FirstName: "Jhon", LastName: "Doe", WarehouseId: 99}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)")).
			WithArgs("EMP002", "Jhon", "Doe", 99).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`employees`, CONSTRAINT `employees_ibfk_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))",
			})

		//act
		employee, err := repo.CreateEmployee(attributes)

		//assert
		expectedEmployees := models.Employee{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "warehouse_id", Details: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`employees`, CONSTRAINT `employees_ibfk_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))"}
		require.Error(t, err)
		require.Equal(t, expectedEmployees, employee)
		require.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create employee database error", func(t *testing.T) {
		// Arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP002", FirstName: "Jhon", LastName: "Doe", WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)")).
			WithArgs("EMP002", "Jhon", "Doe", 1).
			WillReturnError(sql.ErrConnDone)

		// Act
		result, err := repo.CreateEmployee(attributes)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Equal(t, models.Employee{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestUpdateEmployee(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	id := 1

	var employeeAttributes = models.EmployeeAttributes{
		CardNumberId: "12345",
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseId:  2,
	}

	fields := []struct {
		name  string
		patch models.EmployeePatchRequestBody
		query string
		args  []driver.Value
	}{
		{
			name: "card_number_id, first_name, last_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.FirstName, employeeAttributes.LastName, employeeAttributes.WarehouseId, id},
		}, {
			name: "first_name, last_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `first_name` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.FirstName, employeeAttributes.LastName, employeeAttributes.WarehouseId, id},
		},
		{
			name: "card_number_id, last_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.LastName, employeeAttributes.WarehouseId, id},
		},
		{
			name: "card_number_id, first_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.FirstName, employeeAttributes.WarehouseId, id},
		},
		{
			name: "card_number_id, first_name, last_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `last_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.FirstName, employeeAttributes.LastName, id},
		}, {
			name: "last_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `last_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.LastName, employeeAttributes.WarehouseId, id},
		}, {
			name: "first_name, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `first_name` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.FirstName, employeeAttributes.WarehouseId, id},
		}, {
			name: "first_name, last_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `first_name` = ?, `last_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.FirstName, employeeAttributes.LastName, id},
		}, {
			name: "card_number_id, warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.WarehouseId, id},
		}, {
			name: "card_number_id, last_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `last_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.LastName, id},
		}, {
			name: "card_number_id, first_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `card_number_id` = ?, `first_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, employeeAttributes.FirstName, id},
		}, {
			name: "card_number_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `card_number_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.CardNumberId, id},
		}, {
			name: "first_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `first_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.FirstName, id},
		}, {
			name: "last_name",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
			query: "UPDATE employees SET `last_name` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.LastName, id},
		}, {
			name: " warehouse_id",
			patch: models.EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
			query: "UPDATE employees SET `warehouse_id` = ? WHERE id = ?",
			args:  []driver.Value{employeeAttributes.WarehouseId, id},
		},
	}

	for _, field := range fields {
		t.Run("Succes - should return nil if "+field.name+" is/are filled", func(t *testing.T) {

			rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
				AddRow(id, 12345, "John", "Doe", 2)

			// First expectation: UPDATE query ( _, err = sql_utils.Update(r.db, query, args) )
			mock.ExpectExec(regexp.QuoteMeta(field.query)).
				WithArgs(field.args...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			// Second expectation: SELECT query (updatedEmployee, err = r.GetEmployeeById(id))
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?")).
				WithArgs(id).
				WillReturnRows(rows)

			// Act
			employee, err := repo.UpdateEmployeeById(id, field.patch)

			// Assert
			expectedEmployee := models.Employee{Id: 1, EmployeeAttributes: employeeAttributes}

			require.NoError(t, err)
			assert.Equal(t, expectedEmployee, employee)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("Error employee not found", func(t *testing.T) {
		// Arrange

		patchRequest := models.EmployeePatchRequestBody{
			CardNumberId: &employeeAttributes.CardNumberId,
			FirstName:    &employeeAttributes.FirstName,
			LastName:     &employeeAttributes.LastName,
			WarehouseId:  &employeeAttributes.WarehouseId,
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?")).
			WithArgs("12345", "John", "Doe", 2, 99).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(sql.ErrNoRows)

		// Act
		employee, err := repo.UpdateEmployeeById(99, patchRequest)

		// Assert
		expectedEmployee := models.Employee{}

		require.Error(t, err)
		assert.Equal(t, expectedEmployee, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Update error warehouse id not found", func(t *testing.T) {

		// Arrange

		warehouseId := 99

		patchRequest := models.EmployeePatchRequestBody{
			CardNumberId: &employeeAttributes.CardNumberId,
			FirstName:    &employeeAttributes.FirstName,
			LastName:     &employeeAttributes.LastName,
			WarehouseId:  &warehouseId,
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?")).
			WithArgs("12345", "John", "Doe", 99, id).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`employees`, CONSTRAINT `employees_ibfk_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))",
			})

		//act
		employee, err := repo.UpdateEmployeeById(id, patchRequest)

		//assert
		expectedEmployees := models.Employee{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "warehouse_id", Details: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`employees`, CONSTRAINT `employees_ibfk_1` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))"}
		require.Error(t, err)
		require.Equal(t, expectedEmployees, employee)
		require.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("Update database error", func(t *testing.T) {
		// Arrange
		id := 1

		patchRequest := models.EmployeePatchRequestBody{
			CardNumberId: &employeeAttributes.CardNumberId,
			FirstName:    &employeeAttributes.FirstName,
			LastName:     &employeeAttributes.LastName,
			WarehouseId:  &employeeAttributes.WarehouseId,
		}

		mock.ExpectExec(regexp.QuoteMeta("UPDATE employees SET `card_number_id` = ?, `first_name` = ?, `last_name` = ?, `warehouse_id` = ? WHERE id = ?")).
			WithArgs("12345", "John", "Doe", 2, id).
			WillReturnError(sql.ErrConnDone)

		// Act
		employee, err := repo.UpdateEmployeeById(id, patchRequest)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Equal(t, models.Employee{}, employee)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDeleteEmployee(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	t.Run("Delete employee success", func(t *testing.T) {
		// Arrange
		id := 1
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `employees` where `id` = ?")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err := repo.DeleteEmployee(id)

		// Assert
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete employee not found", func(t *testing.T) {
		// Arrange
		id := 99
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `employees` where `id` = ?")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		err := repo.DeleteEmployee(id)

		// Assert
		expectedError := custom_errors.ErrNotFound
		require.Error(t, err)
		require.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete employee database error", func(t *testing.T) {
		// Arrange
		id := 1
		mock.ExpectExec("DELETE FROM `employees` where `id` = ?").
			WithArgs(id).
			WillReturnError(sql.ErrConnDone)

		// Act
		err := repo.DeleteEmployee(id)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestGetReportInboundOrders(t *testing.T) {

	repo, mock, cleanup := setupEmployeeRepository(t)
	defer cleanup()

	t.Run("Success get report with id", func(t *testing.T) {
		// Arrange
		employee_id := 1

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id", "inbound_orders_count",
		}).AddRow(1, "12345", "Raul", "García", 1, 1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id WHERE e.id = ? GROUP BY e.id;")).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		reports, err := repo.GetReportInboundOrders(&employee_id)

		// Assert

		expectedReports := []models.InboundOrderEmployee{{

			Employee: models.Employee{
				Id: 1,
				EmployeeAttributes: models.EmployeeAttributes{

					CardNumberId: "12345",
					FirstName:    "Raul",
					LastName:     "García",
					WarehouseId:  1,
				},
			},
			InboundOrdersCount: 1,
		}}

		require.NoError(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success get report empty report", func(t *testing.T) {
		// Arrange

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id", "inbound_orders_count"}).
			AddRow(1, "12345", "Raul", "García", 1, 1).
			AddRow(2, "67890", "Pepe", "Sierra", 2, 3)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id GROUP BY e.id;")).
			WillReturnRows(rows)

		// Act
		reports, err := repo.GetReportInboundOrders(nil)

		// Assert

		expectedReports := []models.InboundOrderEmployee{{

			Employee: models.Employee{
				Id: 1,
				EmployeeAttributes: models.EmployeeAttributes{

					CardNumberId: "12345",
					FirstName:    "Raul",
					LastName:     "García",
					WarehouseId:  1,
				},
			},
			InboundOrdersCount: 1,
		}, {

			Employee: models.Employee{
				Id: 2,
				EmployeeAttributes: models.EmployeeAttributes{

					CardNumberId: "67890",
					FirstName:    "Pepe",
					LastName:     "Sierra",
					WarehouseId:  2,
				},
			},
			InboundOrdersCount: 3,
		},
		}

		require.NoError(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Employee report Not found id", func(t *testing.T) {

		// Arrange
		employee_id := 1

		rows := sqlmock.NewRows([]string{
			"id", "card_number_id", "first_name", "last_name", "warehouse_id", "inbound_orders_count",
		})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id WHERE e.id = ? GROUP BY e.id;")).
			WithArgs(1).
			WillReturnRows(rows).
			WillReturnError(sql.ErrNoRows)

		// Act
		reports, err := repo.GetReportInboundOrders(&employee_id)

		// Assert

		expectedReports := []models.InboundOrderEmployee(nil)
		expectedError := custom_errors.ErrNotFound

		require.Error(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Employee report database_error", func(t *testing.T) {

		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta("SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id GROUP BY e.id;")).
			WillReturnError(sql.ErrConnDone)

		// Act
		reports, err := repo.GetReportInboundOrders(nil)

		// Assert
		expectedReports := []models.InboundOrderEmployee(nil)

		require.Error(t, err)
		assert.Equal(t, expectedReports, reports)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Empty(t, reports)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
