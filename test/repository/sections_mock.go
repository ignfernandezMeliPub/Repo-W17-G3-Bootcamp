package repository

import (
	"app/pkg/models"
	"github.com/stretchr/testify/mock"
)

type SectionsMock struct {
	mock.Mock
}

func NewSectionsMock() *SectionsMock {
	return &SectionsMock{}
}

func (m *SectionsMock) GetAllSections() ([]models.Section, error) {
	args := m.Called()
	return args.Get(0).([]models.Section), args.Error(1)
}

func (m *SectionsMock) GetSectionById(id int) (models.Section, error) {
	args := m.Called(id)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionsMock) CreateSection(section models.Section) (models.Section, error) {
	args := m.Called(section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionsMock) UpdateSectionById(section models.Section) (models.Section, error) {
	args := m.Called(section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *SectionsMock) DeleteSectionById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *SectionsMock) GetProductBatchBySection(sectionId *int) (prod []models.ProductBatchResponse, err error) {
	args := m.Called(sectionId)
	return args.Get(0).([]models.ProductBatchResponse), args.Error(1)
}
