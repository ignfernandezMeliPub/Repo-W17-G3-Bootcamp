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

func (c *EmployeesController) GetAllEmployees(w http.ResponseWriter, _ *http.Request) {

	res, err := c.svEmployee.GetAllEmployees()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
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

func (c *EmployeesController) CreateEmployee(w http.ResponseWriter, r *http.Request) {

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

func (c *EmployeesController) PatchEmployee(w http.ResponseWriter, r *http.Request) {

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

	employee, err := c.svEmployee.UpdateEmployeeById(id, newEmployeeAttributes)

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

func (c *EmployeesController) GetReportInboundOrders(w http.ResponseWriter, r *http.Request) {

	id, idError := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if id != nil {

		// id format invalid
		if idError != nil {

			utils.ResponseHttpError(w, idError)
			return

		}

		inboundOrder, err := c.svEmployee.GetReportInboundOrderByEmployee(*id)

		if err != nil {

			utils.ResponseHttpError(w, err)
			return

		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": inboundOrder,
		})

	} else {

		res, err := c.svEmployee.GetReportInboundOrders()

		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{"data": res})

	}

}
