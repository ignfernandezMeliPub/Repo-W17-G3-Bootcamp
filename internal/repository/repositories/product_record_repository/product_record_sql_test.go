package product_record_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setupProductRecordRepository(t *testing.T) (*ProductRecordRepositorySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return NewProductRecordRepositorySQL(db), mock, func() { db.Close() }
}

func TestProductRecordRepositorySQL_GetAllProductRecords(t *testing.T) {
	repo, mock, cleanup := setupProductRecordRepository(t)
	defer cleanup()

	t.Run("should return all product records successfully", func(t *testing.T) {
		// Arrange
		expectedProductRecords := []models.ProductRecord{
			{ID: 1, LastUpdateDate: "2021-01-01", PurchasePrice: 100, SalePrice: 150, ProductID: 1},
		}

		// Act
		mock.ExpectQuery("SELECT id, last_update_date, purchase_price, sale_price, product_id FROM product_records").WillReturnRows(sqlmock.NewRows([]string{"id", "last_update_date", "purchase_price", "sale_price", "product_id"}).
			AddRow(1, "2021-01-01", 100, 150, 1))

		// Assert
		productRecords, err := repo.GetAllProductRecords()
		require.NoError(t, err)
		require.Equal(t, expectedProductRecords, productRecords)
		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return not found error when database has no product records", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "last_update_date", "purchase_price", "sale_price", "product_id"})
		mock.ExpectQuery("SELECT id, last_update_date, purchase_price, sale_price, product_id FROM product_records").WillReturnRows(rows)

		// Act
		_, err := repo.GetAllProductRecords()
		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRecordRepositorySQL_CreateProductRecord(t *testing.T) {
	repo, mock, cleanup := setupProductRecordRepository(t)
	defer cleanup()

	t.Run("should create a product record successfully", func(t *testing.T) {
		// Arrange
		expectedProductRecord := models.ProductRecord{ID: 1, LastUpdateDate: "2021-01-01", PurchasePrice: 100, SalePrice: 150, ProductID: 1}

		// Act
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
			WithArgs(expectedProductRecord.LastUpdateDate, expectedProductRecord.PurchasePrice, expectedProductRecord.SalePrice, expectedProductRecord.ProductID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Assert
		productRecord, err := repo.CreateProductRecord(expectedProductRecord)
		require.NoError(t, err)
		require.Equal(t, expectedProductRecord, productRecord)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when product id does not exist", func(t *testing.T) {
		// Arrange
		expectedProductRecord := models.ProductRecord{ID: 1, LastUpdateDate: "2021-01-01", PurchasePrice: 100, SalePrice: 150, ProductID: 1}

		// Act
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES (?, ?, ?, ?)")).
			WithArgs(expectedProductRecord.LastUpdateDate, expectedProductRecord.PurchasePrice, expectedProductRecord.SalePrice, expectedProductRecord.ProductID).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails (`fresh_db`.`product_records`, CONSTRAINT `product_records_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`))",
			})

		// Assert
		productRecord, err := repo.CreateProductRecord(expectedProductRecord)
		require.Error(t, err)
		require.IsType(t, custom_errors.ErrForeignKeyViolation, err)
		require.Equal(t, models.ProductRecord{}, productRecord)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}
