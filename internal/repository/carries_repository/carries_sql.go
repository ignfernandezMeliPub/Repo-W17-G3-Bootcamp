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

const queryStore = `INSERT INTO carries (
				cid, 
				company_name, 
				address, 
				telephone, 
				locality_id
				) VALUES(?, ?, ?, ?, ?)`

//const queryGetById = ``

func (r *CarriesSql) CreateCarries(c models.Carries) (int64, error) {
	args := []any{c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId}
	return sql_utils.Insert(r.db, queryStore, args)
}
