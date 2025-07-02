package models

import (
	"app/pkg/custom_errors"
)

type Buyer struct {
	Id           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerPatch struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
}

func (b BuyerPatch) VerifyMandatoryFieldsPresence() error {
	if b.CardNumberId == nil && b.FirstName == nil && b.LastName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id or first_name or last_name"}
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
	return
}

type BuyerCreateRequest struct {
	CardNumberId *string `json:"card_number_id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
}

func (b BuyerCreateRequest) VerifyMandatoryFieldsPresence() error {
	if b.CardNumberId == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id"}
	}

	if b.FirstName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "first_name"}
	}

	if b.LastName == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "last_name"}
	}

	return nil
}

func (b *BuyerCreateRequest) ToBuyer() Buyer {
	return Buyer{
		CardNumberId: *b.CardNumberId,
		FirstName:    *b.FirstName,
		LastName:     *b.LastName,
	}
}
