package service

import (
	"app/internal/repository/repositories/carries_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"

	"github.com/go-sql-driver/mysql"
)

type ICarriesService interface {
	CreateCarrie(c models.Carries) (models.Carries, error)
}

func NewCarriesService(rp carries_repository.CarriesRepository) *CarriesService {
	return &CarriesService{rp: rp}
}

type CarriesService struct {
	rp carries_repository.CarriesRepository
}

func (s *CarriesService) CreateCarrie(c models.Carries) (models.Carries, error) {

	carrie, err := s.rp.CreateCarrie(c)
	if err != nil {
		if isUniqueConstraintError(err) {
			return models.Carries{}, custom_errors.ErrUniqueAttributeViolationError
		}
		if isForeingKeyError(err) {
			return models.Carries{}, custom_errors.ErrInvalidArgs
		}
		return models.Carries{}, err
	}
	return carrie, nil
}

func isUniqueConstraintError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
func isForeingKeyError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1452
	}
	return false
}
