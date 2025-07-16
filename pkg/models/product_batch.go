package models

import (
	"app/pkg/custom_errors"
	"reflect"
	"time"
)

type ProductBatch struct {
	ID                 int    `json:"id" db:"id"`
	BatchNumber        int    `json:"batch_number" db:"batch_number"`
	CurrentQuantity    int    `json:"current_quantity" db:"current_quantity"`
	CurrentTemperature int    `json:"current_temperature" db:"current_temperature"`
	DueDate            string `json:"due_date" db:"due_date"`
	InitialQuantity    int    `json:"initial_quantity" db:"initial_quantity"`
	ManufacturingDate  string `json:"manufacturing_date" db:"manufacturing_date"`
	ManufacturingHour  int    `json:"manufacturing_hour" db:"manufacturing_hour"`
	MinimumTemperature int    `json:"minimum_temperature" db:"minimum_temperature"`
	ProductId          int    `json:"product_id" db:"product_id"`
	SectionId          int    `json:"section_id" db:"section_id"`
}

type ProductBatchRequest struct {
	Data *ProductBatchData `json:"data"`
}

type ProductBatchData struct {
	BatchNumber        *int    `json:"batch_number"`
	CurrentQuantity    *int    `json:"current_quantity"`
	CurrentTemperature *int    `json:"current_temperature"`
	DueDate            *string `json:"due_date"`
	InitialQuantity    *int    `json:"initial_quantity"`
	ManufacturingDate  *string `json:"manufacturing_date"`
	ManufacturingHour  *int    `json:"manufacturing_hour"`
	MinimumTemperature *int    `json:"minimum_temperature"`
	ProductId          *int    `json:"product_id"`
	SectionId          *int    `json:"section_id"`
}

func (p ProductBatchRequest) Verify() error {

	mandatoryFields := map[string]any{
		"batch_number":        p.Data.BatchNumber,
		"current_quantity":    p.Data.CurrentQuantity,
		"current_temperature": p.Data.CurrentTemperature,
		"due_date":            p.Data.DueDate,
		"initial_quantity":    p.Data.InitialQuantity,
		"manufacturing_date":  p.Data.ManufacturingDate,
		"manufacturing_hour":  p.Data.ManufacturingHour,
		"minimum_temperature": p.Data.MinimumTemperature,
		"product_id":          p.Data.ProductId,
		"section_id":          p.Data.SectionId,
	}

	for field, value := range mandatoryFields {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return &custom_errors.MandatoryArgMissingErr{Argument: field}
		}
	}

	dateLayout := "2006-01-02"
	if _, err := time.Parse(dateLayout, *p.Data.DueDate); err != nil {
		return &custom_errors.InvalidArgValueErr{Argument: "due_date", Value: *p.Data.DueDate, ExtraInfo: "Invalid date format."}
	}
	if _, err := time.Parse(dateLayout, *p.Data.ManufacturingDate); err != nil {
		return &custom_errors.InvalidArgValueErr{Argument: "manufacturing_date", Value: *p.Data.ManufacturingDate, ExtraInfo: "Invalid date format."}
	}

	return nil
}

type ProductBatchResponse struct {
	SectionID     int `json:"section_id" db:"section_id"`
	SectionNumber int `json:"section_number" db:"section_number"`
	ProductsCount int `json:"products_count" db:"products_count"`
}
