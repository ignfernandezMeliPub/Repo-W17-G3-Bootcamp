package sql_utils

import (
	"app/pkg/custom_errors"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStruct es una estructura de prueba para los tests
type TestStruct struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

// TestStructWithEmbedded es una estructura con structs embebidos para probar recursividad
type TestStructWithEmbedded struct {
	TestStruct
	Email string `db:"email"`
}

// TestStructWithoutTags es una estructura sin tags db para probar comportamiento con campos sin mapear
type TestStructWithoutTags struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Age   int    `db:"age"`
	Cache string // Sin tag db
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

func TestQuery(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("should return multiple rows successfully when query returns data", func(t *testing.T) {
		// Arrange
		expectedResults := []TestStruct{
			{ID: 1, Name: "John", Age: 25},
			{ID: 2, Name: "Jane", Age: 30},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25).
			AddRow(2, "Jane", 30)

		mock.ExpectQuery(`SELECT id, name, age FROM users`).
			WillReturnRows(rows)

		// Act
		results, err := Query[TestStruct](db, `SELECT id, name, age FROM users`, []any{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return empty slice and no rows error when query returns no data", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "age"})

		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = ?`).
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		results, err := Query[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{999})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Empty(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle struct with embedded structs successfully", func(t *testing.T) {
		// Arrange
		expectedResults := []TestStructWithEmbedded{
			{TestStruct: TestStruct{ID: 1, Name: "John", Age: 25}, Email: "john@example.com"},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "age", "email"}).
			AddRow(1, "John", 25, "john@example.com")

		mock.ExpectQuery(`SELECT id, name, age, email FROM users`).
			WillReturnRows(rows)

		// Act
		results, err := Query[TestStructWithEmbedded](db, `SELECT id, name, age, email FROM users`, []any{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle struct with fields without db tags", func(t *testing.T) {
		// Arrange
		expectedResults := []TestStructWithoutTags{
			{ID: 1, Name: "John", Age: 25, Cache: ""}, // Cache field should be empty as it has no db tag
		}

		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25)

		mock.ExpectQuery(`SELECT id, name, age FROM users`).
			WillReturnRows(rows)

		// Act
		results, err := Query[TestStructWithoutTags](db, `SELECT id, name, age FROM users`, []any{})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResults, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT id, name, age FROM users`).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		results, err := Query[TestStruct](db, `SELECT id, name, age FROM users`, []any{})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Empty(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when initInstanceWithRows fails during iteration", func(t *testing.T) {
		// Arrange
		// Create a struct that will cause scanning issues
		type InvalidStruct struct {
			ID   int    `db:"id"`
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		// Create rows with a value that can't be scanned into an int field
		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", "not_a_number"). // This will cause scanning error
			AddRow(2, "Jane", 30)

		mock.ExpectQuery(`SELECT id, name, age FROM users`).
			WillReturnRows(rows)

		// Act
		results, err := Query[InvalidStruct](db, `SELECT id, name, age FROM users`, []any{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when rows has error after scanning", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25).
			RowError(0, &mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		mock.ExpectQuery(`SELECT id, name, age FROM users`).
			WillReturnRows(rows)

		// Act
		results, err := Query[TestStruct](db, `SELECT id, name, age FROM users`, []any{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestQueryRow(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("should return single row successfully when query returns data", func(t *testing.T) {
		// Arrange
		expectedResult := TestStruct{ID: 1, Name: "John", Age: 25}

		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25)

		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := QueryRow[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return no rows error when query returns no data", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{"id", "name", "age"})

		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = \?`).
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		result, err := QueryRow[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{999})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Equal(t, TestStruct{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle struct with embedded structs successfully", func(t *testing.T) {
		// Arrange
		expectedResult := TestStructWithEmbedded{TestStruct: TestStruct{ID: 1, Name: "John", Age: 25}, Email: "john@example.com"}

		rows := sqlmock.NewRows([]string{"id", "name", "age", "email"}).
			AddRow(1, "John", 25, "john@example.com")

		mock.ExpectQuery(`SELECT id, name, age, email FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := QueryRow[TestStructWithEmbedded](db, `SELECT id, name, age, email FROM users WHERE id = ?`, []any{1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when initInstanceWithRows fails during scanning", func(t *testing.T) {
		// Arrange
		// Create a mock that will fail when calling rows.Columns()
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		// Create a custom rows mock that will fail on Columns() call
		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25).
			RowError(0, &mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		result, err := QueryRow[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{1})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, TestStruct{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		result, err := QueryRow[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{1})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, TestStruct{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestInsert(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("should return inserted id successfully when insert operation succeeds", func(t *testing.T) {
		// Arrange
		expectedID := int64(123)
		mock.ExpectExec(`INSERT INTO users \(name, age\) VALUES \(\?, \?\)`).
			WithArgs("John", 25).
			WillReturnResult(sqlmock.NewResult(expectedID, 1))

		// Act
		id, err := Insert(db, `INSERT INTO users (name, age) VALUES (?, ?)`, []any{"John", 25})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when referenced record does not exist", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`INSERT INTO users \(name, age, department_id\) VALUES \(\?, \?, \?\)`).
			WithArgs("John", 25, 999).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			})

		// Act
		id, err := Insert(db, `INSERT INTO users (name, age, department_id) VALUES (?, ?, ?)`, []any{"John", 25, 999})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, int64(0), id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestUpdate(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("should return affected rows count successfully when update operation succeeds", func(t *testing.T) {
		// Arrange
		expectedAffectedRows := int64(1)
		mock.ExpectExec(`UPDATE users SET name = \?, age = \? WHERE id = \?`).
			WithArgs("John Updated", 30, 1).
			WillReturnResult(sqlmock.NewResult(0, expectedAffectedRows))

		// Act
		affectedRows, err := Update(db, `UPDATE users SET name = ?, age = ? WHERE id = ?`, []any{"John Updated", 30, 1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedAffectedRows, affectedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return zero affected rows when no records match update criteria", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`UPDATE users SET name = \?, age = \? WHERE id = \?`).
			WithArgs("John Updated", 30, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		affectedRows, err := Update(db, `UPDATE users SET name = ?, age = ? WHERE id = ?`, []any{"John Updated", 30, 999})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, int64(0), affectedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return unique violation error when updated value violates unique constraint", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`UPDATE users SET name = \?, age = \? WHERE id = \?`).
			WithArgs("Existing Name", 30, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'Existing Name' for key 'users.name'",
			})

		// Act
		affectedRows, err := Update(db, `UPDATE users SET name = ?, age = ? WHERE id = ?`, []any{"Existing Name", 30, 1})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, int64(0), affectedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when updated reference does not exist", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`UPDATE users SET name = \?, age = \?, department_id = \? WHERE id = \?`).
			WithArgs("John Updated", 30, 999, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			})

		// Act
		affectedRows, err := Update(db, `UPDATE users SET name = ?, age = ?, department_id = ? WHERE id = ?`, []any{"John Updated", 30, 999, 1})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, int64(0), affectedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDelete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("should return deleted rows count successfully when delete operation succeeds", func(t *testing.T) {
		// Arrange
		expectedDeletedRows := int64(1)
		mock.ExpectExec(`DELETE FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, expectedDeletedRows))

		// Act
		deletedRows, err := Delete(db, `DELETE FROM users WHERE id = ?`, []any{1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedDeletedRows, deletedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return zero deleted rows when no records match delete criteria", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`DELETE FROM users WHERE id = \?`).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// Act
		deletedRows, err := Delete(db, `DELETE FROM users WHERE id = ?`, []any{999})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, int64(0), deletedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when deleted record has dependent records", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`DELETE FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{
				Number:  1451,
				Message: "Cannot delete or update a parent row: a foreign key constraint fails",
			})

		// Act
		deletedRows, err := Delete(db, `DELETE FROM users WHERE id = ?`, []any{1})

		// Assert
		assert.Error(t, err)
		assert.IsType(t, &mysql.MySQLError{}, err)
		assert.Equal(t, int64(0), deletedRows)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDBExecutorInterface(t *testing.T) {
	t.Run("should work with sql db instance", func(t *testing.T) {
		// Arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		expectedResult := TestStruct{ID: 1, Name: "John", Age: 25}
		rows := sqlmock.NewRows([]string{"id", "name", "age"}).
			AddRow(1, "John", 25)

		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = ?`).
			WithArgs(1).
			WillReturnRows(rows)

		// Act - Using *sql.DB as DBExecutor
		result, err := QueryRow[TestStruct](db, `SELECT id, name, age FROM users WHERE id = ?`, []any{1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should work with sql tx instance", func(t *testing.T) {
		// Arrange
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()

		// Configure expectations for transaction
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT id, name, age FROM users WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(1, "John", 25))
		mock.ExpectRollback()

		// Start transaction
		tx, err := db.Begin()
		require.NoError(t, err)

		expectedResult := TestStruct{ID: 1, Name: "John", Age: 25}

		// Act - Using *sql.Tx as DBExecutor
		result, err := QueryRow[TestStruct](tx, `SELECT id, name, age FROM users WHERE id = ?`, []any{1})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		// Rollback the transaction
		err = tx.Rollback()
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestHandleSqlError(t *testing.T) {
	t.Run("should return nil when error is nil", func(t *testing.T) {
		// Act
		result := HandleSqlError(nil)

		// Assert
		assert.Nil(t, result)
	})

	t.Run("should return err not found when error is sql err no rows", func(t *testing.T) {
		// Act
		result := HandleSqlError(sql.ErrNoRows)

		// Assert
		assert.Equal(t, custom_errors.ErrNotFound, result)
	})

	t.Run("should return original error when not mysql error", func(t *testing.T) {
		// Arrange
		originalError := fmt.Errorf("some error")

		// Act
		result := HandleSqlError(originalError)

		// Assert
		assert.Equal(t, originalError, result)
	})

	t.Run("should return unique violation error when mysql error 1062", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'test@example.com' for key 'users.email'",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, result)
		uniqueErr := result.(*custom_errors.UniqueAttributeViolationErr)
		assert.Equal(t, "email", uniqueErr.AttributeName)
		assert.Equal(t, "test@example.com", uniqueErr.Value)
	})

	t.Run("should return unique violation error when mysql error 1062 with simple key", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  1062,
			Message: "Duplicate entry 'test' for key 'unique_key'",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, result)
		uniqueErr := result.(*custom_errors.UniqueAttributeViolationErr)
		assert.Equal(t, "unique_key", uniqueErr.AttributeName)
		assert.Equal(t, "test", uniqueErr.Value)
	})

	t.Run("should return foreign key violation error when mysql error 1451", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  1451,
			Message: "Cannot delete or update a parent row: a foreign key constraint fails (`test_db`.`users`, CONSTRAINT `fk_user_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`))",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.IsType(t, &custom_errors.ForeignKeyViolationError{}, result)
		fkErr := result.(*custom_errors.ForeignKeyViolationError)
		assert.Equal(t, "department_id", fkErr.ConstraintName)
		assert.True(t, fkErr.IsParentRow)
		assert.Equal(t, mysqlErr.Message, fkErr.Details)
	})

	t.Run("should_return_foreign_key_violation_error_when_mysql_error_1452", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails (`test_db`.`users`, CONSTRAINT `fk_user_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`))",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.IsType(t, &custom_errors.ForeignKeyViolationError{}, result)
		fkErr := result.(*custom_errors.ForeignKeyViolationError)
		assert.Equal(t, "department_id", fkErr.ConstraintName)
		assert.False(t, fkErr.IsParentRow)
		assert.Equal(t, mysqlErr.Message, fkErr.Details)
	})

	t.Run("should_return_foreign_key_violation_error_when_mysql_error_1452_with_unknown_constraint", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  1452,
			Message: "Cannot add or update a child row: a foreign key constraint fails",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.IsType(t, &custom_errors.ForeignKeyViolationError{}, result)
		fkErr := result.(*custom_errors.ForeignKeyViolationError)
		assert.Equal(t, "unknown", fkErr.ConstraintName)
		assert.False(t, fkErr.IsParentRow)
		assert.Equal(t, mysqlErr.Message, fkErr.Details)
	})

	t.Run("should_return_original_error_when_mysql_error_number_not_handled", func(t *testing.T) {
		// Arrange
		mysqlErr := &mysql.MySQLError{
			Number:  9999,
			Message: "Some other MySQL error",
		}

		// Act
		result := HandleSqlError(mysqlErr)

		// Assert
		assert.Equal(t, mysqlErr, result)
	})
}
