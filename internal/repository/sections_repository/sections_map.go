package sections_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type SectionsRepositoryImpl struct {
	sections map[int]models.Section
}

func NewSectionsRepositoryMap() *SectionsRepositoryImpl {
	return &SectionsRepositoryImpl{
		sections: map[int]models.Section{},
	}
}

func (s *SectionsRepositoryImpl) PoblateSectionsRepo(sec []models.Section) error {
	for _, section := range sec {
		s.sections[section.ID] = section
	}
	if len(s.sections) == 0 {
		s.sections = make(map[int]models.Section)
	}
	return nil
}

func (s *SectionsRepositoryImpl) GetSections() ([]models.Section, error) {
	if len(s.sections) == 0 {
		return []models.Section{}, &custom_errors.ResourceNotFoundError{}
	}
	var sections []models.Section
	for _, section := range s.sections {
		sections = append(sections, section)
	}
	return sections, nil
}

func (s *SectionsRepositoryImpl) GetSectionById(id int) (models.Section, error) {
	if section, ok := s.sections[id]; ok {
		return section, nil
	}
	return models.Section{}, &custom_errors.ResourceNotFoundError{}
}

func (s *SectionsRepositoryImpl) CreateSection(section models.Section) (models.Section, error) {
	if _, ok := s.sections[section.ID]; ok {
		return models.Section{}, &custom_errors.UniqueAttributeViolationErr{Value: section.ID, AttributeName: "id"}
	}
	for _, sec := range s.sections {
		if sec.SectionNumber == section.SectionNumber {
			return models.Section{}, &custom_errors.UniqueAttributeViolationErr{Value: section.SectionNumber, AttributeName: "section_number"}
		}
	}
	if section.ID == 0 {
		section.ID = len(s.sections) + 1
	}
	s.sections[section.ID] = section
	return section, nil
}

func (s *SectionsRepositoryImpl) UpdateSection(section models.Section) (models.Section, error) {
	if _, ok := s.sections[section.ID]; ok {
		for _, sec := range s.sections {
			if sec.SectionNumber == section.SectionNumber {
				return models.Section{}, &custom_errors.UniqueAttributeViolationErr{Value: section.SectionNumber, AttributeName: "section_number"}
			}
		}
		s.sections[section.ID] = section
		return section, nil
	}
	return models.Section{}, &custom_errors.ResourceNotFoundError{}
}

func (s *SectionsRepositoryImpl) DeleteSectionById(id int) error {
	if _, ok := s.sections[id]; ok {
		delete(s.sections, id)
		return nil
	}
	return &custom_errors.ResourceNotFoundError{}
}
