package service

import (
	"app/internal/repository/repositories/sections_repository"
	"app/pkg/models"
)

type SectionsService interface {
	GetAllSections() ([]models.Section, error)
	GetSectionById(id int) (models.Section, error)
	CreateSection(section models.SectionRequest) (models.Section, error)
	UpdateSectionById(id int, section models.SectionRequest) (models.Section, error)
	DeleteSectionById(id int) error

	GetAllProductBatchesBySection() (prods []models.ProductBatchResponse, err error)
	GetProductBatchBySectionId(sectionId int) (prod models.ProductBatchResponse, err error)
}

type SectionsServiceImpl struct {
	rp sections_repository.SectionsRepository
}

func NewSectionsService(rp sections_repository.SectionsRepository) *SectionsServiceImpl {
	return &SectionsServiceImpl{rp: rp}
}

func (s *SectionsServiceImpl) GetAllSections() ([]models.Section, error) {
	return s.rp.GetAllSections()
}

func (s *SectionsServiceImpl) GetSectionById(id int) (models.Section, error) {
	return s.rp.GetSectionById(id)
}

func (s *SectionsServiceImpl) CreateSection(section models.SectionRequest) (models.Section, error) {

	newSection := models.Section{
		SectionNumber:      *section.SectionNumber,
		CurrentTemperature: *section.CurrentTemperature,
		MinimumTemperature: *section.MinimumTemperature,
		CurrentCapacity:    *section.CurrentCapacity,
		MinimumCapacity:    *section.MinimumCapacity,
		MaximumCapacity:    *section.MaximumCapacity,
		WarehouseId:        *section.WarehouseId,
		ProductTypeId:      *section.ProductTypeId,
	}
	return s.rp.CreateSection(newSection)
}

func (s *SectionsServiceImpl) UpdateSectionById(id int, section models.SectionRequest) (models.Section, error) {
	oldSec, err := s.rp.GetSectionById(id)
	if err != nil {
		return models.Section{}, err
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
		oldSec.WarehouseId = *section.WarehouseId
	}
	if section.ProductTypeId != nil {
		oldSec.ProductTypeId = *section.ProductTypeId
	}
	return s.rp.UpdateSectionById(oldSec)
}

func (s *SectionsServiceImpl) DeleteSectionById(id int) error {
	return s.rp.DeleteSectionById(id)
}

func (s *SectionsServiceImpl) GetAllProductBatchesBySection() (prods []models.ProductBatchResponse, err error) {
	return s.rp.GetAllProductBatchesBySection()
}

func (s *SectionsServiceImpl) GetProductBatchBySectionId(sectionId int) (prod models.ProductBatchResponse, err error) {
	return s.rp.GetProductBatchBySectionId(sectionId)
}
