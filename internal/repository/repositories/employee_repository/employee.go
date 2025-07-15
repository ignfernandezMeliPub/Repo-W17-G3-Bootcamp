package employee_repository

import "app/pkg/models"

type EmployeeRepository interface {
	GetAllEmployees() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error)
	UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error)
	DeleteEmployee(id int) (err error)

	GetReportInboundOrderByEmployee(id int) (inboundOrder models.InboundOrderEmployee, err error)
	GetReportInboundOrders() (inboundOrders []models.InboundOrderEmployee, err error)
}
