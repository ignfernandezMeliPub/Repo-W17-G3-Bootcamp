package inbound_order_repository

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

func setupInboundOrder(t *testing.T) (InboundOrderRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewInboundOrderDb(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}
func TestCreateInboundOrder(t *testing.T) {

	repo, mock, cleanup := setupInboundOrder(t)
	defer cleanup()

	t.Run("Success create inbound order", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD002", EmployeeId: 1, ProductBatchId: 1, WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD002", 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.CreateInboundOrder(attributes)

		// Assert
		expectedInboundOrder := models.InboundOrder{Id: 1, InboundOrderDetails: attributes}

		require.NoError(t, err)
		require.Equal(t, expectedInboundOrder, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create inbound order error duplicated order number", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD001", EmployeeId: 1, ProductBatchId: 1, WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD001", 1, 1, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Error Code: 1062. Duplicate entry 'ORD001' for key 'inbound_orders.order_number'",
			})

		// Act
		employee, err := repo.CreateInboundOrder(attributes)

		// Assert

		expectedInboundOrders := models.InboundOrder{}
		expectedError := &custom_errors.UniqueAttributeViolationErr{AttributeName: "order_number", Value: "ORD001"}
		require.Error(t, err)
		require.Equal(t, expectedInboundOrders, employee)
		require.Equal(t, err, expectedError)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create inbound order error product employee id not found", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD002", EmployeeId: 99, ProductBatchId: 1, WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD002", 99, 1, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`))",
			})

		//act
		employee, err := repo.CreateInboundOrder(attributes)

		//assert
		expectedInboundOrders := models.InboundOrder{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "employee_id", Details: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_1` FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`))"}
		require.Error(t, err)
		require.Equal(t, expectedInboundOrders, employee)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create inbound order error product batch id not found", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD002", EmployeeId: 1, ProductBatchId: 99, WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD002", 1, 99, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_2` FOREIGN KEY (`product_batch_id`) REFERENCES `product_batch` (`id`))",
			})

		//act
		employee, err := repo.CreateInboundOrder(attributes)

		//assert
		expectedInboundOrders := models.InboundOrder{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "product_batch_id", Details: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_2` FOREIGN KEY (`product_batch_id`) REFERENCES `product_batch` (`id`))"}
		require.Error(t, err)
		require.Equal(t, expectedInboundOrders, employee)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create inbound order error warehouse id not found", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD002", EmployeeId: 1, ProductBatchId: 1, WarehouseId: 99}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD002", 1, 1, 99).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_3` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))",
			})

		//act
		employee, err := repo.CreateInboundOrder(attributes)

		//assert
		expectedInboundOrders := models.InboundOrder{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "warehouse_id", Details: "Error Code: 1452. Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`inbound_orders`, CONSTRAINT `inbound_orders_ibfk_3` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`))"}
		require.Error(t, err)
		require.Equal(t, expectedInboundOrders, employee)
		require.Equal(t, expectedError, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Create employee database error", func(t *testing.T) {
		// Arrange
		attributes := models.InboundOrderDetails{OrderDate: "2021-04-04", OrderNumber: "ORD002", EmployeeId: 1, ProductBatchId: 1, WarehouseId: 1}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inbound_orders` (`order_date`,`order_number`,`employee_id`,`product_batch_id`,`warehouse_id`) VALUES (?,?,?,?,?)")).
			WithArgs("2021-04-04", "ORD002", 1, 1, 1).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err := repo.CreateInboundOrder(attributes)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
