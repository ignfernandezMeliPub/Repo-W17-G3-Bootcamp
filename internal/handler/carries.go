package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"net/http"
	"strings"

	"github.com/bootcamp-go/web/response"
)

type CarriesHandler struct {
	sv service.CarriesService
}

func NewCarriesHandler(sv service.CarriesService) *CarriesHandler {
	return &CarriesHandler{sv: sv}
}

func (h *CarriesHandler) CreateCarrie(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateCarrie", logger.LogStatusInProgress)

	var c models.Carries
	c, err = utils.InstantiateVarFromBody(&r.Body, c)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateCarrie", err, httpStatus)
		return nil
	}
	err = validateCarriesAttributes(c)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateCarrie", err, httpStatus)
		return nil
	}
	data, err := h.sv.CreateCarrie(c)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateCarrie", err, httpStatus)
		return nil
	}

	utils.Log(r, "CreateCarrie", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
	return nil
}

func validateCarriesAttributes(c models.Carries) error {
	if strings.TrimSpace(c.Cid) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "cid"}
	}
	if strings.TrimSpace(c.CompanyName) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "company_name"}
	}
	if strings.TrimSpace(c.Address) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "address"}
	}
	if strings.TrimSpace(c.Telephone) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "telephone"}
	}
	if strings.TrimSpace(c.LocalityId) == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "locality_id"}
	}
	return nil
}
