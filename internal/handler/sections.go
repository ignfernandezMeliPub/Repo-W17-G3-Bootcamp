package handler

import (
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
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

func (c *SectionsController) GetSections(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetSections", logger.LogStatusInProgress)

	res, err := c.sv.GetAllSections()
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetSections", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetSections", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return nil
}

func (c *SectionsController) GetSectionById(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetSectionById", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetSectionById", err, httpStatus)
		return nil
	}
	sec, err := c.sv.GetSectionById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetSectionById", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetSectionById", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{"data": sec})
	return nil
}

func (c *SectionsController) CreateSection(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateSection", logger.LogStatusInProgress)

	var section models.SectionRequest
	section, err = utils.InstantiateVarFromBody(&r.Body, section)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateSection", err, httpStatus)
		return nil
	}
	res, err := c.sv.CreateSection(section)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateSection", err, httpStatus)
		return nil
	}

	utils.Log(r, "CreateSection", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{"data": res})
	return nil
}

func (c *SectionsController) PatchSection(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "PatchSection", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSection", err, httpStatus)
		return nil
	}
	var section models.SectionRequest
	if err = request.JSON(r, &section); err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSection", err, httpStatus)
		return nil
	}
	res, err := c.sv.UpdateSectionById(id, section)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "PatchSection", err, httpStatus)
		return nil
	}

	utils.Log(r, "PatchSection", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{"data": res})
	return nil
}

func (c *SectionsController) DeleteSection(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "DeleteSection", logger.LogStatusInProgress)

	id, err := utils.GetURLParamAs(r, "id", strconv.Atoi)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteSection", err, httpStatus)
		return nil
	}
	err = c.sv.DeleteSectionById(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "DeleteSection", err, httpStatus)
		return nil
	}

	utils.Log(r, "DeleteSection", logger.LogStatusSuccess)

	response.JSON(w, http.StatusNoContent, nil)
	return nil
}

func (c *SectionsController) GetAllProductBatchesBySection(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetAllProductBatchesBySection", logger.LogStatusInProgress)

	id, err := utils.GetQueryParamAs(r, "id", strconv.Atoi)

	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllProductBatchesBySection", err, httpStatus)
		return nil
	}

	var data []models.ProductBatchResponse
	data, err = c.sv.GetProductBatchBySection(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetAllProductBatchesBySection", err, httpStatus)
		return nil
	}

	utils.Log(r, "GetAllProductBatchesBySection", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": data,
	})
	return nil
}
