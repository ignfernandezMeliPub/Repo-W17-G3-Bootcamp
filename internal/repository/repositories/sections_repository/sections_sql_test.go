package sections_repository

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

func setupSectionsRepository(t *testing.T) (*SectionsRepositorySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewSectionsRepositorySQL(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}
func TestSectionsServiceImpl_GetAllSections(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("get all_ok", func(t *testing.T) {

		expected := []models.Section{
			{
				ID:                 1,
				SectionNumber:      1,
				CurrentTemperature: 1,
				MinimumTemperature: 1,
				CurrentCapacity:    1,
				MinimumCapacity:    1,
				MaximumCapacity:    1,
				WarehouseId:        1,
				ProductTypeId:      1,
			},
			{
				ID:                 2,
				SectionNumber:      2,
				CurrentTemperature: 2,
				MinimumTemperature: 2,
				CurrentCapacity:    2,
				MinimumCapacity:    2,
				MaximumCapacity:    2,
				WarehouseId:        2,
				ProductTypeId:      2,
			},
		}
		rows := sqlmock.NewRows([]string{
			"id", "section_number", "current_temperature", "minimum_temperature", "current_capacity",
			"minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id",
		}).AddRow(
			1, 1, 1, 1, 1, 1, 1, 1, 1,
		).AddRow(
			2, 2, 2, 2, 2, 2, 2, 2, 2,
		)

		mock.ExpectQuery("SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections").
			WillReturnRows(rows)

		section, err := rp.GetAllSections()
		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, expected, section)
	})

	t.Run("get all_not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections").
			WillReturnError(sql.ErrNoRows)
		section, err := rp.GetAllSections()
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
		require.Nil(t, section)
	})
}

func TestSectionsServiceImpl_GetSectionById(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("get by id_ok", func(t *testing.T) {
		expected := models.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}
		rows := sqlmock.NewRows([]string{
			"id", "section_number", "current_temperature", "minimum_temperature", "current_capacity",
			"minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id",
		}).AddRow(
			1, 1, 1, 1, 1, 1, 1, 1, 1,
		)

		mock.ExpectQuery("SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnRows(rows)

		section, err := rp.GetSectionById(1)
		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, expected, section)
	})

	t.Run("get by id_not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections WHERE id = ?").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		_, err := rp.GetSectionById(2)
		require.Error(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSectionsRepositorySQL_CreateSection(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("create_ok", func(t *testing.T) {
		res := models.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sections` (`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`) VALUES (?,?,?,?,?,?,?,?)")).
			WithArgs(1, float64(1), float64(1), 1, 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		section, err := rp.CreateSection(res)
		require.NoError(t, err)
		require.NotNil(t, section)
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

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sections` (`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`) VALUES (?,?,?,?,?,?,?,?)")).
			WithArgs(1, float64(1), float64(1), 1, 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'section_number' for key '1'"})

		_, err := rp.CreateSection(req)
		resErr := &custom_errors.UniqueAttributeViolationErr{}

		require.NotNil(t, err)
		require.IsType(t, resErr, err)
	})

	t.Run("create_foreignKeyViolation", func(t *testing.T) {
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

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sections` (`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`) VALUES (?,?,?,?,?,?,?,?)")).
			WithArgs(1, float64(1), float64(1), 1, 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{Number: 1451, Message: "FOREIGN KEY \\(`warehouse_id`\\)"})

		_, err := rp.CreateSection(req)
		errExpected := &custom_errors.ForeignKeyViolationError{}
		require.NotNil(t, err)
		require.IsType(t, err, errExpected)
	})
}

func TestSectionsRepositorySQL_UpdateSection(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("update_ok", func(t *testing.T) {
		res := models.Section{
			ID:                 1,
			SectionNumber:      1,
			CurrentTemperature: 1,
			MinimumTemperature: 1,
			CurrentCapacity:    1,
			MinimumCapacity:    1,
			MaximumCapacity:    1,
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE sections SET
			section_number = ?,
			current_temperature = ?,
			minimum_temperature = ?,
			current_capacity = ?,
			minimum_capacity = ?,
			maximum_capacity = ?,
			warehouse_id = ?,
			product_type_id = ?
			WHERE id = ?`)).
			WithArgs(1, float64(1), float64(1), 1, 1, 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		section, err := rp.UpdateSectionById(res)
		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
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
			WarehouseId:        1,
			ProductTypeId:      1,
		}

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE sections SET
			section_number = ?,
			current_temperature = ?,
			minimum_temperature = ?,
			current_capacity = ?,
			minimum_capacity = ?,
			maximum_capacity = ?,
			warehouse_id = ?,
			product_type_id = ?
			WHERE id = ?`)).
			WithArgs(1, float64(1), float64(1), 1, 1, 1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(&mysql.MySQLError{Number: 1451, Message: "FOREIGN KEY \\(`warehouse_id`\\)"})

		_, err := rp.UpdateSectionById(req)
		errExpected := &custom_errors.ForeignKeyViolationError{}
		require.NotNil(t, err)
		require.IsType(t, err, errExpected)
	})
}

func TestSectionsRepositorySQL_DeleteSection(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("delete_ok", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sections WHERE id = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := rp.DeleteSectionById(1)
		require.NoError(t, err)
	})

	t.Run("delete_not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sections WHERE id = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0)).
			WillReturnError(sql.ErrNoRows)

		err := rp.DeleteSectionById(1)

		require.NotNil(t, err)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
	})

	t.Run("delete_no rows affected", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM sections WHERE id = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := rp.DeleteSectionById(1)

		require.NotNil(t, err)
		require.ErrorIs(t, err, &custom_errors.ResourceNotFoundError{})
	})
}

func TestSectionsRepositorySQL_GetProductBatchBySection(t *testing.T) {
	rp, mock, cleanup := setupSectionsRepository(t)
	defer cleanup()

	t.Run("get all_ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"section_id", "section_number", "products_count",
		}).AddRow(
			1, 1, 100,
		).AddRow(
			2, 2, 200,
		)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id GROUP BY section_id")).
			WillReturnRows(rows)

		section, err := rp.GetProductBatchBySection(nil)

		res := []models.ProductBatchResponse{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 100,
			},
			{
				SectionID:     2,
				SectionNumber: 2,
				ProductsCount: 200,
			},
		}

		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
	})

	t.Run("get by id_ok", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"section_id", "section_number", "products_count",
		}).AddRow(
			1, 1, 100,
		)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id  WHERE section_id = ? GROUP BY section_id")).
			WithArgs(1).
			WillReturnRows(rows)
		id := 1
		section, err := rp.GetProductBatchBySection(&id)

		res := []models.ProductBatchResponse{
			{
				SectionID:     1,
				SectionNumber: 1,
				ProductsCount: 100,
			},
		}

		require.NoError(t, err)
		require.NotNil(t, section)
		require.Equal(t, res, section)
	})

	t.Run("get by id_not found", func(t *testing.T) {

		mock.ExpectQuery(regexp.QuoteMeta("SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id  WHERE section_id = ? GROUP BY section_id")).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)
		id := 1
		section, err := rp.GetProductBatchBySection(&id)
		require.Error(t, err)
		require.Equal(t, err, &custom_errors.ResourceNotFoundError{})
		require.Equal(t, []models.ProductBatchResponse(nil), section)
	})
}
