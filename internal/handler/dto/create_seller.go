package dto

import "app/pkg/custom_errors"

type CreateSellerDto struct {
	CompanyId   *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
}

func (c CreateSellerDto) VerifyMandatoryFieldsPresence() error {
	if c.CompanyId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "cid"}
	}

	if c.CompanyName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "company_name"}
	}

	if c.Address == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "address"}
	}

	if c.Telephone == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "telephone"}
	}

	return nil
}
