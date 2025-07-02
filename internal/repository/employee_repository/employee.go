package employee_repository

import "app/pkg/models"

type EmployeeRepository interface {
	GetEmployeesList() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	ValidateUniqueCardNumberID(cardNumber string) (err error)
	SaveEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error)
	UpdateEmployee(employee models.Employee) (updatedEmployee models.Employee, err error)
	DeleteEmployee(id int) (err error)
}
