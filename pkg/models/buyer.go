package models

type Buyer struct {
	Id             int    `json:"id"`
	Card_number_id int    `json:"card_number_id"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
}

type BuyerPatch struct {
	Card_number_id *int    `json:"card_number_id"`
	First_name     *string `json:"first_name"`
	Last_name      *string `json:"last_name"`
}

func (b *Buyer) Patch(id int, patch BuyerPatch) {
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
