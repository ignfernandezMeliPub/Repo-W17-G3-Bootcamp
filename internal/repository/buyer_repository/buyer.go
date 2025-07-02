package buyer_repository

import "app/pkg/models"

type BuyerRepository interface {
	FindAllBuyers() (b []models.Buyer, err error)
	FindBuyerByID(id int) (b models.Buyer, err error)
	FindBuyerByCardNumberID(cardNumberId string) (b models.Buyer, err error)
	CreateBuyer(_b models.Buyer) (b models.Buyer, err error)
	UpdateBuyer(_b models.Buyer) (b models.Buyer, err error)
	DeleteBuyerByID(id int) (e error)
}
