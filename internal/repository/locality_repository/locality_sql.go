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

// CreateLocality Creates a new locality
func (r *LocalityRepositorySql) CreateLocality(locality models.Locality) (models.Locality, error) {
	_, err := sql_utils.Insert(r.db, "INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)", []any{locality.Id, locality.LocalityName, locality.ProvinceName, locality.CountryName})
	if err != nil {
		return locality, err
	}

	return locality, nil
}
