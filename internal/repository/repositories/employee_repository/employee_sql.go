package employee_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

func NewEmployeeDb(db *sql.DB) EmployeeRepository {

	return &EmployeeDb{db: db}

}

type EmployeeDb struct {
	db *sql.DB
}

func (r *EmployeeDb) GetAllEmployees() (employees []models.Employee, err error) {
	sql_utils.Log("GetAllEmployees", logger.LogStatusInProgress, "Select all employees")

	employees, err = sql_utils.Query[models.Employee](r.db, "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees", nil)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetAllEmployees", "Select all employees", err)
		return
	}

	sql_utils.Log("GetAllEmployees", logger.LogStatusSuccess, "Select all employees")
	return

}

func (r *EmployeeDb) GetEmployeeById(id int) (employee models.Employee, err error) {
	sql_utils.Log("GetEmployeeById", logger.LogStatusInProgress, "Select employee by id "+strconv.Itoa(id))

	args := make([]any, 1)
	args[0] = id

	employee, err = sql_utils.QueryRow[models.Employee](r.db, "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetEmployeeById", "Select employee by id "+strconv.Itoa(id), err)
		return
	}

	sql_utils.Log("GetEmployeeById", logger.LogStatusSuccess, "Select employee by id "+strconv.Itoa(id))
	return

}

func (r *EmployeeDb) CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error) {
	sql_utils.LogAudit("CreateEmployee", logger.LogStatusInProgress, "Insert employee")

	args := []any{attributes.CardNumberId, attributes.FirstName, attributes.LastName, attributes.WarehouseId}

	lastId, err := sql_utils.Insert(r.db, "INSERT INTO `employees` (`card_number_id`,`first_name`,`last_name`,`warehouse_id`) VALUES (?,?,?,?)", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateEmployee", "Insert employee", err)
		return
	}

	newEmployee = models.Employee{

		Id:                 int(lastId),
		EmployeeAttributes: attributes,
	}

	sql_utils.LogAudit("CreateEmployee", logger.LogStatusSuccess, "Insert employee. Id: "+strconv.Itoa(newEmployee.Id))
	return

}

func (r *EmployeeDb) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {
	sql_utils.LogAudit("UpdateEmployeeById", logger.LogStatusInProgress, "Update employee by id: "+strconv.Itoa(id))

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

	_, err = sql_utils.Update(r.db, query, args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("UpdateEmployeeById", "Update employee by id: "+strconv.Itoa(id), err)
		return
	}

	updatedEmployee, err = r.GetEmployeeById(id)

	if err == nil {
		sql_utils.LogAudit("UpdateEmployeeById", logger.LogStatusSuccess, "Update employee by id: "+strconv.Itoa(id))
	}

	return

}

func (r *EmployeeDb) DeleteEmployee(id int) (err error) {
	sql_utils.LogAudit("DeleteEmployee", logger.LogStatusInProgress, "Delete employee by id: "+strconv.Itoa(id))

	args := make([]any, 1)
	args[0] = id

	rowsAffected, err := sql_utils.Delete(r.db, "DELETE FROM `employees` where `id` = ?", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteEmployee", "Delete employee by id: "+strconv.Itoa(id), err)
		return
	}

	if rowsAffected == 0 {
		err = custom_errors.ErrNotFound
		sql_utils.LogAuditError("DeleteEmployee", "Delete employee by id: "+strconv.Itoa(id), err)
		return
	}

	sql_utils.LogAudit("DeleteEmployee", logger.LogStatusSuccess, "Delete employee by id: "+strconv.Itoa(id))
	return

}

func (r *EmployeeDb) GetReportInboundOrders(id *int) (inboundOrders []models.InboundOrderEmployee, err error) {

	if id != nil {
		sql_utils.Log("GetReportInboundOrders", logger.LogStatusInProgress, "Select employee inbound orders report by id: "+strconv.Itoa(*id))
	} else {
		sql_utils.Log("GetReportInboundOrders", logger.LogStatusInProgress, "Select employee inbound orders report")
	}

	query := "SELECT e.id as id, e.card_number_id as card_number_id, e.first_name as first_name, e.last_name as last_name, e.warehouse_id as warehouse_id, COUNT(io.id) as inbound_orders_count FROM employees e LEFT JOIN inbound_orders io ON e.id = io.employee_id"

	if id != nil {

		query += " WHERE e.id = ? GROUP BY e.id;"
		inboundOrders, err = sql_utils.Query[models.InboundOrderEmployee](r.db, query, []any{*id})

	} else {

		query += " GROUP BY e.id;"
		inboundOrders, err = sql_utils.Query[models.InboundOrderEmployee](r.db, query, nil)

	}

	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetReportInboundOrders", "Select employee inbound orders report", err)
		return
	}

	if id != nil {
		sql_utils.Log("GetReportInboundOrders", logger.LogStatusSuccess, "Select employee inbound orders report by id: "+strconv.Itoa(*id))
	} else {
		sql_utils.Log("GetReportInboundOrders", logger.LogStatusSuccess, "Select employee inbound orders report")
	}

	return

}
