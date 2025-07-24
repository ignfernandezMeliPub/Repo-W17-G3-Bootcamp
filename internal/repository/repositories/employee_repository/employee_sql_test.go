package employee_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-txdb"
)

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "fresh_db_test",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestSectionsServiceImpl_GetAllEmployees(t *testing.T) {

	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewEmployeeDb(db)

	t.Run("Success get all employees", func(t *testing.T) {

		//arrange

		//act
		employees, err := rp.GetAllEmployees()

		//assert

		expectedEmployees := []models.Employee{{Id: 1, EmployeeAttributes: models.EmployeeAttributes{CardNumberId: "EMP001", FirstName: "Raul", LastName: "García", WarehouseId: 1}}}

		require.Equal(t, expectedEmployees, employees)
		require.NoError(t, err)
		require.NotNil(t, employees)
		require.Greater(t, len(employees), 0)
	})

	t.Run("Error get all not found employees", func(t *testing.T) {

		//arrange
		err := rp.DeleteEmployee(1)
		require.NoError(t, err)

		//act
		employees, err := rp.GetAllEmployees()

		//assert
		expectedEmployees := []models.Employee(nil)
		require.Equal(t, expectedEmployees, employees)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
		require.Nil(t, employees)
	})
}

func TestGetEmployeeById(t *testing.T) {

	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewEmployeeDb(db)

	t.Run("Success get employee", func(t *testing.T) {

		//arrange
		id := 1

		//act
		employee, err := rp.GetEmployeeById(id)

		//assert

		expectedEmployees := models.Employee{Id: 1, EmployeeAttributes: models.EmployeeAttributes{CardNumberId: "EMP001", FirstName: "Raul", LastName: "García", WarehouseId: 1}}

		require.Equal(t, expectedEmployees, employee)
		require.NoError(t, err)
		require.NotNil(t, employee)
	})

	t.Run("Error employee not found", func(t *testing.T) {

		//arrange
		id := 5

		//act
		employee, err := rp.GetEmployeeById(id)

		//assert
		expectedEmployees := models.Employee{}
		require.Equal(t, expectedEmployees, employee)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})

	})

}

func TestCreateEmployee(t *testing.T) {

	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewEmployeeDb(db)

	t.Run("Success create employee", func(t *testing.T) {

		//arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP002", FirstName: "Pepe", LastName: "Doe", WarehouseId: 1}

		//act
		employee, err := rp.CreateEmployee(attributes)

		//assert

		expectedEmployee := models.Employee{Id: 2, EmployeeAttributes: attributes}

		require.Equal(t, expectedEmployee, employee)
		require.NoError(t, err)
		require.NotNil(t, employee)
	})

	t.Run("Error duplicated card number id", func(t *testing.T) {

		//arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP001", FirstName: "Pepe", LastName: "Doe", WarehouseId: 1}

		//act
		employee, err := rp.CreateEmployee(attributes)

		//assert
		expectedEmployees := models.Employee{}
		expectedError := &custom_errors.UniqueAttributeViolationErr{AttributeName: "card_number_id", Value: "EMP001"}
		require.Equal(t, expectedEmployees, employee)
		require.Equal(t, err, expectedError)

	})

	t.Run("Error warehouse id not found", func(t *testing.T) {

		//arrange
		attributes := models.EmployeeAttributes{CardNumberId: "EMP002", FirstName: "Pepe", LastName: "Doe", WarehouseId: 99}

		//act
		employee, err := rp.CreateEmployee(attributes)

		//assert
		expectedEmployees := models.Employee{}
		expectedError := &custom_errors.ForeignKeyViolationError{ConstraintName: "warehouse_id"}
		require.Equal(t, expectedEmployees, employee)
		require.IsType(t, expectedError, err)

	})

}

func TestUpdateEmployee(t *testing.T) {

	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewEmployeeDb(db)

	employeeAttributes := models.EmployeeAttributes{

		CardNumberId: "EMP002", FirstName: "Pepe", LastName: "Doe", WarehouseId: 1,
	}

	t.Run("Success updated employee", func(t *testing.T) {

		//arrange
		id := 1
		attributesPatch := models.EmployeePatchRequestBody{
			CardNumberId: &employeeAttributes.CardNumberId,
			FirstName:    &employeeAttributes.FirstName,
			LastName:     &employeeAttributes.LastName,
			WarehouseId:  &employeeAttributes.WarehouseId,
		}

		//act
		employee, err := rp.UpdateEmployeeById(id, attributesPatch)

		//assert
		expectedEmployee := models.Employee{Id: 1, EmployeeAttributes: employeeAttributes}

		require.Equal(t, expectedEmployee, employee)
		require.NoError(t, err)
		require.NotNil(t, employee)
	})

}
