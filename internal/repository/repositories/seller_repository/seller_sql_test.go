package seller_repository

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

func setupSellerRepository(t *testing.T) (*SellerRepositorySql, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewSellerRepositorySql(db)

	cleanup := func() {
		db.Close()
	}

	return &repo, mock, cleanup
}

func TestSellerRepositorySql_CreateSeller(t *testing.T) {
	repo, mock, cleanup := setupSellerRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		seller := models.Seller{
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		expectedSeller := seller
		expectedSeller.Id = 1

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(1001, "Test Company", "123 Test St", "555-1234", "LOC001").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.CreateSeller(seller)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_cid", func(t *testing.T) {
		// Arrange
		seller := models.Seller{
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(1001, "Test Company", "123 Test St", "555-1234", "LOC001").
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry '1001' for key 'sellers.cid'",
			})

		// Act
		result, err := repo.CreateSeller(seller)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.Equal(t, seller, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		seller := models.Seller{
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(1001, "Test Company", "123 Test St", "555-1234", "LOC001").
			WillReturnError(sql.ErrConnDone)

		// Act
		result, err := repo.CreateSeller(seller)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Equal(t, seller, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepositorySql_GetSellerById(t *testing.T) {
	repo, mock, cleanup := setupSellerRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(
			1, 1001, "Test Company", "123 Test St", "555-1234", "LOC001",
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		seller, err := repo.GetSellerById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, seller)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		_, err := repo.GetSellerById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err := repo.GetSellerById(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepositorySql_GetAllSellers(t *testing.T) {
	repo, mock, cleanup := setupSellerRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		expectedSellers := []models.Seller{
			{
				Id:          1,
				CompanyId:   1001,
				CompanyName: "Company A",
				Address:     "123 Test St",
				Telephone:   "555-1111",
				LocalityId:  "LOC001",
			},
			{
				Id:          2,
				CompanyId:   1002,
				CompanyName: "Company B",
				Address:     "456 Test Ave",
				Telephone:   "555-2222",
				LocalityId:  "LOC002",
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(
			1, 1001, "Company A", "123 Test St", "555-1111", "LOC001",
		).AddRow(
			2, 1002, "Company B", "456 Test Ave", "555-2222", "LOC002",
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers`)).
			WillReturnRows(rows)

		// Act
		sellers, err := repo.GetAllSellers()

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSellers, sellers)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no_rows", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers`)).
			WillReturnRows(rows)

		// Act
		sellers, err := repo.GetAllSellers()

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, sellers)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers`)).
			WillReturnError(sql.ErrConnDone)

		// Act
		sellers, err := repo.GetAllSellers()

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Empty(t, sellers)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepositorySql_DeleteSellerById(t *testing.T) {
	repo, mock, cleanup := setupSellerRepository(t)
	defer cleanup()

	t.Run("success", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err := repo.DeleteSellerById(1)

		// Assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("seller_not_found", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sellers WHERE id = ?`)).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		err := repo.DeleteSellerById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnError(sql.ErrConnDone)

		// Act
		err := repo.DeleteSellerById(1)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepositorySql_UpdateSellerById(t *testing.T) {
	repo, mock, cleanup := setupSellerRepository(t)
	defer cleanup()

	t.Run("success_update_all_fields", func(t *testing.T) {
		// Arrange
		companyId := 1001
		companyName := "Updated Company"
		address := "456 Updated St"
		telephone := "555-9999"

		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Updated Company",
			Address:     "456 Updated St",
			Telephone:   "555-9999",
			LocalityId:  "LOC001",
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE sellers SET cid = ?, company_name = ?, address = ?, telephone = ? WHERE id = ?`)).
			WithArgs(1001, "Updated Company", "456 Updated St", "555-9999", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetSellerById call that happens after update
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(
			1, 1001, "Updated Company", "456 Updated St", "555-9999", "LOC001",
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.UpdateSellerById(1, &companyId, &companyName, &address, &telephone)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success_update_single_field", func(t *testing.T) {
		// Arrange
		companyName := "Updated Company Name Only"

		expectedSeller := models.Seller{
			Id:          1,
			CompanyId:   1001,
			CompanyName: "Updated Company Name Only",
			Address:     "123 Test St",
			Telephone:   "555-1234",
			LocalityId:  "LOC001",
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE sellers SET company_name = ? WHERE id = ?`)).
			WithArgs("Updated Company Name Only", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Mock the GetSellerById call that happens after update
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		}).AddRow(
			1, 1001, "Updated Company Name Only", "123 Test St", "555-1234", "LOC001",
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := repo.UpdateSellerById(1, nil, &companyName, nil, nil)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSeller, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no_fields_to_update", func(t *testing.T) {
		// Arrange - no mocks needed since method should return early

		// Act
		_, err := repo.UpdateSellerById(1, nil, nil, nil, nil)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.MandatoryArgMissingErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("seller_not_found", func(t *testing.T) {
		// Arrange
		companyName := "Updated Company"

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE sellers SET company_name = ? WHERE id = ?`)).
			WithArgs("Updated Company", 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Mock the GetSellerById call that happens after update - should return not found
		rows := sqlmock.NewRows([]string{
			"id", "cid", "company_name", "address", "telephone", "locality_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?`)).
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		_, err := repo.UpdateSellerById(999, nil, &companyName, nil, nil)

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate_cid_constraint", func(t *testing.T) {
		// Arrange
		companyId := 1002 // existing company ID

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE sellers SET cid = ? WHERE id = ?`)).
			WithArgs(1002, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry '1002' for key 'sellers.cid'",
			})

		// Act
		_, err := repo.UpdateSellerById(1, &companyId, nil, nil, nil)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database_error", func(t *testing.T) {
		// Arrange
		companyName := "Updated Company"

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE sellers SET company_name = ? WHERE id = ?`)).
			WithArgs("Updated Company", 1).
			WillReturnError(sql.ErrConnDone)

		// Act
		_, err := repo.UpdateSellerById(1, nil, &companyName, nil, nil)

		// Assert
		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
