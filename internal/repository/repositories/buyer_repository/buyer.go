package buyer_repository

import "app/pkg/models"

type BuyerRepository interface {
	GetAllBuyers() (b []models.Buyer, err error)
	GetBuyerById(id int) (b models.Buyer, err error)
	GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error)
	CreateBuyer(_b models.Buyer) (b models.Buyer, err error)
	UpdateBuyer(_b models.Buyer) (b models.Buyer, err error)
	DeleteBuyerById(id int) (e error)
	GetBuyerPurchaseOrdersCount(buyerId int) (p []models.BuyerPurchaseOrdersCount, err error)
	GetBuyersPurchaseOrdersCount() (p []models.BuyerPurchaseOrdersCount, err error)
}
