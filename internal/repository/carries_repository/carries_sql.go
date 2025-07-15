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

func (r *CarriesSql) CreateCarrie(c models.Carries) (models.Carries, error) {
	args := []any{c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId}
	newId, err := sql_utils.Insert(r.db, queryCreateCarrie, args)
	if err != nil {
		return c, err
	}
	c.Id = int(newId)
	return c, nil
}
