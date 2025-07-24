package sections_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
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

func TestSectionsServiceImpl_GetAllSections(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("get all_ok", func(t *testing.T) {

		section, err := rp.GetAllSections()
		require.NoError(t, err)
		require.NotNil(t, section)
		require.Greater(t, len(section), 0)
	})

	t.Run("get all_not found", func(t *testing.T) {
		err := rp.DeleteSectionById(1)
		require.NoError(t, err)
		section, err := rp.GetAllSections()
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
		require.Nil(t, section)
	})
}

func TestSectionsServiceImpl_GetSectionById(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("get by id_ok", func(t *testing.T) {
		section, err := rp.GetSectionById(1)

		res := models.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: -10,
			MinimumTemperature: -20,
			CurrentCapacity:    300,
			MinimumCapacity:    100,
			MaximumCapacity:    500,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
	})

	t.Run("get by id_not found", func(t *testing.T) {
		section, err := rp.GetSectionById(2)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
		require.Equal(t, models.Section{}, section)
	})
}

func TestSectionsRepositorySQL_CreateSection(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("create_ok", func(t *testing.T) {
		req := models.Section{
			ID:                 0,
			SectionNumber:      6,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		section, err := rp.CreateSection(req)
		require.NoError(t, err)
		require.NotNil(t, section)
		res := req
		res.ID = section.ID
		require.Equal(t, res, section)
	})

	t.Run("create_conflict", func(t *testing.T) {
		req := models.Section{
			ID:                 0,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		section, err := rp.CreateSection(req)
		res := models.Section{}
		resErr := &custom_errors.UniqueAttributeViolationErr{
			AttributeName: "section_number",
			Value:         "1",
		}

		require.NotNil(t, err)
		require.Equal(t, resErr, err)
		require.Equal(t, res, section)
	})

	t.Run("create_foreignKeyViolation", func(t *testing.T) {
		req := models.Section{
			ID:                 1,
			SectionNumber:      2,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        2,
			ProductTypeId:      1,
		}

		section, err := rp.CreateSection(req)
		errExpected := &custom_errors.ForeignKeyViolationError{
			ConstraintName: "warehouse_id",
			IsParentRow:    false,
			Details:        "Cannot add or update a child row: a foreign key constraint fails (`fresh_db_test`.`sections`, CONSTRAINT `fk_sections_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE)",
		}
		require.NotNil(t, err)
		require.Equal(t, err, errExpected)
		require.Equal(t, models.Section{}, section)
	})
}

func TestSectionsRepositorySQL_UpdateSection(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("update_ok", func(t *testing.T) {
		req := models.Section{
			ID:                 1,
			SectionNumber:      6,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		section, err := rp.UpdateSectionById(req)
		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, req, section)
	})

	t.Run("update_foreignKeyViolation", func(t *testing.T) {
		req := models.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        2,
			ProductTypeId:      1,
		}

		section, err := rp.UpdateSectionById(req)
		errExpected := &custom_errors.ForeignKeyViolationError{
			ConstraintName: "warehouse_id",
			IsParentRow:    false,
			Details:        "Cannot add or update a child row: a foreign key constraint fails (`fresh_db_test`.`sections`, CONSTRAINT `fk_sections_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE)",
		}
		require.NotNil(t, err)
		require.Equal(t, err, errExpected)
		require.Equal(t, models.Section{}, section)
	})
}

func TestSectionsRepositorySQL_DeleteSection(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("delete_ok", func(t *testing.T) {

		err := rp.DeleteSectionById(1)
		require.NoError(t, err)
	})

	t.Run("delete_not found", func(t *testing.T) {
		err := rp.DeleteSectionById(2)

		require.NotNil(t, err)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
	})
}

func TestSectionsRepositorySQL_GetProductBatchBySection(t *testing.T) {
	db, err := sql.Open("txdb", "fresh_db_test")
	require.NoError(t, err)
	defer db.Close()
	rp := NewSectionsRepositorySQL(db)

	t.Run("get all_ok", func(t *testing.T) {
		section, err := rp.GetProductBatchBySection(nil)

		res := []models.ProductBatchResponse{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 500,
			},
		}

		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
	})

	t.Run("get by id_ok", func(t *testing.T) {
		id := 1
		section, err := rp.GetProductBatchBySection(&id)

		res := []models.ProductBatchResponse{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 500,
			},
		}

		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
	})

	t.Run("get by id_not found", func(t *testing.T) {
		id := 2
		res, err := rp.GetProductBatchBySection(&id)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
		require.Equal(t, []models.ProductBatchResponse(nil), res)
	})
}
