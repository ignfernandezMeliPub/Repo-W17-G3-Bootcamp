package employee_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

func NewEmployeeMap(db map[int]models.Employee) *EmployeeMap {
	// default db
	defaultDb := make(map[int]models.Employee)
	if db != nil {
		defaultDb = db
	}
	return &EmployeeMap{db: defaultDb, lastId: len(db)}
}

type EmployeeMap struct {
	db     map[int]models.Employee
	lastId int
}

func (m *EmployeeMap) GetEmployeesList() (employees []models.Employee, err error) {

	employees = make([]models.Employee, len(m.db))

	i := 0
	for _, value := range m.db {
		employees[i] = value
		i++
	}

	if len(employees) == 0 {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}

	return

}

func (m *EmployeeMap) GetEmployeeById(id int) (employee models.Employee, err error) {

	employee, ok := m.db[id]

	if !ok {

		err = &custom_errors.ResourceNotFoundError{}

	}

	return

}

func (m *EmployeeMap) ValidateUniqueCardNumberID(cardNumber string) (err error) {

	for _, employee := range m.db {

		if employee.CardNumberId == cardNumber {

			err = &custom_errors.UniqueAttributeViolationErr{AttributeName: "Card_number_id", Value: cardNumber}
			return
		}

	}

	return

}

func (m *EmployeeMap) SaveEmployee(attributes models.EmployeeAttributes) (newEmployee models.Employee, err error) {

	m.lastId = m.lastId + 1
	newEmployee = models.Employee{Id: m.lastId,
		EmployeeAttributes: models.EmployeeAttributes{

			CardNumberId: attributes.CardNumberId,
			FirstName:    attributes.FirstName,
			LastName:     attributes.LastName,
			WarehouseId:  attributes.WarehouseId,
		}}

	m.db[m.lastId] = newEmployee

	return

}

func (m *EmployeeMap) UpdateEmployee(employee models.Employee) (updatedEmployee models.Employee, err error) {

	updatedEmployee = employee
	m.db[employee.Id] = updatedEmployee

	return

}

func (m *EmployeeMap) DeleteEmployee(id int) (err error) {

	delete(m.db, id)

	return

}
