package product_batch_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func setupSectionsRepository(t *testing.T) (*ProductBatchRepositorySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewProductBatchRepositorySQL(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestSectionsRepositorySQL_CreateProductBatch(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("create_ok", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2021-01-01",
			InitialQuantity:    1,
			ManufacturingDate:  "2021-01-01",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductId:          1,
			SectionId:          1,
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `product_batches` (`batch_number`,`current_quantity`,`current_temperature`,`due_date`,`initial_quantity`,`manufacturing_date`,`manufacturing_hour`,`minimum_temperature`,`product_id`,`section_id`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
			WithArgs(1, 1, 1, "2021-01-01", 1, "2021-01-01", 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		prodBatch, err := rp.CreateProductBatch(req)
		require.NoError(t, err)
		require.NotNil(t, prodBatch)
		require.Equal(t, req, prodBatch)
	})

	t.Run("create_conflict", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2021-01-01",
			InitialQuantity:    1,
			ManufacturingDate:  "2021-01-01",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductId:          1,
			SectionId:          1,
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `product_batches` (`batch_number`,`current_quantity`,`current_temperature`,`due_date`,`initial_quantity`,`manufacturing_date`,`manufacturing_hour`,`minimum_temperature`,`product_id`,`section_id`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
			WithArgs(1, 1, 1, "2021-01-01", 1, "2021-01-01", 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'section_number' for key '1'"})

		section, err := rp.CreateProductBatch(req)
		res := models.ProductBatch{}
		resErr := &custom_errors.UniqueAttributeViolationErr{}

		require.NotNil(t, err)
		require.IsType(t, resErr, err)
		require.Equal(t, res, section)
	})

	t.Run("create_foreignKeyViolation", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        1,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2021-01-01",
			InitialQuantity:    1,
			ManufacturingDate:  "2021-01-01",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductId:          1,
			SectionId:          1,
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `product_batches` (`batch_number`,`current_quantity`,`current_temperature`,`due_date`,`initial_quantity`,`manufacturing_date`,`manufacturing_hour`,`minimum_temperature`,`product_id`,`section_id`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
			WithArgs(1, 1, 1, "2021-01-01", 1, "2021-01-01", 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{Number: 1451, Message: "FOREIGN KEY \\(`product_id`\\)"})

		section, err := rp.CreateProductBatch(req)
		errExpected := &custom_errors.ForeignKeyViolationError{}
		require.NotNil(t, err)
		require.IsType(t, err, errExpected)
		require.Equal(t, models.ProductBatch{}, section)
	})
}
