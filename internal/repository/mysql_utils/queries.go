package mysql_utils

import (
	"database/sql"
)

// QueryRow executes a SQL query that returns a single row and automatically maps
// the result to a Go struct using 'db' tags.
//
// Type Parameters:
//   - T: The struct type to map the result to. Must be a struct with fields
//     that have 'db:"column_name"' tags.
//
// Parameters:
//   - db: Pointer to the sql.DB database connection
//   - query: SQL query to execute. Can contain placeholders (?)
//   - args: Variable arguments to replace placeholders in the query
//
// Returns:
//   - T: An instance of type T with fields populated from the database
//   - error: Error if any problem occurs during execution or mapping
//
// Notes:
//   - Only fields with 'db' tags that are exported (capitalized) are mapped
//   - The order of fields in SELECT must match the order of 'db' tags in the struct
//   - If no row is found, returns sql.ErrNoRows
//   - Column order in SELECT must match the order of 'db' tags in struct T
//
// Example of compatible struct:
//
//	type Product struct {
//	    ID          int     `db:"id"`           // Will be mapped
//	    Name        string  `db:"name"`         // Will be mapped
//	    Price       float64 `db:"price"`        // Will be mapped
//	    Description string  `db:"description"`  // Will be mapped
//	    internal    string  `db:"internal"`     // Won't be mapped (not exported)
//	    Cache       string                      // Won't be mapped (no db tag)
//	}
func QueryRow[T any](db *sql.DB, query string, args []any) (T, error) {
	var instance T
	row := db.QueryRow(query, args...)
	return instance, initializeInstanceFromScanner(row, &instance)
}

// Query executes a SQL SELECT statement that returns multiple rows and automatically
// maps each result row to a Go struct using 'db' tags.
//
// This function uses reflection to automatically populate struct fields based on
// their 'db' tags, eliminating the need for manual scanning of database columns.
// It's designed for SELECT queries that return zero or more rows.
//
// Type Parameters:
//   - T: The struct type to map each result row to. Must be a struct with fields
//     that have 'db:"column_name"' tags corresponding to the SELECT columns.
//
// Parameters:
//   - db: Pointer to the sql.DB database connection
//   - query: SELECT SQL statement to execute. Can contain placeholders (?)
//   - args: Slice of arguments to replace placeholders in the query, in order
//
// Returns:
//   - []T: A slice containing instances of type T, each populated from a database row.
//     Returns empty slice if no rows match the query conditions.
//   - error: Error if any problem occurs during execution, iteration, or mapping
//
// Struct Requirements:
//
//   - Fields must be exported (start with uppercase letter)
//
//   - Fields must have 'db' tags matching SELECT column names
//
//   - Field types must be compatible with database column types
//
//   - Important! Attributes order in the struct must match column order. Otherwise, it will fail. Example:
//
//     // CORRECT: Struct field order matches SELECT column order
//     type User struct {
//     ID   int    `db:"id"`     // 1st column in SELECT
//     Name string `db:"name"`   // 2nd column in SELECT
//     Age  int    `db:"age"`    // 3rd column in SELECT
//     }
//     // SELECT id, name, age FROM users  -- Same order!
//
//     // WRONG: Struct field order doesn't match SELECT column order
//     type User struct {
//     Name string `db:"name"`   // This will receive 'id' value
//     ID   int    `db:"id"`     // This will receive 'name' value
//     Age  int    `db:"age"`    // This will receive 'age' value
//     }
//     // SELECT id, name, age FROM users  -- Different order!
//
// Example of compatible struct:
//
//	type Product struct {
//	    ID          int     `db:"id"`           // Will be mapped
//	    Name        string  `db:"name"`         // Will be mapped
//	    Price       float64 `db:"price"`        // Will be mapped
//	    Description string  `db:"description"`  // Will be mapped
//	    internal    string  `db:"internal"`     // Won't be mapped (not exported)
//	    Cache       string                      // Won't be mapped (no db tag)
//	}
//
// Behavior:
//   - Returns empty slice (not nil) if query returns no rows
//   - Processes all rows returned by the query
//   - Stops processing and returns error if any row fails to scan
//
// Performance considerations:
//   - Uses reflection, which has some performance overhead. Suitable for most applications, but consider manual scanning for high-performance scenarios
//
// Error handling:
//   - Query execution errors: returned immediately
//   - Row scanning errors: stops processing and returns error
//   - Row iteration errors: checked after processing all rows
//   - No rows found: returns empty slice with no error (not sql.ErrNoRows)
func Query[T any](db *sql.DB, query string, args []any) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // TODO Controlar el error que puede devolver esto

	var results []T

	for rows.Next() {
		var instance T

		err := initializeInstanceFromScanner(rows, &instance)
		if err != nil {
			return nil, err
		}

		results = append(results, instance)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
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
//
// Notes:
//   - Returns 0 if no rows match the WHERE condition
//   - Returns the actual number of rows modified, not rows examined
//   - If SET values are the same as existing values, behavior depends on database:
//   - MySQL: May return 0 (no actual change) or row count (rows examined)
//   - PostgreSQL: Returns rows that matched the condition
//   - Always use parameterized queries (?) to prevent SQL injection
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
//
// Notes:
//   - Returns 0 if no rows match the WHERE condition
//   - BE VERY CAREFUL with DELETE without WHERE clause - it deletes ALL rows
//   - Consider using soft deletes (UPDATE with deleted flag) for important data
//   - Foreign key constraints may prevent deletion if related records exist
//   - Always use parameterized queries (?) to prevent SQL injection
func Delete(db *sql.DB, query string, args []any) (int64, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
