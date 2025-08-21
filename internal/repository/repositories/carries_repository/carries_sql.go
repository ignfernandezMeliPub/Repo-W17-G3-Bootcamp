package carries_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/models"
	"database/sql"
	"strconv"
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
	sql_utils.LogAudit("CreateCarrie", logger.LogStatusInProgress, "Insert carrie")

	args := []any{c.Cid, c.CompanyName, c.Address, c.Telephone, c.LocalityId}
	newId, err := sql_utils.Insert(r.db, queryCreateCarrie, args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateCarrie", "Insert carrie", err)
		return c, err
	}
	c.Id = int(newId)

	sql_utils.LogAudit("CreateCarrie", logger.LogStatusSuccess, "Insert carrie. Id: "+strconv.Itoa(c.Id))
	return c, nil
}
