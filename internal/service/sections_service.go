package service

import (
	"app/internal/repository/sections_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"strconv"
)

type SectionsService interface {
	GetSections() ([]models.Section, error)
	GetSectionByID(id int) (models.Section, error)
	CreateSection(section models.Section) (models.Section, error)
	UpdateSection(id int, section models.SectionPatch) (models.Section, error)
	DeleteSection(id int) error
}

type SectionsServiceImpl struct {
	rp              sections_repository.SectionsRepository
	sv_warehouse    IWarehouseService
	sv_product_type ProductTypeServices
}

func NewSectionsService(rp sections_repository.SectionsRepository, wh IWarehouseService, pt ProductTypeServices) *SectionsServiceImpl {
	return &SectionsServiceImpl{rp: rp, sv_warehouse: wh, sv_product_type: pt}
}

func (s *SectionsServiceImpl) GetSections() ([]models.Section, error) {
	return s.rp.GetSections()
}

func (s *SectionsServiceImpl) GetSectionByID(id int) (models.Section, error) {
	return s.rp.GetSectionById(id)
}

func (s *SectionsServiceImpl) CreateSection(section models.Section) (models.Section, error) {
	if section.ID != 0 {
		_, err := s.rp.GetSectionById(section.ID)
		if err == nil {
			return models.Section{}, &custom_errors.ResourceConflictError{Argument: "id", Value: strconv.Itoa(section.ID)}
		}
	}
	err := s.ValidateSection(section.ID, section.SectionNumber)
	if err != nil {
		return models.Section{}, err
	}
	_, err = s.sv_warehouse.FindWarehouseById(section.WarehouseId)
	if err != nil {
		return models.Section{}, &custom_errors.InvalidArgValueErr{Argument: "warehouse_id", Value: section.WarehouseId}
	}
	_, err = s.sv_product_type.GetProductTypeById(section.ProductTypeId)
	if err != nil {
		return models.Section{}, &custom_errors.InvalidArgValueErr{Argument: "product_type_id", Value: section.ProductTypeId}
	}
	return s.rp.CreateSection(section)
}

func (s *SectionsServiceImpl) UpdateSection(id int, section models.SectionPatch) (models.Section, error) {
	oldSec, err := s.rp.GetSectionById(id)
	if err != nil {
		return models.Section{}, err
	}
	if section.SectionNumber != nil {
		err = s.ValidateSection(id, *section.SectionNumber)
		if err != nil {
			return models.Section{}, err
		}
	}

	if section.SectionNumber != nil {
		oldSec.SectionNumber = *section.SectionNumber
	}
	if section.CurrentTemperature != nil {
		oldSec.CurrentTemperature = *section.CurrentTemperature
	}
	if section.MinimumTemperature != nil {
		oldSec.MinimumTemperature = *section.MinimumTemperature
	}
	if section.CurrentCapacity != nil {
		oldSec.CurrentCapacity = *section.CurrentCapacity
	}
	if section.MinimumCapacity != nil {
		oldSec.MinimumCapacity = *section.MinimumCapacity
	}
	if section.MaximumCapacity != nil {
		oldSec.MaximumCapacity = *section.MaximumCapacity
	}
	if section.WarehouseId != nil {
		_, err = s.sv_warehouse.FindWarehouseById(*section.WarehouseId)
		if err != nil {
			return models.Section{}, &custom_errors.InvalidArgValueErr{Argument: "warehouse_id", Value: *section.WarehouseId}
		}
		oldSec.WarehouseId = *section.WarehouseId
	}
	if section.ProductTypeId != nil {
		_, err = s.sv_product_type.GetProductTypeById(*section.ProductTypeId)
		if err != nil {
			return models.Section{}, &custom_errors.InvalidArgValueErr{Argument: "product_type_id", Value: *section.ProductTypeId}
		}
		oldSec.ProductTypeId = *section.ProductTypeId
	}

	return s.rp.UpdateSection(oldSec)
}

func (s *SectionsServiceImpl) DeleteSection(id int) error {
	_, err := s.rp.GetSectionById(id)
	if err != nil {
		return err
	}
	return s.rp.DeleteSectionById(id)
}

func (s *SectionsServiceImpl) ValidateSection(id int, sectionNumber string) error {
	sections, err := s.rp.GetSections()
	if err != nil {
		return err
	}
	for _, section := range sections {
		if section.SectionNumber == sectionNumber && id != section.ID {
			return &custom_errors.ResourceConflictError{Argument: "section_number", Value: sectionNumber}
		}
	}
	return nil
}
