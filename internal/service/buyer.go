package service

import (
	"app/internal/repository/buyer_repository"
	"app/internal/service/utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"strings"
)

type BuyerService interface {
	GetAllBuyers() (b []models.Buyer, err error)
	GetBuyerById(id int) (b models.Buyer, err error)
	GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error)
	CreateBuyer(_b models.Buyer) (b models.Buyer, err error)
	UpdateBuyerById(id int, _b models.BuyerPatch) (b models.Buyer, err error)
	DeleteBuyerById(id int) (err error)
}

type BuyerServiceDefault struct {
	rp buyer_repository.BuyerRepository
}

func NewBuyerDefault(rp buyer_repository.BuyerRepository) *BuyerServiceDefault {
	return &BuyerServiceDefault{rp: rp}
}

func (s *BuyerServiceDefault) GetAllBuyers() (b []models.Buyer, err error) {
	b, err = s.rp.GetAllBuyers()
	return
}
func (s *BuyerServiceDefault) GetBuyerById(id int) (b models.Buyer, err error) {
	b, err = s.rp.GetBuyerById(id)
	return
}
func (s *BuyerServiceDefault) GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error) {
	b, err = s.rp.GetBuyerByCardNumberId(cardNumberId)
	return
}
func (s *BuyerServiceDefault) CreateBuyer(_b models.Buyer) (b models.Buyer, err error) {

	if err = validateBuyerAttributes(_b); err != nil {
		return
	}

	_, err = s.rp.GetBuyerByCardNumberId(_b.CardNumberId)

	if err = utils.ExpectError(
		err,
		custom_errors.ErrNotFound,
		&custom_errors.UniqueAttributeViolationErr{
			AttributeName: "card_number_id",
			Value:         _b.CardNumberId,
		}); err != nil {
		return
	}

	b, err = s.rp.CreateBuyer(_b)
	return
}

func (s *BuyerServiceDefault) UpdateBuyerById(id int, _b models.BuyerPatch) (b models.Buyer, err error) {

	buyer, err := s.rp.GetBuyerById(id)

	if err != nil {
		return
	}

	buyer.Patch(_b)

	if err = validateBuyerAttributes(buyer); err != nil {
		return
	}

	if _b.CardNumberId != nil {

		oldBuyer, e := s.rp.GetBuyerByCardNumberId(*_b.CardNumberId)
		err = e

		if err = utils.ExpectErrorOrNilCondition(
			err,
			oldBuyer.Id != id,
			custom_errors.ErrNotFound,
			&custom_errors.UniqueAttributeViolationErr{
				AttributeName: "card_number_id",
				Value:         *_b.CardNumberId,
			}); err != nil {
			return
		}
	}

	b, err = s.rp.UpdateBuyer(buyer)
	return
}

func (s *BuyerServiceDefault) DeleteBuyerById(id int) (err error) {

	_, err = s.rp.GetBuyerById(id)

	if err != nil {
		return
	}

	err = s.rp.DeleteBuyerById(id)
	return
}

func validateBuyerAttributes(buyer models.Buyer) error {

	if strings.TrimSpace(buyer.CardNumberId) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "card_number_id",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if strings.TrimSpace(buyer.FirstName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "first_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}

	if strings.TrimSpace(buyer.LastName) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "last_name",
			Value:     "",
			ExtraInfo: "Value must be non-empty",
		}
	}
	return nil
}
