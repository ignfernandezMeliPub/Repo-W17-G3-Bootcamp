package sections_repository

import "app/pkg/models"

type SectionsRepository interface {
	PoblateSectionsRepo([]models.Section) error
	GetAllSections() ([]models.Section, error)
	GetSectionById(id int) (models.Section, error)
	CreateSection(section models.Section) (models.Section, error)
	UpdateSectionById(section models.Section) (models.Section, error)
	DeleteSectionById(id int) error
}
