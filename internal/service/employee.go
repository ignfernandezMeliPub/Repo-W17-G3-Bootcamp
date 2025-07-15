package service

import (
	"app/internal/repository/employee_repository"
	"app/pkg/models"
)

type EmployeeServiceInterface interface {
	GetAllEmployees() (employees []models.Employee, err error)
	GetEmployeeById(id int) (employee models.Employee, err error)
	CreateEmployee(attributes models.EmployeePostRequestBody) (newEmployee models.Employee, err error)
	UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (employee models.Employee, err error)
	DeleteEmployee(id int) (err error)

	GetReportInboundOrderByEmployee(id int) (inboundOrder models.InboundOrderEmployee, err error)
	GetReportInboundOrders() (inboundOrders []models.InboundOrderEmployee, err error)
}

func NewEmployeeService(repository employee_repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repository: repository}
}

type EmployeeService struct {
	repository employee_repository.EmployeeRepository
}

func (s *EmployeeService) GetAllEmployees() (employees []models.Employee, err error) {

	employees, err = s.repository.GetAllEmployees()
	return

}

func (s *EmployeeService) GetEmployeeById(id int) (employee models.Employee, err error) {

	employee, err = s.repository.GetEmployeeById(id)
	return

}

func (s *EmployeeService) CreateEmployee(attributes models.EmployeePostRequestBody) (newEmployee models.Employee, err error) {

	newAttributes := models.EmployeeAttributes{
		CardNumberId: *attributes.CardNumberId,
		FirstName:    *attributes.FirstName,
		LastName:     *attributes.LastName,
		WarehouseId:  *attributes.WarehouseId,
	}

	newEmployee, err = s.repository.CreateEmployee(newAttributes)
	return

}

func (s *EmployeeService) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {

	updatedEmployee, err = s.repository.UpdateEmployeeById(id, attributes)
	return

}

func (s *EmployeeService) DeleteEmployee(id int) (err error) {

	err = s.repository.DeleteEmployee(id)

	return
}

func (s *EmployeeService) GetReportInboundOrders() (inboundOrders []models.InboundOrderEmployee, err error) {

	inboundOrders, err = s.repository.GetReportInboundOrders()
	return

}

func (s *EmployeeService) GetReportInboundOrderByEmployee(id int) (inboundOrder models.InboundOrderEmployee, err error) {

	inboundOrder, err = s.repository.GetReportInboundOrderByEmployee(id)
	return

}
