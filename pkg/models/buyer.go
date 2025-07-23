package models

import (
	"app/pkg/custom_errors"
	"strings"
)

type Buyer struct {
	Id           int    `json:"id" db:"id"`
	CardNumberId string `json:"card_number_id" db:"card_number_id"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
}

type BuyerPatch struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
}

func (b BuyerPatch) Verify() error {
	if b.CardNumberId == nil && b.FirstName == nil && b.LastName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id or first_name or last_name"}
	}

	if b.CardNumberId != nil && strings.TrimSpace(*b.CardNumberId) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "card_number_id",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if b.FirstName != nil && strings.TrimSpace(*b.FirstName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "first_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if b.LastName != nil && strings.TrimSpace(*b.LastName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "last_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}
	return nil
}

func (b *Buyer) Patch(patch BuyerPatch) {
	if patch.CardNumberId != nil {
		b.CardNumberId = *patch.CardNumberId
	}
	if patch.FirstName != nil {
		b.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		b.LastName = *patch.LastName
	}
}

type BuyerCreateRequest struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
}

func (b BuyerCreateRequest) Verify() error {
	if b.CardNumberId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id"}
	}

	if b.FirstName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "first_name"}
	}

	if b.LastName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "last_name"}
	}

	if strings.TrimSpace(*b.CardNumberId) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "card_number_id",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if strings.TrimSpace(*b.FirstName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "first_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if strings.TrimSpace(*b.LastName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "last_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	return nil
}

func (b BuyerCreateRequest) ToBuyer() Buyer {
	return Buyer{
		CardNumberId: *b.CardNumberId,
		FirstName:    *b.FirstName,
		LastName:     *b.LastName,
	}
}

type BuyerPurchaseOrdersCount struct {
	Id                  int    `json:"id" db:"id"`
	CardNumberId        string `json:"card_number_id" db:"card_number_id"`
	FirstName           string `json:"first_name" db:"first_name"`
	LastName            string `json:"last_name" db:"last_name"`
	PurchaseOrdersCount int    `json:"purchase_orders_count" db:"purchase_orders_count"`
}
