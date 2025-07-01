package service

import (
	employee_repository "app/internal/repository/employee_repository"
	"app/pkg/models"
)

type EmployeeServiceInterface interface {
	GetEmployeesList() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	CreateEmployee(attributes models.EmployeeRequestBody) (newEmployee models.Employee, err error)
	UpdateEmployee(id int, attributes models.EmployeeRequestBody) (employee models.Employee, err error)
	DeleteEmployee(id int) (err error)
}

func NewEmployeeService(repository employee_repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repository: repository}
}

type EmployeeService struct {
	repository employee_repository.EmployeeRepository
}

func (s *EmployeeService) GetEmployeesList() (employees []models.Employee, err error) {

	employees, err = s.repository.GetEmployeesList()

	return

}

func (s *EmployeeService) GetEmployeeById(id int) (employee models.Employee, err error) {

	employee, err = s.repository.GetEmployeeById(id)
	return

}

func (s *EmployeeService) CreateEmployee(attributes models.EmployeeRequestBody) (newEmployee models.Employee, err error) {

	err = s.repository.ValidateUniqueCardNumberID(*attributes.CardNumberId)

	if err != nil {

		return

	}

	// add validation wharehouse id

	newAttributes := models.EmployeeAttributes{
		CardNumberId: *attributes.CardNumberId,
		FirstName:    *attributes.FirstName,
		LastName:     *attributes.LastName,
		WarehouseId:  *attributes.WarehouseId}

	newEmployee, err = s.repository.SaveEmployee(newAttributes)
	return

}

func (s *EmployeeService) UpdateEmployee(id int, attributes models.EmployeeRequestBody) (updatedEmployee models.Employee, err error) {

	employee, err := s.repository.GetEmployeeById(id)

	if err != nil {

		return

	}

	if attributes.CardNumberId != nil {

		err = s.repository.ValidateUniqueCardNumberID(*attributes.CardNumberId)

		if err != nil {

			return

		}

	}

	// add validation wharehouse id

	employee.Patch(attributes)
	updatedEmployee, err = s.repository.UpdateEmployee(employee)

	return

}

func (s *EmployeeService) DeleteEmployee(id int) (err error) {

	_, err = s.repository.GetEmployeeById(id)

	if err != nil {

		return

	}

	s.repository.DeleteEmployee(id)

	return

}
