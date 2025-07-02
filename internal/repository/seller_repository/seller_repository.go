package seller_repository

import (
	"app/pkg/models"
)

type SellerRepository interface {
	GetById(id int) (models.Seller, error)
	CompanyIdIsUsed(companyId int) (bool, error)
	Save(seller models.Seller) (models.Seller, error)
	GetAll() ([]models.Seller, error)
	Delete(id int) error
	Update(seller models.Seller) (models.Seller, error)
}
