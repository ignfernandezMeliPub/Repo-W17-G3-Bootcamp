package locality_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
)

type LocalityRepositorySql struct {
	db *sql.DB
}

func NewLocalityRepositorySql(db *sql.DB) LocalityRepositorySql {
	return LocalityRepositorySql{db: db}
}

const queryGetCarriesReportById = `SELECT 
			l.id AS locality_id, 
			l.locality_name, 
			COUNT(c.id) AS carries_count
			FROM localities l 
			LEFT JOIN carries c ON l.id = c.locality_id
			WHERE l.id = ?
			GROUP BY l.id, l.locality_name`

const queryGetCarriesReport = `SELECT 
		l.id AS locality_id, 
		l.locality_name, 
		COUNT(c.id) AS carries_count
		FROM localities l 
		LEFT JOIN carries c ON l.id = c.locality_id
		GROUP BY l.id, l.locality_name`

// CreateLocality Creates a new locality
func (r *LocalityRepositorySql) CreateLocality(locality models.Locality) (models.Locality, error) {
	_, err := sql_utils.Insert(r.db, "INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)", []any{locality.Id, locality.LocalityName, locality.ProvinceName, locality.CountryName})
	if err != nil {
		return locality, err
	}

	return locality, nil
}

func (r *LocalityRepositorySql) GetCarriesReport(localityId string) ([]models.CarriesReport, error) {
	var query string
	var args []any

	if localityId != "" {
		query = queryGetCarriesReportById
		args = []any{localityId}

	} else {
		query = queryGetCarriesReport
	}
	return sql_utils.Query[models.CarriesReport](r.db, query, args)
}
