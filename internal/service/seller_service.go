package service

import (
	"app/internal/repository/seller_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type SellerService interface {
	Create(companyId int, companyName string, address string, telephone string) (models.Seller, error)
	GetById(id int) (models.Seller, error)
	GetAll() ([]models.Seller, error)
	Delete(id int) error
	Patch(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error)
}

type SellerServiceImpl struct {
	repository seller_repository.SellerRepository
}

func NewSellerService(repository seller_repository.SellerRepository) SellerServiceImpl {
	return SellerServiceImpl{repository: repository}
}

// Create Creates a new seller
func (s *SellerServiceImpl) Create(companyId int, companyName string, address string, telephone string) (models.Seller, error) {
	isUsed, err := s.repository.CompanyIdIsUsed(companyId)
	if err != nil {
		return models.Seller{}, err
	}
	if isUsed {
		return models.Seller{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "companyId", Value: companyId}
	}

	seller := models.Seller{Id: -1, CompanyId: companyId, CompanyName: companyName, Address: address, Telephone: telephone}
	return s.repository.Save(seller)
}

// GetById retrieves a seller by its ID
func (s *SellerServiceImpl) GetById(id int) (models.Seller, error) {
	return s.repository.GetById(id)
}

// GetAll retrieves all sellers
func (s *SellerServiceImpl) GetAll() ([]models.Seller, error) {
	return s.repository.GetAll()
}

// Patch allows to patch a resource's atts
func (s *SellerServiceImpl) Patch(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error) {
	coincidence, err := s.GetById(id)
	if err != nil {
		return models.Seller{}, err
	}

	if companyId != nil {
		isUsed, err := s.repository.CompanyIdIsUsed(*companyId)
		if err != nil {
			return models.Seller{}, err
		}
		if isUsed {
			return models.Seller{}, &custom_errors.UniqueAttributeViolationErr{AttributeName: "companyId", Value: id}
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

	return s.repository.Update(coincidence)
}

// Delete removes a seller by its ID
func (s *SellerServiceImpl) Delete(id int) error {
	return s.repository.Delete(id)
}
