package models

import (
	"app/pkg/custom_errors"
	"strings"
	"time"
)

type PurchaseOrderDetail struct {
	Id              int `json:"id" db:"id"`
	ProductRecordId int `json:"product_record_id" db:"product_record_id"`
	Quantity        int `json:"quantity" db:"quantity"`
}

type PurchaseOrder struct {
	Id                   int                   `json:"id" db:"id"`
	OrderNumber          string                `json:"order_number" db:"order_number"`
	OrderDate            time.Time             `json:"order_date" db:"order_date"`
	TrackingCode         string                `json:"tracking_code" db:"tracking_code"`
	BuyerId              int                   `json:"buyer_id" db:"buyer_id"`
	PurchaseOrderDetails []PurchaseOrderDetail `json:"purchase_order_details" db:"purchase_order_details"`
}

type PurchaseOrderData struct {
	OrderNumber          *string                      `json:"order_number"`
	OrderDate            *string                      `json:"order_date"`
	TrackingCode         *string                      `json:"tracking_code"`
	BuyerId              *int                         `json:"buyer_id"`
	PurchaseOrderDetails []PurchaseOrderDetailRequest `json:"purchase_order_details"`
}

type PurchaseOrderDetailRequest struct {
	ProductRecordId *int `json:"product_record_id"`
	Quantity        *int `json:"quantity"`
}

type PurchaseOrderCreateRequest struct {
	Data *PurchaseOrderData `json:"data"`
}

func (pr PurchaseOrderCreateRequest) Verify() error {
	if pr.Data == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "data"}
	}

	p := pr.Data

	if p.OrderNumber == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "order_number"}
	}

	if strings.TrimSpace(*p.OrderNumber) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "order_number",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if p.OrderDate == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "order_date"}
	}

	_, err := time.Parse("2006-01-02", *p.OrderDate)
	if err != nil {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "order_date",
			Value:     *p.OrderDate,
			ExtraInfo: "Value must be in the format YYYY-MM-DD",
		}
	}

	if p.TrackingCode == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "tracking_code"}
	}

	if strings.TrimSpace(*p.TrackingCode) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "tracking_code",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if p.BuyerId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "buyer_id"}
	}

	if len(p.PurchaseOrderDetails) == 0 {
		return &custom_errors.MandatoryArgMissingErr{Argument: "purchase_order_details"}
	}

	for _, detail := range p.PurchaseOrderDetails {
		if detail.ProductRecordId == nil {
			return &custom_errors.MandatoryArgMissingErr{Argument: "product_record_id"}
		}
		if *detail.ProductRecordId <= 0 {
			return &custom_errors.InvalidArgValueErr{
				Argument:  "product_record_id",
				Value:     detail.ProductRecordId,
				ExtraInfo: "Value must be greater than 0",
			}
		}

		if detail.Quantity == nil {
			return &custom_errors.MandatoryArgMissingErr{Argument: "quantity"}
		}
		if *detail.Quantity <= 0 {
			return &custom_errors.InvalidArgValueErr{
				Argument:  "quantity",
				Value:     detail.Quantity,
				ExtraInfo: "Value must be greater than 0",
			}
		}
	}

	return nil
}

func (pr PurchaseOrderCreateRequest) ToPurchaseOrder() PurchaseOrder {
	p := pr.Data
	orderDate, _ := time.Parse("2006-01-02", *p.OrderDate)

	details := make([]PurchaseOrderDetail, len(p.PurchaseOrderDetails))

	for i, detail := range p.PurchaseOrderDetails {
		details[i] = PurchaseOrderDetail{
			ProductRecordId: *detail.ProductRecordId,
			Quantity:        *detail.Quantity,
		}
	}

	return PurchaseOrder{
		OrderNumber:          *p.OrderNumber,
		OrderDate:            orderDate,
		TrackingCode:         *p.TrackingCode,
		BuyerId:              *p.BuyerId,
		PurchaseOrderDetails: details,
	}
}
