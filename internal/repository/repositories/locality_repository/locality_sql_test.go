package locality_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
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
		require.Equal(t, locality, result)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.Equal(t, locality, result)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Equal(t, locality, result)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Equal(t, expectedResult, result)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Equal(t, expectedResult, result)
		require.NoError(t, mock.ExpectationsWereMet())
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
		_, err := repo.GetLocalitySellerCount("NONEXISTENT")

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id WHERE l.id = ? GROUP BY l.id`)).
			WithArgs("LOC001").
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err := repo.GetLocalitySellerCount("LOC001")

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Equal(t, expectedResults, results)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Equal(t, expectedResults, results)
		require.NoError(t, mock.ExpectationsWereMet())
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
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, results)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.id, l.locality_name, COUNT(s.id) AS sellers_count FROM sellers s RIGHT JOIN localities l ON s.locality_id = l.id GROUP BY l.id`)).
			WillReturnError(sql.ErrConnDone)

		// Act
		results, err := repo.GetLocalitiesSellerCount()

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Empty(t, results)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestLocalityRepositorySql_GetCarriesReport(t *testing.T) {
	repo, mock, cleanup := setupLocalityRepository(t)
	defer cleanup()
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

	t.Run("GetCarriesReport with localityId specified", func(t *testing.T) {
		mockRowsById := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
			AddRow("L001", "Locality One", 5)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetCarriesReportById)).
			WithArgs("L001").
			WillReturnRows(mockRowsById)

		results, err := repo.GetCarriesReport("L001")
		require.NoError(t, err)
		require.Len(t, results, 1)

		require.Equal(t, "L001", results[0].LocalityId)
		require.Equal(t, "Locality One", results[0].LocalityName)
		require.Equal(t, 5, results[0].CarriesCount)
	})

	t.Run("GetCarriesReport with empty localityId", func(t *testing.T) {
		mockRowsAll := sqlmock.NewRows([]string{"locality_id", "locality_name", "carries_count"}).
			AddRow("L001", "Locality One", 5).
			AddRow("L002", "Locality Two", 3)

		mock.ExpectQuery(regexp.QuoteMeta(queryGetCarriesReport)).
			WillReturnRows(mockRowsAll)

		results, err := repo.GetCarriesReport("")
		require.NoError(t, err)
		require.Len(t, results, 2)
	})

	t.Run("GetCarriesReport returns sql error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryGetCarriesReport)).
			WillReturnError(sql.ErrConnDone)

		results, err := repo.GetCarriesReport("")
		require.Error(t, err)
		require.Equal(t, 0, len(results))
	})
}
