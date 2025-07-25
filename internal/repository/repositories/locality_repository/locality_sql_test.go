package locality_repository

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

func setupLocalityRepository(t *testing.T) (*LocalityRepositorySql, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewLocalityRepositorySql(db)

	cleanup := func() {
		db.Close()
	}

	return &repo, mock, cleanup
}

func TestLocalityRepositorySql_CreateLocality(t *testing.T) {
	repo, mock, cleanup := setupLocalityRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		locality := models.Locality{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)`)).
			WithArgs("LOC001", "Buenos Aires", "Buenos Aires", "Argentina").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.CreateLocality(locality)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, locality, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_id", func(t *testing.T) {
		// Arrange
		locality := models.Locality{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)`)).
			WithArgs("LOC001", "Buenos Aires", "Buenos Aires", "Argentina").
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'LOC001' for key 'localities.PRIMARY'",
			})

		// Act
		result, err := repo.CreateLocality(locality)

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		assert.Equal(t, locality, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		locality := models.Locality{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			ProvinceName: "Buenos Aires",
			CountryName:  "Argentina",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO localities (id, locality_name, province_name, country_name) VALUES (?, ?, ?, ?)`)).
			WithArgs("LOC001", "Buenos Aires", "Buenos Aires", "Argentina").
			WillReturnError(sql.ErrConnDone)

		// Act
		result, err := repo.CreateLocality(locality)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Equal(t, locality, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLocalityRepositorySql_GetLocalitySellerCount(t *testing.T) {
	repo, mock, cleanup := setupLocalityRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedResult := models.LocalitySellerCount{
			Id:           "LOC001",
			LocalityName: "Buenos Aires",
			SellersCount: 5,
		}

		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		}).AddRow(
			"LOC001", "Buenos Aires", 5,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id`)).
			WithArgs("LOC001").
			WillReturnRows(rows)

		// Act
		result, err := repo.GetLocalitySellerCount("LOC001")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("locality_with_zero_sellers", func(t *testing.T) {
		// Arrange
		expectedResult := models.LocalitySellerCount{
			Id:           "LOC002",
			LocalityName: "Córdoba",
			SellersCount: 0,
		}

		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		}).AddRow(
			"LOC002", "Córdoba", 0,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id`)).
			WithArgs("LOC002").
			WillReturnRows(rows)

		// Act
		result, err := repo.GetLocalitySellerCount("LOC002")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("locality_not_found", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id`)).
			WithArgs("NONEXISTENT").
			WillReturnRows(rows)

		// Act
		result, err := repo.GetLocalitySellerCount("NONEXISTENT")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrNotFound, err)
		assert.Equal(t, models.LocalitySellerCount{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id`)).
			WithArgs("LOC001").
			WillReturnError(sql.ErrConnDone)

		// Act
		result, err := repo.GetLocalitySellerCount("LOC001")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Equal(t, models.LocalitySellerCount{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLocalityRepositorySql_GetLocalitiesSellerCount(t *testing.T) {
	repo, mock, cleanup := setupLocalityRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedResults := []models.LocalitySellerCount{
			{
				Id:           "LOC001",
				LocalityName: "Buenos Aires",
				SellersCount: 5,
			},
			{
				Id:           "LOC002",
				LocalityName: "Córdoba",
				SellersCount: 3,
			},
			{
				Id:           "LOC003",
				LocalityName: "Rosario",
				SellersCount: 0,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		}).AddRow(
			"LOC001", "Buenos Aires", 5,
		).AddRow(
			"LOC002", "Córdoba", 3,
		).AddRow(
			"LOC003", "Rosario", 0,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id`)).
			WillReturnRows(rows)

		// Act
		results, err := repo.GetLocalitiesSellerCount()

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("single_locality", func(t *testing.T) {
		// Arrange
		expectedResults := []models.LocalitySellerCount{
			{
				Id:           "LOC001",
				LocalityName: "Buenos Aires",
				SellersCount: 2,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		}).AddRow(
			"LOC001", "Buenos Aires", 2,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id`)).
			WillReturnRows(rows)

		// Act
		results, err := repo.GetLocalitiesSellerCount()

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no_localities_found", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "locality_name", "sellers_count",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id`)).
			WillReturnRows(rows)

		// Act
		results, err := repo.GetLocalitiesSellerCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrNotFound, err)
		assert.Empty(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id`)).
			WillReturnError(sql.ErrConnDone)

		// Act
		results, err := repo.GetLocalitiesSellerCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
		assert.Empty(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
