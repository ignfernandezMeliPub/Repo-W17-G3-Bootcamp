package employee_repository

import "app/pkg/models"

type EmployeeRepository interface {
	GetAllEmployees() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	ValidateUniqueCardNumberID(cardNumber string) (err error)
	CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error)
	UpdateEmployeeById(employee models.Employee) (updatedEmployee models.Employee, err error)
	DeleteEmployee(id int) (err error)
}
