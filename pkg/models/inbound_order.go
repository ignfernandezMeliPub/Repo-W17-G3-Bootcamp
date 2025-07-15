package models

import "app/pkg/custom_errors"

type InboundOrderDetails struct {
	OrderDate      string `json:"order_date" db:"order_date"`
	OrderNumber    string `json:"order_number" db:"order_number"`
	EmployeeId     int    `json:"employee_id" db:"employee_id"`
	ProductBatchId int    `json:"product_batch_id" db:"product_batch_id"`
	WarehouseId    int    `json:"warehouse_id" db:"warehouse_id"`
}

type InboundOrder struct {
	Id int `json:"id" db:"id"`
	InboundOrderDetails
}

type InboundOrderRequestBody struct {
	Data *InboundOrderData `json:"data"`
}

type InboundOrderData struct {
	OrderDate      *string `json:"order_date"`
	OrderNumber    *string `json:"order_number"`
	EmployeeId     *int    `json:"employee_id"`
	ProductBatchId *int    `json:"product_batch_id"`
	WarehouseId    *int    `json:"warehouse_id"`
}

func (c InboundOrderRequestBody) VerifyMandatoryFieldsPresence() error {
	if c.Data == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data"}
	}
	if c.Data.OrderDate == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.order_date"}
	}

	if c.Data.OrderNumber == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.order_number"}
	}

	if c.Data.EmployeeId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.employee_id"}
	}

	if c.Data.ProductBatchId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.product_batch_id"}
	}

	if c.Data.WarehouseId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data.warehouse_id"}
	}

	return nil
}
