package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
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

func (c *EmployeesController) GetAllEmployees(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetAllEmployees", logger.LogStatusInProgress)

	res, err := c.svEmployee.GetAllEmployees()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllEmployees", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetAllEmployees", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return nil
}

func (c *EmployeesController) GetEmployeeById(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetEmployeeById", logger.LogStatusInProgress)

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {
		httpStatus := utils.ResponseHttpError(w, idError)
		utils.LogError(r, "GetEmployeeById", idError, httpStatus)
		return nil
	}

	employee, err := c.svEmployee.GetEmployeeById(id)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetEmployeeById", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetEmployeeById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})
	return nil

}

func (c *EmployeesController) CreateEmployee(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateEmployee", logger.LogStatusInProgress)

	var newEmployeeAttributes models.EmployeePostRequestBody

	newEmployeeAttributes, err = utils.InstantiateVarFromBody(&r.Body, newEmployeeAttributes)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateEmployee", err, httpStatus)
		return nil
	}

	employee, err := c.svEmployee.CreateEmployee(newEmployeeAttributes)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateEmployee", err, httpStatus)
		return nil
	}

	utils.Log(r, "CreateEmployee", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": employee,
	})
	return nil

}

func (c *EmployeesController) PatchEmployee(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "PatchEmployee", logger.LogStatusInProgress)

	var newEmployeeAttributes models.EmployeePatchRequestBody

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {
		httpStatus := utils.ResponseHttpError(w, idError)
		utils.LogError(r, "PatchEmployee", idError, httpStatus)
		return nil
	}

	newEmployeeAttributes, err = utils.InstantiateVarFromBody(&r.Body, newEmployeeAttributes)

	// Some of the fields sent have the wrong type
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchEmployee", err, httpStatus)
		return nil
	}

	employee, err := c.svEmployee.UpdateEmployeeById(id, newEmployeeAttributes)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchEmployee", err, httpStatus)
		return nil
	}

	utils.Log(r, "PatchEmployee", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": employee,
	})
	return nil
}

func (c *EmployeesController) DeleteEmployee(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "DeleteEmployee", logger.LogStatusInProgress)

	id, idError := utils.GetURLParamAs(r, "id", strconv.Atoi)

	// id format invalid
	if idError != nil {
		httpStatus := utils.ResponseHttpError(w, idError)
		utils.LogError(r, "DeleteEmployee", idError, httpStatus)
		return nil
	}

	err = c.svEmployee.DeleteEmployee(id)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteEmployee", err, httpStatus)
		return nil
	}

	utils.Log(r, "DeleteEmployee", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, map[string]any{
		"message": "Deleted succesfully",
	})
	return nil
}

func (c *EmployeesController) GetReportInboundOrders(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetReportInboundOrders", logger.LogStatusInProgress)

	id, idError := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if idError != nil {
		httpStatus := utils.ResponseHttpError(w, idError)
		utils.LogError(r, "GetReportInboundOrders", idError, httpStatus)
		return nil
	}

	res, err := c.svEmployee.GetReportInboundOrders(id)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetReportInboundOrders", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetReportInboundOrders", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return nil
}
