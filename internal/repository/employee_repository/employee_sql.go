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

func (r *EmployeeDb) GetAllEmployees() (employees []models.Employee, err error) {

	employees, err = sql_utils.Query[models.Employee](r.db, "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees", nil)

	if err != nil {

		return nil, err

	}

	if len(employees) == 0 {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}

	return

}

func (r *EmployeeDb) GetEmployeeById(id int) (employee models.Employee, err error) {

	args := make([]any, 1)
	args[0] = id

	employee, err = sql_utils.QueryRow[models.Employee](r.db, "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?", args)

	if err != nil {

		return models.Employee{}, err

	}

	return

}

func (r *EmployeeDb) CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error) {

	args := []any{attributes.CardNumberId, attributes.FirstName, attributes.LastName, attributes.WarehouseId}

	lastId, err := sql_utils.Insert(r.db, "INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)", args)

	if err != nil {

		return models.Employee{}, err

	}

	newEmployee = models.Employee{

		Id:                 int(lastId),
		EmployeeAttributes: attributes,
	}

	return

}

func (r *EmployeeDb) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {

	query := "UPDATE employees SET "
	var args []any

	if attributes.CardNumberId != nil {
		query += "`card_number_id` = ?, "
		args = append(args, *attributes.CardNumberId)
	}
	if attributes.FirstName != nil {
		query += "`first_name` = ?, "
		args = append(args, *attributes.FirstName)
	}
	if attributes.LastName != nil {
		query += "`last_name` = ?, "
		args = append(args, *attributes.LastName)
	}
	if attributes.WarehouseId != nil {
		query += "`warehouse_id` = ?, "
		args = append(args, *attributes.WarehouseId)
	}

	query = query[:len(query)-2]
	query += " WHERE id = ?"
	args = append(args, id)

	rowsAffected, err := sql_utils.Update(r.db, query, args)

	if err != nil {

		return models.Employee{}, err

	}

	if rowsAffected == 0 {

		return models.Employee{}, &custom_errors.ResourceNotFoundError{}

	}

	updatedEmployee, err = r.GetEmployeeById(id)

	return

}

func (r *EmployeeDb) DeleteEmployee(id int) (err error) {

	args := make([]any, 1)
	args[0] = id

	rowsAffected, err := sql_utils.Delete(r.db, "DELETE FROM `employees` where `id` = ?", args)

	if err != nil {

		return

	}

	if rowsAffected == 0 {

		err = &custom_errors.ResourceNotFoundError{}

	}

	return

}

func (r *EmployeeDb) GetReportInboundOrderByEmployee(id int) (inboundOrder models.InboundOrderEmployee, err error) {

	args := make([]any, 1)
	args[0] = id

	inboundOrder, err = sql_utils.QueryRow[models.InboundOrderEmployee](r.db, "SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id WHERE e.id = ? GROUP BY e.id, e.card_number_id, e.first_name, e.last_name, e.warehouse_id;", args)

	if err != nil {

		return models.InboundOrderEmployee{}, err

	}

	return

}

func (r *EmployeeDb) GetReportInboundOrders() (inboundOrders []models.InboundOrderEmployee, err error) {

	inboundOrders, err = sql_utils.Query[models.InboundOrderEmployee](r.db, "SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id GROUP BY e.id;", nil)

	if err != nil {

		return nil, err

	}

	return

}
