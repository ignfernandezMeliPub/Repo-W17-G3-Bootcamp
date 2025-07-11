package sql_utils

import (
	"database/sql"
)

// QueryRow executes a query that is expected to return a single row and
// automatically maps the result to a struct instance of type T.
//
// It uses dbTagName tags to map columns to struct fields, meaning the column
// order in the SELECT statement does not need to match the field order
// in the struct.
//
// Type Parameters:
//   - T: The struct type to which the result will be mapped.
//
// Parameters:
//   - db: A pointer to the active sql.DB database connection.
//   - query: The SQL query to execute, which may contain placeholders (?).
//   - args: A slice of arguments to replace the placeholders in the query.
//
// Returns:
//   - An instance of type T populated with data from the database.
//   - An error if the query fails, the mapping fails, or if no rows are
//     returned (sql.ErrNoRows).
//
// Example of compatible struct:
//
//	type Product struct {
//	    ID          int     `db:"id"`           // Will be mapped
//	    Name        string  `db:"name"`         // Will be mapped
//	    Price       float64 `db:"price"`        // Will be mapped
//	    Description string  `db:"description"`  // Will be mapped
//	    internal    string  `db:"internal"`     // Won't be mapped (not exported)
//	    Cache       string                      // Won't be mapped (no dbTagName tag)
//	}
func QueryRow[T any](db *sql.DB, query string, args []any) (T, error) {
	var instance T

	// ? Se usa db.Query() para obtener un *sql.Rows, que nos da acceso a los nombres de las columnas.
	rows, err := db.Query(query, args...)
	if err != nil {
		return instance, err
	}
	defer rows.Close() // TODO Controlar el error que puede devolver esto

	// ? Se avanza a la primera y única esperada fila.
	success := rows.Next()

	// ? Si no hubo exito en el avance, verificamos el porque
	if !success {
		// ? Se comprueba si fue porque hubo un error.
		if err := rows.Err(); err != nil {
			return instance, err
		}

		// ? Si no hubo error, significa que la consulta no devolvió resultados.
		return instance, sql.ErrNoRows
	}

	// ? Inicializamos la instancia
	if err := initInstanceWithRows(rows, &instance); err != nil {
		return instance, err
	}

	// ? Se retorna la instancia y se chequea por última vez si hubo algún error en 'rows'.
	return instance, rows.Err()
}

// Query executes a query that can return multiple rows and maps each row to a
// struct instance, returning a slice of the results.
//
// It uses dbTagName tags to map columns to struct fields, meaning the column
// order in the SELECT statement does not need to match the field order
// in the struct.
//
// Type Parameters:
//   - T: The struct type to which each row will be mapped.
//
// Parameters:
//   - db: A pointer to the active sql.DB database connection.
//   - query: The SQL query to execute, which may contain placeholders (?).
//   - args: A slice of arguments to replace the placeholders in the query.
//
// Returns:
//   - A slice of type T containing the populated structs from each row.
//   - If the query returns no rows, it returns an empty slice and a nil error.
//   - An error if the query or the mapping process fails.
//
// Example of compatible struct:
//
//	type Product struct {
//	    ID          int     `db:"id"`           // Will be mapped
//	    Name        string  `db:"name"`         // Will be mapped
//	    Price       float64 `db:"price"`        // Will be mapped
//	    Description string  `db:"description"`  // Will be mapped
//	    internal    string  `db:"internal"`     // Won't be mapped (not exported)
//	    Cache       string                      // Won't be mapped (no dbTagName tag)
//	}
func Query[T any](db *sql.DB, query string, args []any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // TODO Controlar el error que puede devolver esto

	var results []T

	for rows.Next() {
		var instance T

		err := initInstanceWithRows(rows, &instance)
		if err != nil {
			return nil, err
		}

		results = append(results, instance)
	}

	return results, rows.Err()
}

// Insert executes an INSERT SQL statement and returns the ID of the newly inserted record.
//
// This function is designed for INSERT operations that generate auto-increment IDs.
// It uses sql.Result.LastInsertId() to retrieve the generated primary key value.
//
// Parameters:
//   - db: Pointer to the sql.DB database connection
//   - query: INSERT SQL statement to execute. Can contain placeholders (?)
//   - args: Slice of arguments to replace placeholders in the query, in order
//
// Returns:
//   - int64: The auto-generated ID of the inserted record (primary key value)
//   - error: Error if any problem occurs during execution
//
// Notes:
//   - LastInsertId() behavior depends on the database driver:
//   - MySQL: Returns the auto-increment ID of the inserted row
//   - SQLite: Returns the rowid of the inserted row
//   - PostgreSQL: May not be supported, use RETURNING clause instead
//   - If the table doesn't have an auto-increment primary key, behavior may vary
//   - For bulk inserts, only returns the ID of the last inserted row
func Insert(db *sql.DB, query string, args []any) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Update executes an UPDATE SQL statement and returns the number of rows affected.
//
// This function is designed for UPDATE operations that modify existing records.
// It uses sql.Result.RowsAffected() to determine how many rows were actually updated.
//
// Parameters:
//   - db: Pointer to the sql.DB database connection
//   - query: UPDATE SQL statement to execute. Can contain placeholders (?)
//   - args: Slice of arguments to replace placeholders in the query, in order
//
// Returns:
//   - int64: Number of rows that were actually modified by the UPDATE statement
//   - error: Error if any problem occurs during execution
func Update(db *sql.DB, query string, args []any) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Delete executes a DELETE SQL statement and returns the number of rows deleted.
//
// This function is designed for DELETE operations that remove existing records.
// It uses sql.Result.RowsAffected() to determine how many rows were actually deleted.
//
// Parameters:
//   - db: Pointer to the sql.DB database connection
//   - query: DELETE SQL statement to execute. Can contain placeholders (?)
//   - args: Slice of arguments to replace placeholders in the query, in order
//
// Returns:
//   - int64: Number of rows that were actually deleted by the DELETE statement
//   - error: Error if any problem occurs during execution
func Delete(db *sql.DB, query string, args []any) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
