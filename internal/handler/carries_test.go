package handler

import (
	"app/pkg/models"
	"app/test/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/stretchr/testify/require"
)

func Test_CreateCarrie(t *testing.T) {

	t.Run("should create a new carry successfully", func(t *testing.T) {
		carrie := `{
			"cid": "1",
			"company_name": "Company 1",
			"address": "Street 123",
			"telephone": "1234567890",
			"locality_id": "L001"
		}`
		inputCarrie := models.Carries{
			Cid:         "1",
			CompanyName: "Company 1",
			Address:     "Street 123",
			Telephone:   "1234567890",
			LocalityId:  "L001",
		}
		returned := inputCarrie
		returned.Id = 1

		mockService := new(service.MockCarriesService)
		mockService.On("CreateCarrie", inputCarrie).Return(returned, nil)

		hd := NewCarriesHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/carries", bytes.NewBufferString(carrie))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateCarrie(response, request)

		require.Equal(t, http.StatusCreated, response.Code)

		var res map[string]models.Carries
		err := json.Unmarshal(response.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, returned, res["data"])

		mockService.AssertExpectations(t)
	})
	t.Run("should return error if the request body is invalid", func(t *testing.T) {
		carrie := `{
			"cid": "1",
			"company_name": "Company 1",
			"address": "Street 123",
			"telephone": "1234567890",
			"locality_id": "L001",
		}`
		mockService := new(service.MockCarriesService)
		hd := NewCarriesHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/carries", bytes.NewBufferString(carrie))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateCarrie(response, request)

		require.Equal(t, http.StatusBadRequest, response.Code)
		mockService.AssertNotCalled(t, "CreateCarrie")
	})

	t.Run("should return error for validateCarriesAttributes", func(t *testing.T) {
		carrie := `{
			"cid": " ",
			"company_name": "Company 1",
			"address": "Street 123",
			"telephone": "1234567890",
			"locality_id": "L001"
		}`
		mockService := new(service.MockCarriesService)
		hd := NewCarriesHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/carries", bytes.NewBufferString(carrie))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateCarrie(response, request)

		require.Equal(t, http.StatusUnprocessableEntity, response.Code)
		require.Contains(t, response.Body.String(), "cid")

		mockService.AssertNotCalled(t, "CreateCarrie")
	})

	t.Run("should handle error from service", func(t *testing.T) {
		carrie := `{
			"cid": "1",
			"company_name": "Company 1",
			"address": "Street 123",
			"telephone": "1234567890",
			"locality_id": "L001"
		}`
		inputCarrie := models.Carries{
			Cid:         "1",
			CompanyName: "Company 1",
			Address:     "Street 123",
			Telephone:   "1234567890",
			LocalityId:  "L001",
		}
		mockService := new(service.MockCarriesService)
		mockService.On("CreateCarrie", inputCarrie).Return(models.Carries{}, errors.New("database connection failed"))
		hd := NewCarriesHandler(mockService)

		request := httptest.NewRequest(http.MethodPost, "/carries", bytes.NewBufferString(carrie))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		hd.CreateCarrie(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
		mockService.AssertExpectations(t)
	})
}
func Test_validateCarriesAttributes(t *testing.T) {

	carrie := models.Carries{
		Cid:         "1",
		CompanyName: "Company 1",
		Address:     "Street 123",
		Telephone:   "1234567890",
		LocalityId:  "L001",
	}
	tests := []struct {
		name          string
		fieldToModify func(c models.Carries) models.Carries
		errorExpected string
	}{
		{
			name: "should return error if the cid is empty",
			fieldToModify: func(c models.Carries) models.Carries {
				c.Cid = ""
				return c
			},
			errorExpected: "cid",
		},
		{
			name: "should return error if the company name is empty",
			fieldToModify: func(c models.Carries) models.Carries {
				c.CompanyName = ""
				return c
			},
			errorExpected: "company_name",
		},
		{
			name: "should return error if the address is empty",
			fieldToModify: func(c models.Carries) models.Carries {
				c.Address = ""
				return c
			},
			errorExpected: "address",
		},
		{
			name: "should return error if the telephone is empty",
			fieldToModify: func(c models.Carries) models.Carries {
				c.Telephone = ""
				return c
			},
			errorExpected: "telephone",
		},
		{
			name: "should return error if the locality id is empty",
			fieldToModify: func(c models.Carries) models.Carries {
				c.LocalityId = ""
				return c
			},
			errorExpected: "locality_id",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputCarrie := test.fieldToModify(carrie)
			err := validateCarriesAttributes(inputCarrie)
			require.Error(t, err)
			require.Contains(t, err.Error(), test.errorExpected)
		})
	}
}
