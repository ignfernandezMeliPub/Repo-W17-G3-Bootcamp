package repository

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type EmployeeMock struct {
	mock.Mock
}

func NewEmployeesMock() *EmployeeMock {
	return &EmployeeMock{}
}

func (m *EmployeeMock) GetAllEmployees() (employees []models.Employee, err error) {

	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)

}

func (m *EmployeeMock) GetEmployeeById(id int) (employee models.Employee, err error) {

	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)

}

func (m *EmployeeMock) CreateEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error) {

	args := m.Called(attributes)
	return args.Get(0).(models.Employee), args.Error(1)

}

func (m *EmployeeMock) UpdateEmployeeById(id int, attributes models.EmployeePatchRequestBody) (updatedEmployee models.Employee, err error) {

	args := m.Called(id, attributes)
	return args.Get(0).(models.Employee), args.Error(1)

}

func (m *EmployeeMock) DeleteEmployee(id int) (err error) {

	args := m.Called(id)
	return args.Error(0)

}

func (m *EmployeeMock) GetReportInboundOrders(id *int) (inboundOrders []models.InboundOrderEmployee, err error) {

	args := m.Called(id)
	return args.Get(0).([]models.InboundOrderEmployee), args.Error(1)

}
