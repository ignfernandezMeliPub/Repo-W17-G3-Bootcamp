package service

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) GetAllEmployees() (employees []models.Employee, err error) {

	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)

}

func (m *MockEmployeeService) GetEmployeeById(id int) (employee models.Employee, err error) {
	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) CreateEmployee(attributes models.EmployeePostRequestBody) (newEmployee models.Employee, err error) {
	args := m.Called(attributes)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {
	args := m.Called(id, attributes)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *MockEmployeeService) DeleteEmployee(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeService) GetReportInboundOrders(id *int) (inboundOrders []models.InboundOrderEmployee, err error) {
	args := m.Called(id)
	return args.Get(0).([]models.InboundOrderEmployee), args.Error(1)
}
