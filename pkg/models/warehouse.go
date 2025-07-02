package models

import "errors"

type Warehouse struct {
	Id                 int      `json:"id"`
	WarehouseCode      string   `json:"warehouse_code"`
	Address            string   `json:"address"`
	Telephone          string   `json:"telephone"`
	MinimumCapacity    int      `json:"minimum_capacity"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
}

func (w Warehouse) VerifyMandatoryFieldsPresence() error {
	if w.WarehouseCode == "" {
		return errors.New("warehouse code is required")
	}
	if w.Address == "" {
		return errors.New("address is required")
	}
	if w.Telephone == "" {
		return errors.New("telephone is required")
	}
	if w.MinimumCapacity == 0 {
		return errors.New("minimun capacity is required")
	}
	if w.MinimumTemperature == nil {
		return errors.New("minimun temperature is required")
	}
	return nil
}
