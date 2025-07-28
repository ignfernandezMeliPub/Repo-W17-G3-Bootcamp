package sql_utils

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFields(t *testing.T) {
	t.Run("should return error when instance is not a pointer", func(t *testing.T) {
		// Arrange
		var instance TestStruct
		fieldMap := make(map[string]reflect.Value)

		// Act
		result, err := getFields(instance, fieldMap)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "instance must be a pointer", err.Error())
		assert.Empty(t, result)
	})

	t.Run("should return error when instance is not a pointer to struct", func(t *testing.T) {
		// Arrange
		var instance int = 42
		fieldMap := make(map[string]reflect.Value)

		// Act
		result, err := getFields(&instance, fieldMap)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "instance must be a pointer to a struct", err.Error())
		assert.Empty(t, result)
	})

	t.Run("should success with valid instance", func(t *testing.T) {
		// Arrange
		instance := TestStruct{ID: 1, Name: "John", Age: 25}

		fieldMap := make(map[string]reflect.Value)
		expected := map[string]interface{}{
			"id":   instance.ID,
			"name": instance.Name,
			"age":  instance.Age,
		}

		// Act
		result, err := getFields(&instance, fieldMap)

		// Assert
		require.NoError(t, err)
		for k, v := range expected {
			gotVal, ok := result[k]
			require.True(t, ok, k)
			require.Equal(t, v, gotVal.Interface(), k)
		}
		require.Equal(t, len(expected), len(result))
	})

	t.Run("should success with valid instance with embedded struct", func(t *testing.T) {
		// Arrange
		instance := TestStructWithEmbedded{TestStruct: TestStruct{ID: 1, Name: "John", Age: 25}, Email: "john@example.com"}

		fieldMap := make(map[string]reflect.Value)
		expected := map[string]interface{}{
			"id":    instance.ID,
			"name":  instance.Name,
			"age":   instance.Age,
			"email": instance.Email,
		}

		// Act
		result, err := getFields(&instance, fieldMap)

		// Assert
		require.NoError(t, err)
		for k, v := range expected {
			gotVal, ok := result[k]
			require.True(t, ok, k)
			require.Equal(t, v, gotVal.Interface(), k)
		}
		require.Equal(t, len(expected), len(result))
	})
}
func TestInitInstanceWithRows(t *testing.T) {
	t.Run("should return error when rows.Columns() fails. Example: Rows are empty", func(t *testing.T) {
		// Arrange

		instance := &TestStruct{}
		rows := sql.Rows{}

		// Act
		err := initInstanceWithRows(&rows, instance)

		// Assert
		assert.Error(t, err)
	})

	t.Run("should return error when instance is not a pointer to struct", func(t *testing.T) {
		// Arrange
		// Usamos sqlmock para crear un *sql.Rows simulado
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()
		instance := TestStruct{}
		mockRows := sqlmock.NewRows([]string{"id", "name", "extra"}).
			AddRow(1, "John", "extra")
		mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)
		rows, err := db.Query("SELECT id, name, extra FROM users")
		require.NoError(t, err)
		defer rows.Close()

		// Act
		err = initInstanceWithRows(rows, instance)

		// Assert
		assert.Error(t, err)
	})

	t.Run("should fail if do not use rows.Next() before calling initInstanceWithRows", func(t *testing.T) {
		// Arrange
		// Usamos sqlmock para crear un *sql.Rows simulado
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()
		instance := &TestStruct{ID: 1}
		mockRows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)
		mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)
		rows, err := db.Query("SELECT id FROM users")
		require.NoError(t, err)
		defer rows.Close()

		// Act
		err = initInstanceWithRows(rows, instance)

		// Assert
		assert.Error(t, err)
	})

	t.Run("should success with struct with rows that have columns that are not in the struct", func(t *testing.T) {
		// Arrange
		// Usamos sqlmock para crear un *sql.Rows simulado
		db, mock, cleanup := setupMockDB(t)
		defer cleanup()
		expectedInstance := &TestStruct{ID: 1, Name: "John", Age: 25}
		instance := &TestStruct{}
		mockRows := sqlmock.NewRows([]string{"id", "name", "age", "extra"}).
			AddRow(1, "John", 25, "extra")
		mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)
		rows, err := db.Query("SELECT id, name, age, extra FROM users")
		require.NoError(t, err)
		defer rows.Close()

		success := rows.Next()
		require.True(t, success)

		// Act
		err = initInstanceWithRows(rows, instance)

		// Assert
		assert.NoError(t, err)
		assert.EqualValues(t, expectedInstance, instance)
	})
}
