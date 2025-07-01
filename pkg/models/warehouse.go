package models

import "errors"

type Warehouse struct {
	Id                  int      `json:"id"`
	Warehouse_code      string   `json:"warehouse_code"`
	Address             string   `json:"address"`
	Telephone           string   `json:"telephone"`
	Minimun_capacity    int      `json:"minimum_capacity"`
	Minimun_temperature *float64 `json:"minimum_temperature"`
}

func (w Warehouse) VerifyMandatoryFieldsPresence() error {
	if w.Warehouse_code == "" {
		return errors.New("warehouse code is required")
	}
	if w.Address == "" {
		return errors.New("address is required")
	}
	if w.Telephone == "" {
		return errors.New("telephone is required")
	}
	if w.Minimun_capacity == 0 {
		return errors.New("minimun capacity is required")
	}
	if w.Minimun_temperature == nil {
		return errors.New("minimun temperature is required")
	}
	return nil
}
