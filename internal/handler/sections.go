package handler

import (
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type SectionsController struct {
	sv service.SectionsService
}

func NewSectionsController(sv service.SectionsService) *SectionsController {
	return &SectionsController{sv: sv}
}

func (c *SectionsController) GetSections(w http.ResponseWriter, r *http.Request) {
	res, err := c.sv.GetSections()
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return
}

func (c *SectionsController) GetSection(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		valueErr := custom_errors.InvalidArgValueErr{Value: idString, Argument: "id"}
		response.Error(w, http.StatusBadRequest, valueErr.Error())
		return
	}
	sec, err := c.sv.GetSectionByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": sec})
	return

}

func (c *SectionsController) CreateSection(w http.ResponseWriter, r *http.Request) {
	var section models.Section
	if err := request.JSON(r, &section); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res, err := c.sv.CreateSection(section)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"data": res})
	return

}

func (c *SectionsController) UpdateSection(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		valueErr := custom_errors.InvalidArgValueErr{Value: idString, Argument: "id"}
		response.Error(w, http.StatusBadRequest, valueErr.Error())
		return
	}
	var section models.SectionPatch
	if err := request.JSON(r, &section); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res, err := c.sv.UpdateSection(id, section)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return

}

func (c *SectionsController) DeleteSection(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		valueErr := custom_errors.InvalidArgValueErr{Value: idString, Argument: "id"}
		response.Error(w, http.StatusBadRequest, valueErr.Error())
		return
	}
	err = c.sv.DeleteSection(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
