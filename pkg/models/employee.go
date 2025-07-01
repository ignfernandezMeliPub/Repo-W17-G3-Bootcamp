package models

import "app/pkg/custom_errors"

type EmployeeAttributes struct {
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

type Employee struct {
	Id int `json:"id"`
	EmployeeAttributes
}

func (e *Employee) Patch(patch EmployeeRequestBody) {
	if patch.CardNumberId != nil {
		e.CardNumberId = *patch.CardNumberId
	}
	if patch.FirstName != nil {
		e.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		e.LastName = *patch.LastName
	}
	if patch.WarehouseId != nil {
		e.WarehouseId = *patch.WarehouseId
	}
}

type EmployeeRequestBody struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	WarehouseId  *int    `json:"warehouse_id"`
}

func (c EmployeeRequestBody) VerifyMandatoryFieldsPresence() error {
	if c.CardNumberId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id"}
	}

	if c.FirstName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "first_name"}
	}

	if c.LastName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "last_name"}
	}

	if c.WarehouseId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "warehouse_id"}
	}

	return nil
}
