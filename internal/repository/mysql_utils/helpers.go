package mysql_utils

import (
	"fmt"
	"reflect"
)

// Scanner interface unifies sql.Row and sql.Rows Scan method
type Scanner interface {
	Scan(dest ...any) error
}

// initializeInstanceFromScanner maps a SQL result to a Go struct using
// reflection and 'db' tags. Works with both sql.Row and sql.Rows.
//
// This function uses reflection to:
// 1. Validate that the instance parameter is a pointer
// 2. Examine the struct fields
// 3. Identify fields with 'db' tags
// 4. Create a slice of pointers to those fields
// 5. Execute scanner.Scan() to populate the fields
//
// Parameters:
//   - scanner: Any type that implements Scanner interface (sql.Row or sql.Rows)
//     Contains the SQL result data to be mapped
//   - instance: Pointer to the struct to be populated with data.
//     MUST be a pointer, otherwise returns an error.
//
// Returns:
//   - error: nil if operation is successful, specific error if it fails
//
// Behavior:
//
//   - Only processes exported fields (starting with uppercase letter)
//
//   - Only processes fields that have non-empty 'db' tag
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
// Possible errors:
//   - "instance must be a pointer": if a non-pointer is passed
//   - sql.ErrNoRows: if the query returns no rows (for sql.Row)
//   - Type conversion errors if Go types don't match database types
//   - Database connection errors
//   - Field count mismatch errors if SELECT columns don't match struct fields
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
func initializeInstanceFromScanner(scanner Scanner, instance any) error {
	// ? Primero verificamos que instance sea un pointer. Si no lo es, error (no deberia ocurrir, esta funcion es para uso interno del paquete)
	metadata := reflect.ValueOf(instance)
	if metadata.Kind() != reflect.Ptr {
		return fmt.Errorf("instance must be a pointer")
	}

	// ? Obtenemos el valor a donde apunta 'instance'
	value := metadata.Elem()

	// ? Verificamos que instance sea un pointer a un struct
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("instance must be a pointer to a struct")
	}

	// ? Obtenemos el tipo del dato al que apunta 'instance'
	instanceDataType := value.Type()

	// ? Generamos un array en el que guardaremos pointers a cada uno de los atributos de 'instance' a inicializar
	var attributesToInitialize []any

	// ? Recorremos cada uno de los atributos del tipo de dato de 'instance'
	for i := 0; i < value.NumField(); i++ {
		// ? Obtenemos el atributo actual
		field := value.Field(i)

		// ? Obtenemos la metadata del atributo actual, y de esta su tag 'db'.
		// ? Se espera que el tag 'db' indique el nombre de la columna que representa ese atributo
		dbTag := instanceDataType.Field(i).Tag.Get("db")

		// ? Verificamos si dbTag estaba seteado y si el campo es modificable (esta exportado y no es readonly)
		if dbTag != "" && field.CanSet() {
			// ? Caso afirmativo, obtenemos el pointer al atributo que estamos analizando y lo agregamos a la lista de atributos a inicializar
			attributesToInitialize = append(attributesToInitialize, field.Addr().Interface())
		}
	}

	// ? Inicializamos los pointers, y por ende, la estructura
	return scanner.Scan(attributesToInitialize)
}
