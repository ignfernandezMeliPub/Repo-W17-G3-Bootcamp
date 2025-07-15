package buyer_repository

import "app/pkg/models"

type BuyerRepository interface {
	GetAllBuyers() (b []models.Buyer, err error)
	GetBuyerById(id int) (b models.Buyer, err error)
	GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error)
	CreateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error)
	UpdateBuyer(buyer models.Buyer) (updatedBuyer models.Buyer, err error)
	DeleteBuyerById(id int) (err error)
	GetBuyerPurchaseOrdersCount(buyerId int) (b []models.BuyerPurchaseOrdersCount, err error)
	GetBuyersPurchaseOrdersCount() (b []models.BuyerPurchaseOrdersCount, err error)
}
