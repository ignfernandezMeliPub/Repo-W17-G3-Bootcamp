package dto

import "app/pkg/custom_errors"

type CreateLocalityDto struct {
	Data *CreateLocalityData `json:"data"`
}

func (c CreateLocalityDto) VerifyMandatoryFieldsPresence() error {
	if c.Data == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data"}
	}

	if c.Data.Id == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.id"}
	}

	if c.Data.LocalityName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.locality_name"}
	}

	if c.Data.ProvinceName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.province_name"}
	}

	if c.Data.CountryName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.country_name"}
	}

	return nil
}

type CreateLocalityData struct {
	Id           *string `json:"id"`
	LocalityName *string `json:"locality_name"`
	ProvinceName *string `json:"province_name"`
	CountryName  *string `json:"country_name"`
}
