package models

import (
	"app/pkg/custom_errors"
)

type Buyer struct {
	Id             int    `json:"id"`
	Card_number_id string `json:"card_number_id"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
}

type BuyerPatch struct {
	Card_number_id *string `json:"card_number_id"`
	First_name     *string `json:"first_name"`
	Last_name      *string `json:"last_name"`
}

func (b BuyerPatch) VerifyMandatoryFieldsPresence() error {
	return nil
}

func (b *Buyer) Patch(patch BuyerPatch) {
	if patch.Card_number_id != nil {
		b.Card_number_id = *patch.Card_number_id
	}
	if patch.First_name != nil {
		b.First_name = *patch.First_name
	}
	if patch.Last_name != nil {
		b.Last_name = *patch.Last_name
	}
	return
}

type BuyerCreateRequest struct {
	Card_number_id *string `json:"card_number_id"`
	First_name     *string `json:"first_name"`
	Last_name      *string `json:"last_name"`
}

func (b BuyerCreateRequest) VerifyMandatoryFieldsPresence() error {
	if b.Card_number_id == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "Card_number_id"}
	}

	if b.First_name == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "First_name"}
	}

	if b.Last_name == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "Last_name"}
	}

	return nil
}

func (b *BuyerCreateRequest) ToBuyer() Buyer {
	return Buyer{
		Card_number_id: *b.Card_number_id,
		First_name:     *b.First_name,
		Last_name:      *b.Last_name,
	}
}
