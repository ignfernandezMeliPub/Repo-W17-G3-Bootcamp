package service

import (
	"app/internal/repository/seller_repository"
	"app/pkg/models"
)

type SellerService interface {
	GetAllSellers() ([]models.Seller, error)
	GetSellerById(id int) (models.Seller, error)
	CreateSeller(companyId int, companyName string, address string, telephone string) (models.Seller, error)
	UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error)
	DeleteSellerById(id int) error
}

type SellerServiceImpl struct {
	repository seller_repository.SellerRepository
}

func NewSellerService(repository seller_repository.SellerRepository) SellerServiceImpl {
	return SellerServiceImpl{repository: repository}
}

// CreateSeller Creates a new seller
func (s *SellerServiceImpl) CreateSeller(companyId int, companyName string, address string, telephone string) (models.Seller, error) {
	return s.repository.CreateSeller(models.Seller{Id: -1, CompanyId: companyId, CompanyName: companyName, Address: address, Telephone: telephone})
}

// GetSellerById retrieves a seller by its ID
func (s *SellerServiceImpl) GetSellerById(id int) (models.Seller, error) {
	return s.repository.GetSellerById(id)
}

// GetAllSellers retrieves all sellers
func (s *SellerServiceImpl) GetAllSellers() ([]models.Seller, error) {
	return s.repository.GetAllSellers()
}

// UpdateSellerById allows patching a seller's attributes
func (s *SellerServiceImpl) UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error) {
	return s.repository.UpdateSellerById(id, companyId, companyName, address, telephone)
}

// DeleteSellerById removes a seller by its ID
func (s *SellerServiceImpl) DeleteSellerById(id int) error {
	return s.repository.DeleteSellerById(id)
}
