package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
)

type EmployeesController struct {
	svEmployee service.EmployeeServiceInterface
}

func NewEmployeeController(svEmployee service.EmployeeServiceInterface) *EmployeesController {
	return &EmployeesController{svEmployee: svEmployee}
}

func (c *EmployeesController) GetEmployeesList(w http.ResponseWriter, r *http.Request) {

	res, err := c.svEmployee.GetEmployeesList()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return

}

func (c *EmployeesController) GetEmployeeById(w http.ResponseWriter, r *http.Request) {

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {

		utils.ResponseHttpError(w, idError)
		return

	}

	employee, err := c.svEmployee.GetEmployeeById(id)

	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) SaveEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeAttributes models.EmployeePostRequestBody

	newEmployeeAttributes, err := utils.InstantiateVarFromBody(&r.Body, newEmployeeAttributes)

	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	employee, err := c.svEmployee.CreateEmployee(newEmployeeAttributes)

	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeAttributes models.EmployeePatchRequestBody

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {
		utils.ResponseHttpError(w, idError)
		return
	}

	newEmployeeAttributes, err := utils.InstantiateVarFromBody(&r.Body, newEmployeeAttributes)

	// Some of the fields sent have the wrong type
	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	employee, err := c.svEmployee.UpdateEmployee(id, newEmployeeAttributes)

	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {

		utils.ResponseHttpError(w, idError)
		return

	}

	err := c.svEmployee.DeleteEmployee(id)

	if err != nil {

		utils.ResponseHttpError(w, err)
		return

	}

	response.JSON(w, http.StatusNoContent, map[string]any{
		"message": "Deleted succesfully",
	})

}
