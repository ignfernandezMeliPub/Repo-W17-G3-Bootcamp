package handler

import (
	"app/internal/handler/dto"
	"app/internal/handler/utils"
	"app/internal/logger"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

type LocalityHandler struct {
	service service.LocalityService
}

func NewLocalityHandler(service service.LocalityService) LocalityHandler {
	return LocalityHandler{service: service}
}

// CreateLocality Creates a new locality
func (h *LocalityHandler) CreateLocality(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "CreateLocality", logger.LogStatusInProgress)

	var createLocalityDto dto.CreateLocalityDto
	createLocalityDto, err = utils.InstantiateVarFromBody(&r.Body, createLocalityDto)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateLocality", err, httpStatus)
		return
	}

	newLocality, err := h.service.CreateLocality(*createLocalityDto.Data.Id, *createLocalityDto.Data.LocalityName, *createLocalityDto.Data.ProvinceName, *createLocalityDto.Data.CountryName)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "CreateLocality", err, httpStatus)
		return
	}

	utils.Log(r, "CreateLocality", logger.LogStatusSuccess)

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": []models.Locality{newLocality},
	})
	return
}

func (h *LocalityHandler) GetCarriesReport(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetCarriesReport", logger.LogStatusInProgress)

	id := r.URL.Query().Get("id")

	report, err := h.service.GetCarriesReport(id)
	if err != nil {
		httpStatus := utils.ResponseHttpError(w, err)
		utils.LogError(r, "GetCarriesReport", err, httpStatus)
		return
	}

	utils.Log(r, "GetCarriesReport", logger.LogStatusSuccess)

	response.JSON(w, http.StatusOK, map[string]any{
		"data": report,
	})
	return
}

// GetLocalitySellerCount retrieves seller count statistics for localities.
func (h *LocalityHandler) GetLocalitySellerCount(w http.ResponseWriter, r *http.Request) (err error) {
	utils.Log(r, "GetLocalitySellerCount", logger.LogStatusInProgress)

	id := r.URL.Query().Get("id")

	// ? LocalitySellerCount de una locality en particular
	if id != "" {
		var result models.LocalitySellerCount
		result, err = h.service.GetLocalitySellerCount(id)
		if err != nil {
			httpStatus := utils.ResponseHttpError(w, err)
			utils.LogError(r, "GetLocalitySellerCount", err, httpStatus)
			return
		}

		utils.Log(r, "GetLocalitySellerCount", logger.LogStatusSuccess)

		response.JSON(w, http.StatusOK, map[string]any{
			"data": []models.LocalitySellerCount{result},
		})
		// ? LocalitySellerCount de cada una de las localities
	} else {
		var result []models.LocalitySellerCount
		result, err = h.service.GetLocalitiesSellerCount()
		if err != nil {
			httpStatus := utils.ResponseHttpError(w, err)
			utils.LogError(r, "GetLocalitySellerCount", err, httpStatus)
			return
		}

		utils.Log(r, "GetLocalitySellerCount", logger.LogStatusSuccess)

		response.JSON(w, http.StatusOK, map[string]any{
			"data": result,
		})
	}
	return
}
