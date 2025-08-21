package locality_repository

import (
	"app/internal/logger"
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
	sql_utils.LogAudit("CreateLocality", logger.LogStatusInProgress, "Insert locality")

	_, err := sql_utils.Insert(r.db, "INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)", []any{locality.Id, locality.LocalityName, locality.ProvinceName, locality.CountryName})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateLocality", "Insert locality", err)
		return locality, err
	}

	sql_utils.LogAudit("CreateLocality", logger.LogStatusSuccess, "Insert locality. Id: "+locality.Id)
	return locality, nil
}

// GetCarriesReport Returns the carries report for a locality
func (r *LocalityRepositorySql) GetCarriesReport(localityId string) ([]models.CarriesReport, error) {
	if localityId != "" {
		sql_utils.Log("GetCarriesReport", logger.LogStatusInProgress, "Select carries report by locality id: "+localityId)
	} else {
		sql_utils.Log("GetCarriesReport", logger.LogStatusInProgress, "Select carries report")
	}

	var query string
	var args []any

	if localityId != "" {
		query = queryGetCarriesReportById
		args = []any{localityId}

	} else {
		query = queryGetCarriesReport
	}
	carriesReports, err := sql_utils.Query[models.CarriesReport](r.db, query, args)
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetCarriesReport", "Select carries report", err)
		return carriesReports, err
	}

	if localityId != "" {
		sql_utils.Log("GetCarriesReport", logger.LogStatusSuccess, "Select carries report by locality id: "+localityId)
	} else {
		sql_utils.Log("GetCarriesReport", logger.LogStatusSuccess, "Select carries report")
	}
	return carriesReports, nil
}

// GetLocalitySellerCount Returns LocalitySellerCount for received localityId
func (r *LocalityRepositorySql) GetLocalitySellerCount(localityId string) (models.LocalitySellerCount, error) {
	sql_utils.Log("GetLocalitySellerCount", logger.LogStatusInProgress, "Select locality seller count by id: "+localityId)

	localitySellerCount, err := sql_utils.QueryRow[models.LocalitySellerCount](r.db, "SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id", []any{localityId})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetLocalitySellerCount", "Select locality seller count by id: "+localityId, err)
		return localitySellerCount, err
	}

	sql_utils.Log("GetLocalitySellerCount", logger.LogStatusSuccess, "Select locality seller count by id: "+localityId)
	return localitySellerCount, nil
}

// GetLocalitiesSellerCount Returns a []LocalitySellerCount with an LocalitySellerCount for every locality
func (r *LocalityRepositorySql) GetLocalitiesSellerCount() ([]models.LocalitySellerCount, error) {
	sql_utils.Log("GetLocalitiesSellerCount", logger.LogStatusInProgress, "Select localities seller count")

	localitySellerCounts, err := sql_utils.Query[models.LocalitySellerCount](r.db, "SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id", []any{})
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetLocalitiesSellerCount", "Select localities seller count", err)
		return localitySellerCounts, err
	}

	sql_utils.Log("GetLocalitiesSellerCount", logger.LogStatusSuccess, "Select localities seller count")
	return localitySellerCounts, nil
}
