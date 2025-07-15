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

// GetLocalitySellerCount Returns LocalitySellerCount for received localityId
func (r *LocalityRepositorySql) GetLocalitySellerCount(localityId string) (models.LocalitySellerCount, error) {
	return sql_utils.QueryRow[models.LocalitySellerCount](r.db, "SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id", []any{localityId})
}

// GetLocalitiesSellerCount Returns a []LocalitySellerCount with an LocalitySellerCount for every locality
func (r *LocalityRepositorySql) GetLocalitiesSellerCount() ([]models.LocalitySellerCount, error) {
	return sql_utils.Query[models.LocalitySellerCount](r.db, "SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id", []any{})
}
