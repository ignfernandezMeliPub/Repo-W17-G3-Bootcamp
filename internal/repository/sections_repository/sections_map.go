package sections_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type SectionsRepositoryImpl struct {
	sections map[int]models.Section
	seed     int
}

func NewSectionsRepositoryMap() *SectionsRepositoryImpl {
	return &SectionsRepositoryImpl{
		sections: map[int]models.Section{},
		seed:     1,
	}
}

func (s *SectionsRepositoryImpl) PoblateSectionsRepo(sec []models.Section) error {
	for _, section := range sec {
		s.sections[section.ID] = section
		s.seed++
	}
	if len(s.sections) == 0 {
		s.sections = make(map[int]models.Section)
	}
	return nil
}

func (s *SectionsRepositoryImpl) GetAllSections() ([]models.Section, error) {
	if len(s.sections) == 0 {
		return []models.Section{}, &custom_errors.ResourceNotFoundError{}
	}
	sections := make([]models.Section, len(s.sections))
	i := 0
	for _, section := range s.sections {
		sections[i] = section
		i++
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
	if section.ID == 0 {
		section.ID = s.seed
	}
	s.sections[section.ID] = section
	s.seed++
	return section, nil
}

func (s *SectionsRepositoryImpl) UpdateSectionById(section models.Section) (models.Section, error) {
	s.sections[section.ID] = section
	return section, nil
}

func (s *SectionsRepositoryImpl) DeleteSectionById(id int) error {
	delete(s.sections, id)
	return nil
}
