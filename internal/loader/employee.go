package loader

import "app/pkg/models"

func NewEmployeeJSONFile(path string) *EmployeeJSONFile {
	return &EmployeeJSONFile{
		path: path,
	}
}

// VehicleJSONFile is a struct that implements the LoaderVehicle interface
type EmployeeJSONFile struct {
	// path is the path to the file that contains the vehicles in JSON format
	path string
}

func (l *EmployeeJSONFile) Load() (employees map[int]models.Employee, err error) {

	employeList, errList := LoadDataFromFile[models.Employee](l.path)

	if errList != nil {

		err = errList
		return

	}

	employees = make(map[int]models.Employee)
	for _, employee := range employeList {
		employees[employee.Id] = models.Employee{
			Id: employee.Id,
			EmployeeAttributes: models.EmployeeAttributes{

				CardNumberId: employee.CardNumberId,
				FirstName:    employee.FirstName,
				LastName:     employee.LastName,
				WarehouseId:  employee.WarehouseId,
			},
		}

	}

	return

}
