package models

import (
	"app/pkg/custom_errors"
)

type Carries struct {
	Id          int    `json:"id" db:"id"`
	Cid         string `json:"cid" db:"cid"`
	CompanyName string `json:"company_name" db:"company_name"`
	Address     string `json:"address" db:"address"`
	Telephone   string `json:"telephone" db:"telephone"`
	LocalityId  string `json:"locality_id" db:"locality_id"`
}

type CarriesReport struct {
	LocalityId   string `json:"locality_id" db:"locality_id"`
	LocalityName string `json:"locality_name" db:"locality_name"`
	CarriesCount int    `json:"carries_count" db:"carries_count"`
}

func (c Carries) VerifyMandatoryFieldsPresence() error {
	if c.Cid == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "cid"}
	}
	if c.CompanyName == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "company_name"}
	}
	if c.Address == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "address"}
	}
	if c.Telephone == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "telephone"}
	}
	if c.LocalityId == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "locality_id"}
	}
	return nil
}
