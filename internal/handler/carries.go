package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"net/http"
	"strings"

	"github.com/bootcamp-go/web/response"
)

type CarriesHandler struct {
	sv service.ICarriesService
}

func NewCarriesHandler(sv service.ICarriesService) *CarriesHandler {
	return &CarriesHandler{sv: sv}
}

func (h *CarriesHandler) CreateCarrie(w http.ResponseWriter, r *http.Request) {
	var c models.Carries
	c, err := utils.InstantiateVarFromBody(&r.Body, c)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	err = validateCarriesAttributes(c)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	data, err := h.sv.CreateCarrie(c)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": data,
	})
}

func (h *CarriesHandler) GetCarriesReport(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var idPtr *string
	if id != "" {
		idPtr = &id
	}
	report, err := h.sv.GetCarriesReport(idPtr)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": report,
	})
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
