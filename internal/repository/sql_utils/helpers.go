package sql_utils

import (
	"database/sql"
	"fmt"
	"reflect"
)

const dbTagName = "db"

// initInstanceWithRows scans the current row from a *sql.Rows result set into a
// destination struct instance.
//
// It maps database columns to struct fields by matching column names with the
// struct's dbTagName tags. This makes the scanning process robust, as it does not
// depend on the column order in the SELECT statement.
//
// The function is designed to be called within a `for rows.Next()` loop.
//
// Parameters:
//   - rows: The active *sql.Rows iterator, positioned on a valid row.
//   - instance: A pointer to the struct to be populated.
//
// Behavior:
//   - Only exported struct fields with a non-empty dbTagName tag are considered for mapping.
//   - Columns returned from the query that do not have a corresponding tagged field
//     in the struct are safely ignored.
func initInstanceWithRows(rows *sql.Rows, instance any) error {

	// ? Obtener los nombres de las columnas del resultado de la query.
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// ? Creamos un mapa que mapeara 'dbTagName -> structAttribute'
	fieldMap := make(map[string]reflect.Value)

	// ? Obtenemos los tags de db de la estructura recibida, se hace uso de la recursividad en caso de reribir un struct con strutcs embebidos
	fieldMap, err = getFields(instance, fieldMap)

	if err != nil {

		return err

	}

	// ? Preparamos el slice de punteros que 'Scan()' inicializara. Este slice DEBE tener el mismo orden y tamaño que las columnas de la query.
	scanDest := make([]any, len(columns))

	// ? Recorremos cada una de las columnas obtenidas en la query
	for i, colName := range columns {
		// ? Buscamos si hay un atributo en nuestro struct para esta columna.
		field, ok := fieldMap[colName]
		if ok {
			// ? Si se encuentra, se pone un puntero a ese campo como destino del escaneo.
			scanDest[i] = field.Addr().Interface()
		} else {
			// ? Si no hay un campo para esta columna, se escanea en un receptor temporal para que sea ignorado. De no hacerlo, Scan() fallaría.
			scanDest[i] = new(sql.RawBytes)
		}
	}

	// ? Ejecutamos 'Scan()' para inicializar los pointers
	return rows.Scan(scanDest...)
}

// getFields extracts struct fields and their reflection values, mapping them by database tag names.
//
// This function recursively traverses a struct (including embedded structs) and builds a map
// where keys are database tag names and values are the corresponding reflect.Value objects.
// It's designed to work with structs that use `db` tags for database field mapping.
//
// Parameters:
//   - instance: Must be a pointer to a struct. The struct to extract fields from
//   - fieldMap: A map to populate with the extracted fields. Pass an empty map or existing map
//     Format: map[dbTagName]reflect.Value
//
// Returns:
//   - map[string]reflect.Value: The populated fieldMap with database tag names as keys
//     and settable reflect.Value objects as values
//   - error: Possible errors include:
//   - "instance must be a pointer": When instance is not a pointer
//   - "instance must be a pointer to a struct": When instance doesn't point to a struct
//   - Recursive errors from embedded struct processing
//
// Behavior:
//   - Only includes fields that have a `db` tag defined
//   - Only includes fields that are settable (CanSet() returns true)
//   - Recursively processes embedded structs by taking their address
//   - Skips fields without `db` tags or that are not settable
//
// Example usage:
//
//	type User struct {
//	    ID   int    `db:"id"`
//	    Name string `db:"name"
//	}
//
//	user := &User{}
//	fieldMap := make(map[string]reflect.Value)
//	result, err := getFields(user, fieldMap)
//	// result["id"] contains reflect.Value for user.ID
//	// result["name"] contains reflect.Value for user.Name
//
// Note: This function modifies the passed fieldMap parameter and also returns it.
func getFields(instance any, fieldMap map[string]reflect.Value) (map[string]reflect.Value, error) {
	// ? Validar que 'instance' sea un puntero
	instancePtr := reflect.ValueOf(instance)
	if instancePtr.Kind() != reflect.Ptr {
		return fieldMap, fmt.Errorf("instance must be a pointer")
	}

	// ? Validar que 'instance' apunte a un struct
	instanceValue := instancePtr.Elem()
	if instanceValue.Kind() != reflect.Struct {
		return fieldMap, fmt.Errorf("instance must be a pointer to a struct")
	}

	// ? Llenamos el mapa con el formato 'dbTagName -> structAttribute'
	instanceType := instanceValue.Type()

	for i := 0; i < instanceValue.NumField(); i++ {
		field := instanceType.Field(i)
		fieldValue := instanceValue.Field(i)
		dbTag := field.Tag.Get(dbTagName)

		// ? Si es un struct embebido, ejecutamos recursivamente getFields()
		if field.Type.Kind() == reflect.Struct {
			fieldPtr := fieldValue.Addr()
			_, err := getFields(fieldPtr.Interface(), fieldMap)

			if err != nil {
				return fieldMap, err
			}
			// ? Si es un atributo de dato primitivo lo agregamos al mapa si tiene el tag db y es seteable
		} else if dbTag != "" && fieldValue.CanSet() {
			fieldMap[dbTag] = instanceValue.Field(i)
		}
	}

	return fieldMap, nil
}
