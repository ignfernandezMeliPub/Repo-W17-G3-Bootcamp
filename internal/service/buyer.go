package service

import (
	"app/internal/repository/buyer_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"errors"
)

type BuyerService interface {
	FindAllBuyers() (b []models.Buyer, err error)
	FindBuyerByID(id int) (b models.Buyer, err error)
	FindBuyerByCardNumberID(card_number_id int) (b models.Buyer, err error)
	CreateBuyer(_b models.Buyer) (b models.Buyer, err error)
	UpdateBuyerByID(id int, _b models.BuyerPatch) (b models.Buyer, err error)
	DeleteBuyerByID(id int) (err error)
}

type BuyerServiceDefault struct {
	rp buyer_repository.BuyerRepository
}

func NewBuyerDefault(rp buyer_repository.BuyerRepository) *BuyerServiceDefault {
	return &BuyerServiceDefault{rp: rp}
}

func (s *BuyerServiceDefault) FindAllBuyers() (b []models.Buyer, err error) {
	b, err = s.rp.FindAllBuyers()
	return
}
func (s *BuyerServiceDefault) FindBuyerByID(id int) (b models.Buyer, err error) {
	b, err = s.rp.FindBuyerByID(id)
	return
}
func (s *BuyerServiceDefault) FindBuyerByCardNumberID(card_number_id int) (b models.Buyer, err error) {
	b, err = s.rp.FindBuyerByCardNumberID(card_number_id)
	return
}
func (s *BuyerServiceDefault) CreateBuyer(_b models.Buyer) (b models.Buyer, err error) {

	_, err = s.rp.FindBuyerByCardNumberID(_b.Card_number_id)

	if err == nil {
		err = &custom_errors.UniqueAttributeViolationErr{AttributeName: "Card_number_id", Value: _b.Card_number_id}
		return
	} else if !errors.As(err, &custom_errors.ErrNotFound) {
		return
	}

	b, err = s.rp.CreateBuyer(_b)
	return
}
func (s *BuyerServiceDefault) UpdateBuyerByID(id int, _b models.BuyerPatch) (b models.Buyer, err error) {

	buyer, err := s.rp.FindBuyerByID(id)

	if err != nil {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}

	if _b.Card_number_id != nil {

		old_buyer, e := s.rp.FindBuyerByCardNumberID(*_b.Card_number_id)
		err = e

		if err == nil && old_buyer.Id != id { // Matches new card_number_id with distinct buyer
			err = &custom_errors.UniqueAttributeViolationErr{AttributeName: "Card_number_id", Value: _b.Card_number_id}
			return
		} else if !errors.As(err, &custom_errors.ErrNotFound) {
			return
		}
	}

	buyer.Patch(id, _b)
	b, err = s.rp.UpdateBuyer(buyer)
	return
}

func (s *BuyerServiceDefault) DeleteBuyerByID(id int) (err error) {

	_, err = s.rp.FindBuyerByID(id)

	if err != nil {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}
	err = s.rp.DeleteBuyerByID(id)
	return
}
