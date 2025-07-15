package service

import (
	"app/internal/repository/repositories/buyer_repository"
	"app/pkg/models"
)

type BuyerService interface {
	GetAllBuyers() (b []models.Buyer, err error)
	GetBuyerById(id int) (b models.Buyer, err error)
	GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error)
	CreateBuyer(_b models.Buyer) (b models.Buyer, err error)
	UpdateBuyerById(id int, _b models.BuyerPatch) (b models.Buyer, err error)
	DeleteBuyerById(id int) (err error)
	GetBuyerPurchaseOrdersCount(buyerId int) (p []models.BuyerPurchaseOrdersCount, err error)
	GetBuyersPurchaseOrdersCount() (p []models.BuyerPurchaseOrdersCount, err error)
}

type BuyerServiceDefault struct {
	rp buyer_repository.BuyerRepository
}

func NewBuyerDefault(rp buyer_repository.BuyerRepository) *BuyerServiceDefault {
	return &BuyerServiceDefault{rp: rp}
}

func (s *BuyerServiceDefault) GetAllBuyers() (b []models.Buyer, err error) {
	return s.rp.GetAllBuyers()
}
func (s *BuyerServiceDefault) GetBuyerById(id int) (b models.Buyer, err error) {
	return s.rp.GetBuyerById(id)
}
func (s *BuyerServiceDefault) GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error) {
	return s.rp.GetBuyerByCardNumberId(cardNumberId)
}
func (s *BuyerServiceDefault) CreateBuyer(_b models.Buyer) (b models.Buyer, err error) {
	return s.rp.CreateBuyer(_b)
}

func (s *BuyerServiceDefault) UpdateBuyerById(id int, _b models.BuyerPatch) (b models.Buyer, err error) {

	buyer, err := s.rp.GetBuyerById(id)

	if err != nil {
		return
	}

	buyer.Patch(_b)

	b, err = s.rp.UpdateBuyer(buyer)
	return
}

func (s *BuyerServiceDefault) DeleteBuyerById(id int) (err error) {
	return s.rp.DeleteBuyerById(id)
}

func (s *BuyerServiceDefault) GetBuyerPurchaseOrdersCount(buyerId int) (p []models.BuyerPurchaseOrdersCount, err error) {
	return s.rp.GetBuyerPurchaseOrdersCount(buyerId)
}

func (s *BuyerServiceDefault) GetBuyersPurchaseOrdersCount() (p []models.BuyerPurchaseOrdersCount, err error) {
	return s.rp.GetBuyersPurchaseOrdersCount()
}
