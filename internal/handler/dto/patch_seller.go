package dto

import "app/pkg/custom_errors"

type PatchSellerDto struct {
	Id          *int    `json:"id"`
	CompanyId   *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
}

func (c PatchSellerDto) VerifyMandatoryFieldsPresence() error {
	if c.Id == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "id"}
	}

	if c.CompanyId == nil && c.CompanyName == nil && c.Address == nil && c.Telephone == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "cid or company_name or address or telephone"}
	}

	return nil
}
