package carries_repository

import "database/sql"

func NewCarriesSql(db *sql.DB) *CarriesSql {
	return &CarriesSql{
		db: db,
	}
}

type CarriesSql struct {
	db *sql.DB
}
