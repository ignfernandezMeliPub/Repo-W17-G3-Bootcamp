package models

type Section struct {
	ID                 int     `json:"id"`
	SectionNumber      string  `json:"section_number"`
	CurrentTemperature float32 `json:"current_temperature"`
	MinimumTemperature float32 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseId        int     `json:"warehouse_id"`
	ProductTypeId      int     `json:"product_type_id"`
}

type SectionPatch struct {
	SectionNumber      *string  `json:"section_number,omitempty"`
	CurrentTemperature *float32 `json:"current_temperature,omitempty"`
	MinimumTemperature *float32 `json:"minimum_temperature,omitempty"`
	CurrentCapacity    *int     `json:"current_capacity,omitempty"`
	MinimumCapacity    *int     `json:"minimum_capacity,omitempty"`
	MaximumCapacity    *int     `json:"maximum_capacity,omitempty"`
	WarehouseId        *int     `json:"warehouse_id,omitempty"`
	ProductTypeId      *int     `json:"product_type_id,omitempty"`
}
