package carries_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
)

func NewCarriesSql(db *sql.DB) *CarriesSql {
	return &CarriesSql{
		db: db,
	}
}

type CarriesSql struct {
	db *sql.DB
}

// ************ const queries ************

const queryCreateCarrie = `INSERT INTO carries (
				cid, 
				company_name, 
				address, 
				telephone, 
				locality_id
				) VALUES(?, ?, ?, ?, ?)`

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

func (r *CarriesSql) CreateCarrie(c models.Carries) (models.Carries, error) {
	args := []any{c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId}
	newId, err := sql_utils.Insert(r.db, queryCreateCarrie, args)
	if err != nil {
		return c, err
	}
	c.Id = int(newId)
	return c, nil
}

func (r *CarriesSql) GetCarriesReport(localityId *string) ([]models.CarriesReport, error) {
	var query string
	var args []any

	if localityId != nil {
		query = queryGetCarriesReportById
		args = []any{*localityId}

	} else {
		query = queryGetCarriesReport
	}
	return sql_utils.Query[models.CarriesReport](r.db, query, args)
}
