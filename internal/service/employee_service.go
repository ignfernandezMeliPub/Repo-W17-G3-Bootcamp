package service

import (
	employee_repository "app/internal/repository/employee_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"errors"
)

type EmployeeServiceInterface interface {
	GetEmployeesList() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	CreateEmployee(attributes models.EmployeePostRequestBody) (newEmployee models.Employee, err error)
	UpdateEmployee(id int, attributes models.EmployeePatchRequestBody) (employee models.Employee, err error)
	DeleteEmployee(id int) (err error)
}

func NewEmployeeService(repository employee_repository.EmployeeRepository, svWahrehouse WarehouseDefault) *EmployeeService {
	return &EmployeeService{repository: repository, svWahrehouse: svWahrehouse}
}

type EmployeeService struct {
	repository   employee_repository.EmployeeRepository
	svWahrehouse WarehouseDefault
}

func (s *EmployeeService) GetEmployeesList() (employees []models.Employee, err error) {

	employees, err = s.repository.GetEmployeesList()

	return

}

func (s *EmployeeService) GetEmployeeById(id int) (employee models.Employee, err error) {

	employee, err = s.repository.GetEmployeeById(id)
	return

}

func (s *EmployeeService) CreateEmployee(attributes models.EmployeePostRequestBody) (newEmployee models.Employee, err error) {

	err = s.repository.ValidateUniqueCardNumberID(*attributes.CardNumberId)

	if err != nil {

		return

	}

	// add validation wharehouse id

	_, err = s.svWahrehouse.FindWarehouseById(*attributes.WarehouseId)

	if err != nil {

		if errors.As(err, &custom_errors.ErrNotFound) {

			err = &custom_errors.InvalidArgValueErr{

				Argument:  "warehouse_id",
				Value:     *attributes.WarehouseId,
				ExtraInfo: "The warehouse sent doesn't exist",
			}

		}

		return

	}

	newAttributes := models.EmployeeAttributes{
		CardNumberId: *attributes.CardNumberId,
		FirstName:    *attributes.FirstName,
		LastName:     *attributes.LastName,
		WarehouseId:  *attributes.WarehouseId}

	newEmployee, err = s.repository.SaveEmployee(newAttributes)
	return

}

func (s *EmployeeService) UpdateEmployee(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {

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

	if attributes.WarehouseId != nil {

		_, err = s.svWahrehouse.FindWarehouseById(*attributes.WarehouseId)

		if err != nil {

			if errors.As(err, &custom_errors.ErrNotFound) {

				err = &custom_errors.InvalidArgValueErr{

					Argument:  "Warehouse",
					Value:     attributes.WarehouseId,
					ExtraInfo: "The harehouse sent doesn't exist",
				}

			}

			return

		}

	}

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
