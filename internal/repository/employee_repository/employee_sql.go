package employee_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
)

func NewEmployeeDb(db *sql.DB) EmployeeRepository {

	return &EmployeeDb{db: db}

}

type EmployeeDb struct {
	db *sql.DB
}

func (m *EmployeeDb) GetAllEmployees() (employees []models.Employee, err error) {

	employees, err = sql_utils.Query[models.Employee](m.db, "SELECT Id, CardNumberId, FirstName, LastName, WarehouseId FROM employees", nil)

	err = sql_utils.HandleSqlError(err)

	return

}

func (m *EmployeeDb) GetEmployeeById(id int) (employee models.Employee, err error) {

	args := make([]any, 1)
	args[0] = id

	employee, err = sql_utils.QueryRow[models.Employee](m.db, "SELECT Id, CardNumberId, FirstName, LastName, WarehouseId FROM employees WHERE Id = ?", args)

	err = sql_utils.HandleSqlError(err)

	return

}

func (m *EmployeeDb) CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error) {

	args := []any{attributes.CardNumberId, attributes.FirstName, attributes.LastName, attributes.WarehouseId}

	lastId, err := sql_utils.Insert(m.db, "INSERT INTO `employees` (`CardNumberId`,`FirstName`,`LastName`,`WarehouseId`) VALUES (?,?,?,?)", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		return
	}

	newEmployee = models.Employee{

		Id:                 int(lastId),
		EmployeeAttributes: attributes,
	}

	return

}

func (m *EmployeeDb) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {

	query := "UPDATE employees SET "
	var args []any

	if attributes.CardNumberId != nil {
		query += "`CardNumberId` = ?, "
		args = append(args, *attributes.CardNumberId)
	}
	if attributes.FirstName != nil {
		query += "`FirstName` = ?, "
		args = append(args, *attributes.FirstName)
	}
	if attributes.LastName != nil {
		query += "`LastName` = ?, "
		args = append(args, *attributes.LastName)
	}
	if attributes.WarehouseId != nil {
		query += "`WarehouseId` = ?, "
		args = append(args, *attributes.WarehouseId)
	}

	query = query[:len(query)-2]
	query += " WHERE id = ?"
	args = append(args, id)

	rowsAffected, err := sql_utils.Update(m.db, query, args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		return
	}

	if rowsAffected == 0 {
		return models.Employee{}, custom_errors.ErrNotFound
	}

	updatedEmployee, err = m.GetEmployeeById(id)

	return

}

func (m *EmployeeDb) DeleteEmployee(id int) (err error) {

	args := make([]any, 1)
	args[0] = id

	_, err = sql_utils.Delete(m.db, "DELETE FROM `employees` where `id` = ?", args)

	err = sql_utils.HandleSqlError(err)
	return

}
