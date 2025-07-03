package service

import (
	"app/internal/repository/seller_repository"
	"app/pkg/custom_errors"
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
	isUsed, err := s.repository.CompanyIdIsUsed(companyId)
	if err != nil {
		return models.Seller{}, err
	}
	if isUsed {
		return models.Seller{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "companyId", Value: companyId}
	}

	seller := models.Seller{Id: -1, CompanyId: companyId, CompanyName: companyName, Address: address, Telephone: telephone}
	return s.repository.CreateSeller(seller)
}

// GetSellerById retrieves a seller by its ID
func (s *SellerServiceImpl) GetSellerById(id int) (models.Seller, error) {
	return s.repository.GetSellerById(id)
}

// GetAllSellers retrieves all sellers
func (s *SellerServiceImpl) GetAllSellers() ([]models.Seller, error) {
	return s.repository.GetAllSellers()
}

// UpdateSellerById allows to patch a resource's atts

func (s *SellerServiceImpl) UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error) {
	coincidence, err := s.GetSellerById(id)
	if err != nil {
		return models.Seller{}, err
	}

	if companyId != nil && *companyId != coincidence.CompanyId {
		println(coincidence.CompanyId)
		println(*companyId)
		println(coincidence.CompanyId == *companyId)
		isUsed, err := s.repository.CompanyIdIsUsed(*companyId)
		if err != nil {
			return models.Seller{}, err
		}
		if isUsed {
			return models.Seller{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "companyId", Value: *companyId}
		}

		coincidence.CompanyId = *companyId
	}

	if companyName != nil {
		coincidence.CompanyName = *companyName
	}

	if address != nil {
		coincidence.Address = *address
	}

	if telephone != nil {
		coincidence.Telephone = *telephone
	}

	return s.repository.UpdateSellerById(coincidence)
}

// DeleteSellerById removes a seller by its ID
func (s *SellerServiceImpl) DeleteSellerById(id int) error {
	return s.repository.DeleteSellerById(id)
}
