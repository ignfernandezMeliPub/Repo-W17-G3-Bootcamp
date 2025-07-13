package seller_repository

import (
	"app/pkg/models"
)

type SellerRepository interface {
	GetSellerById(id int) (models.Seller, error)
	CreateSeller(seller models.Seller) (models.Seller, error)
	GetAllSellers() ([]models.Seller, error)
	DeleteSellerById(id int) error
	UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error)
}
