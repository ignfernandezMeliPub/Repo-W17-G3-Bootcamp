package models

import (
	"app/pkg/custom_errors"
)

type ProductRecord struct {
	ID             int     `json:"id" db:"id"`
	LastUpdateDate string  `json:"last_update_date" db:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price" db:"purchase_price"`
	SalePrice      float64 `json:"sale_price" db:"sale_price"`
	ProductID      int     `json:"product_id" db:"product_id"`
}

type ProductRecordRequest struct {
	LastUpdateDate *string  `json:"last_update_date"`
	PurchasePrice  *float64 `json:"purchase_price"`
	SalePrice      *float64 `json:"sale_price"`
	ProductID      *int     `json:"product_id"`
}

func (p ProductRecordRequest) VerifyMandatoryFieldsPresence() error {
	if p.LastUpdateDate == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "last_update_date"}
	}

	if p.PurchasePrice == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "purchase_price"}
	}

	if p.SalePrice == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "sale_price"}
	}

	if p.ProductID == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "product_id"}
	}
	return nil
}
