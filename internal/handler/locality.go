package handler

import (
	"app/internal/handler/dto"
	"app/internal/handler/utils"
	"app/internal/service"
	"app/pkg/models"
	"github.com/bootcamp-go/web/response"
	"net/http"
)

type LocalityHandler struct {
	service service.LocalityService
}

func NewLocalityHandler(service service.LocalityService) LocalityHandler {
	return LocalityHandler{service: service}
}

// CreateLocality Creates a new locality
func (h *LocalityHandler) CreateLocality(w http.ResponseWriter, r *http.Request) {
	var createLocalityDto dto.CreateLocalityDto
	createLocalityDto, err := utils.InstantiateVarFromBody(&r.Body, createLocalityDto)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	newLocality, err := h.service.CreateLocality(*createLocalityDto.Data.Id, *createLocalityDto.Data.LocalityName, *createLocalityDto.Data.ProvinceName, *createLocalityDto.Data.CountryName)
	if err != nil {
		utils.ResponseHttpError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": []models.Locality{newLocality},
	})
}

// GetLocalitySellerCount retrieves seller count statistics for localities.
func (h *LocalityHandler) GetLocalitySellerCount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// ? LocalitySellerCount de una locality en particular
	if id != "" {
		result, err := h.service.GetLocalitySellerCount(id)
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": []models.LocalitySellerCount{result},
		})
		// ? LocalitySellerCount de cada una de las localities
	} else {
		result, err := h.service.GetLocalitiesSellerCount()
		if err != nil {
			utils.ResponseHttpError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": result,
		})
	}
}
