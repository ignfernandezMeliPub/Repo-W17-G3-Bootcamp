package seller_repository

import (
	"app/pkg/models"
)

type SellerRepository interface {
	GetSellerById(id int) (models.Seller, error)
	CompanyIdIsUsed(companyId int) (bool, error)
	CreateSeller(seller models.Seller) (models.Seller, error)
	GetAllSellers() ([]models.Seller, error)
	DeleteSellerById(id int) error
	UpdateSellerById(seller models.Seller) (models.Seller, error)
}
