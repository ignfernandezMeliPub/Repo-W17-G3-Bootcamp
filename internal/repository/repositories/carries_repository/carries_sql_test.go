package carries_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCarriesRepository(t *testing.T) (*CarriesSql, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewCarriesSql(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestCarriesSQL_CreateCarrie(t *testing.T) {
	repo, mock, cleanup := setupCarriesRepository(t)
	defer cleanup()

	const query = `INSERT INTO carries (
				cid, 
				company_name, 
				address, 
				telephone, 
				locality_id
				) VALUES(?, ?, ?, ?, ?)`
	carrie := models.Carries{
		Cid:         "1",
		CompanyName: "Company 1",
		Address:     "Address 1",
		Telephone:   "1234567890",
		LocalityId:  "1",
	}

	t.Run("create carries success", func(t *testing.T) {
		expectedCarries := carrie
		expectedCarries.Id = 1

		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(carrie.Cid, carrie.CompanyName, carrie.Address, carrie.Telephone, carrie.LocalityId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		createdCarries, err := repo.CreateCarrie(carrie)
		require.NoError(t, err)
		require.Equal(t, expectedCarries, createdCarries)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("duplicate carries error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(carrie.Cid, carrie.CompanyName, carrie.Address, carrie.Telephone, carrie.LocalityId).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry '1' for key 'carries.cid'",
			})
		result, err := repo.CreateCarrie(carrie)
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.Equal(t, carrie, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("database error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(carrie.Cid, carrie.CompanyName, carrie.Address, carrie.Telephone, carrie.LocalityId).
			WillReturnError(sql.ErrConnDone)
		result, err := repo.CreateCarrie(carrie)
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Equal(t, carrie, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
