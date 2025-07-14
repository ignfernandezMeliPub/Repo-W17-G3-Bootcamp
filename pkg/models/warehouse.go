package models

import (
	"app/pkg/custom_errors"
)

type Warehouse struct {
	Id                 int      `json:"id" db:"id"`
	WarehouseCode      string   `json:"warehouse_code" db:"warehouse_code"`
	Address            string   `json:"address" db:"address"`
	Telephone          string   `json:"telephone" db:"telephone"`
	MinimumCapacity    int      `json:"minimum_capacity" db:"minimum_capacity"`
	MinimumTemperature *float64 `json:"minimum_temperature" db:"minimum_temperature"`
}

func (w Warehouse) VerifyMandatoryFieldsPresence() error {
	if w.WarehouseCode == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "warehouse_code"}
	}
	if w.Address == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "address"}
	}
	if w.Telephone == "" {
		return &custom_errors.MandatoryArgMissingErr{Argument: "telephone"}
	}
	if w.MinimumCapacity == 0 {
		return &custom_errors.MandatoryArgMissingErr{Argument: "minimum_capacity"}
	}
	if w.MinimumTemperature == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "minimum_temperature"}
	}
	return nil
}
