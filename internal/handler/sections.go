package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
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
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return
}

func (c *SectionsController) GetSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	sec, err := c.sv.GetSectionByID(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": sec})
	return

}

func (c *SectionsController) CreateSection(w http.ResponseWriter, r *http.Request) {
	var section models.SectionRequest
	if err := request.JSON(r, &section); err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	res, err := c.sv.CreateSection(section)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"data": res})
	return

}

func (c *SectionsController) UpdateSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	var section models.SectionRequest
	if err := request.JSON(r, &section); err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	res, err := c.sv.UpdateSection(id, section)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return

}

func (c *SectionsController) DeleteSection(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	err = c.sv.DeleteSection(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}
