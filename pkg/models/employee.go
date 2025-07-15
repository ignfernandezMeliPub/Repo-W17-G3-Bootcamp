package models

import "app/pkg/custom_errors"

type EmployeeAttributes struct {
	CardNumberId string `json:"card_number_id" db:"card_number_id"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
	WarehouseId  int    `json:"warehouse_id" db:"warehouse_id"`
}

type Employee struct {
	Id int `json:"id" db:"id"`
	EmployeeAttributes
}

type EmployeePatchRequestBody struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	WarehouseId  *int    `json:"warehouse_id"`
}

type EmployeePostRequestBody struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	WarehouseId  *int    `json:"warehouse_id"`
}

type InboundOrderEmployee struct {
	Employee
	InboundOrdersCount int `json:"inbound_orders_count" db:"inbound_orders_count"`
}

func (c EmployeePatchRequestBody) Verify() error {
	if c.CardNumberId == nil && c.FirstName == nil && c.LastName == nil && c.WarehouseId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id or first_name or last_name or warehouse_id"}
	}

	return nil
}

func (c EmployeePostRequestBody) Verify() error {
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
