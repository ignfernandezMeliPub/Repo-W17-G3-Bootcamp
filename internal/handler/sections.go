package handler

import (
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
)

type SectionsController struct {
	sv service.SectionsService
}

func NewSectionsController(sv service.SectionsService) *SectionsController {
	return &SectionsController{sv: sv}
}

func (c *SectionsController) GetSections(w http.ResponseWriter, _ *http.Request) {
	res, err := c.sv.GetAllSections()
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return
}

func (c *SectionsController) GetSectionById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	sec, err := c.sv.GetSectionById(id)
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
		utils.ResponseHttpError(w, &custom_errors.InvalidBodyError{})
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

func (c *SectionsController) PatchSection(w http.ResponseWriter, r *http.Request) {
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
	res, err := c.sv.UpdateSectionById(id, section)
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
	err = c.sv.DeleteSectionById(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func (c *SectionsController) GetAllProductBatchesBySection(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		res, err := c.sv.GetAllProductBatchesBySection()
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{"data": res})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ResponseHttpError(w, &custom_errors.UrlParamDecodeError{UrlParam: "id"})
		return
	}
	res, err := c.sv.GetProductBatchBySectionId(id)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return
}
