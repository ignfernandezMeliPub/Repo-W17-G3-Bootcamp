package models

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
	BatchNumber        *int    `json:"batch_number,omitempty"`
	CurrentQuantity    *int    `json:"current_quantity,omitempty"`
	CurrentTemperature *int    `json:"current_temperature,omitempty"`
	DueDate            *string `json:"due_date,omitempty"`
	InitialQuantity    *int    `json:"initial_quantity,omitempty"`
	ManufacturingDate  *string `json:"manufacturing_date,omitempty"`
	ManufacturingHour  *int    `json:"manufacturing_hour,omitempty"`
	MinimumTemperature *int    `json:"minimum_temperature,omitempty"`
	ProductId          *int    `json:"product_id,omitempty"`
	SectionId          *int    `json:"section_id,omitempty"`
}

type ProductBatchResponse struct {
	SectionID     int `json:"section_id" db:"section_id"`
	SectionNumber int `json:"section_number" db:"section_number"`
	ProductsCount int `json:"products_count" db:"products_count"`
}
