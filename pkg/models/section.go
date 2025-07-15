package models

import (
	"app/pkg/custom_errors"
	"reflect"
)

type Section struct {
	ID                 int     `json:"id" db:"id"`
	SectionNumber      string  `json:"section_number" db:"section_number"`
	CurrentTemperature float32 `json:"current_temperature" db:"current_temperature"`
	MinimumTemperature float32 `json:"minimum_temperature" db:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity" db:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity" db:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity" db:"maximum_capacity"`
	WarehouseId        int     `json:"warehouse_id" db:"warehouse_id"`
	ProductTypeId      int     `json:"product_type_id" db:"product_type_id"`
}

type SectionRequest struct {
	SectionNumber      *string  `json:"section_number"`
	CurrentTemperature *float32 `json:"current_temperature"`
	MinimumTemperature *float32 `json:"minimum_temperature"`
	CurrentCapacity    *int     `json:"current_capacity"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	MaximumCapacity    *int     `json:"maximum_capacity"`
	WarehouseId        *int     `json:"warehouse_id"`
	ProductTypeId      *int     `json:"product_type_id"`
}

func (s SectionRequest) VerifyMandatoryFieldsPresence() error {
	mandatoryFields := map[string]any{
		"section_number":      s.SectionNumber,
		"current_temperature": s.CurrentTemperature,
		"minimum_temperature": s.MinimumTemperature,
		"current_capacity":    s.CurrentCapacity,
		"minimum_capacity":    s.MinimumCapacity,
		"maximum_capacity":    s.MaximumCapacity,
		"warehouse_id":        s.WarehouseId,
		"product_type_id":     s.ProductTypeId,
	}

	for field, value := range mandatoryFields {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return &custom_errors.MandatoryArgMissingErr{Argument: field}
		}
	}
	return nil
}
