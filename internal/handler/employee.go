package handler

import (
	"app/internal/service"
	"app/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

type EmployeesController struct {
	sv service.EmployeeServiceInterface
}

func NewEmployeeController(sv service.EmployeeServiceInterface) *EmployeesController {
	return &EmployeesController{sv: sv}
}

func (c *EmployeesController) GetEmployeesList(w http.ResponseWriter, r *http.Request) {

	res, err := c.sv.GetEmployeesList()
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return

}

func (c *EmployeesController) GetEmployeeById(w http.ResponseWriter, r *http.Request) {

	id, idError := strconv.Atoi(chi.URLParam(r, "id"))

	if idError != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "el formto del id debe ser un entero",
			"error":   idError.Error(),
		})
		return

	}

	employee, err := c.sv.GetEmployeeById(id)

	if err != nil {

		response.JSON(w, http.StatusNotFound, map[string]any{
			"data": err.Error(),
		})
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) SaveEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeAttributes models.EmployeeRequestBody

	if err := json.NewDecoder(r.Body).Decode(&newEmployeeAttributes); err != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "Some of the fields sent have the wrong type",
			"data":    err.Error(),
		})
		return

	}

	if validationError := newEmployeeAttributes.VerifyMandatoryFieldsPresence(); validationError != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "A field is missing",
			"data":    validationError.Error(),
		})
		return

	}

	employee, err := c.sv.CreateEmployee(newEmployeeAttributes)

	if err != nil {

		response.JSON(w, http.StatusInternalServerError, map[string]any{
			"data": err.Error(),
		})
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {

	var newEmployeeAttributes models.EmployeeRequestBody

	id, idError := strconv.Atoi(chi.URLParam(r, "id"))

	if idError != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "el formto del id debe ser un entero",
			"error":   idError.Error(),
		})
		return

	}

	if err := json.NewDecoder(r.Body).Decode(&newEmployeeAttributes); err != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "Some of the fields sent have the wrong type",
			"data":    err.Error(),
		})
		return

	}

	employee, err := c.sv.UpdateEmployee(id, newEmployeeAttributes)

	if err != nil {

		response.JSON(w, http.StatusInternalServerError, map[string]any{
			"data": err.Error(),
		})
		return

	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})

}

func (c *EmployeesController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	id, idError := strconv.Atoi(chi.URLParam(r, "id"))

	if idError != nil {

		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "el formto del id debe ser un entero",
			"error":   idError.Error(),
		})
		return

	}

	err := c.sv.DeleteEmployee(id)

	if err != nil {

		response.JSON(w, http.StatusInternalServerError, map[string]any{
			"data": err.Error(),
		})
		return

	}

	response.JSON(w, http.StatusNoContent, map[string]any{
		"message": "Deleted succesfully",
	})

}
