package product_batch_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "fresh_db_test",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestSectionsRepositorySQL_CreateSection(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewProductBatchRepositorySQL(db)

	t.Run("create_ok", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        2,
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

		section, err := rp.CreateProductBatch(req)
		require.NoError(t, err)
		require.NotNil(t, section)
		res := req
		res.ID = section.ID
		require.Equal(t, res, section)
	})

	t.Run("create_conflict", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        1001,
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

		section, err := rp.CreateProductBatch(req)
		res := models.ProductBatch{}
		resErr := &custom_errors.UniqueAttributeViolationErr{
			AttributeName: "batch_number",
			Value:         "1001",
		}

		require.NotNil(t, err)
		require.Equal(t, resErr, err)
		require.Equal(t, res, section)
	})

	t.Run("create_foreignKeyViolation", func(t *testing.T) {
		req := models.ProductBatch{
			ID:                 1,
			BatchNumber:        1002,
			CurrentQuantity:    1,
			CurrentTemperature: 1,
			DueDate:            "2021-01-01",
			InitialQuantity:    1,
			ManufacturingDate:  "2021-01-01",
			ManufacturingHour:  1,
			MinimumTemperature: 1,
			ProductId:          2,
			SectionId:          1,
		}

		section, err := rp.CreateProductBatch(req)
		errExpected := &custom_errors.ForeignKeyViolationError{
			ConstraintName: "product_id",
			IsParentRow:    false,
			Details:        "Cannot add or update a child row: a foreign key constraint fails (`fresh_db_test`.`product_batches`, CONSTRAINT `fk_product_batches_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE)",
		}
		require.NotNil(t, err)
		require.Equal(t, err, errExpected)
		require.Equal(t, models.ProductBatch{}, section)
	})
}
