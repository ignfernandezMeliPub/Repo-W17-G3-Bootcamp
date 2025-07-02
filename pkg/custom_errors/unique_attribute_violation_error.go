package custom_errors

// Deprecated. Use conflict error instead
import "fmt"

type UniqueAttributeViolationErr struct {
	AttributeName string
	Value         any
}

var ErrUniqueAttributeViolationError = &UniqueAttributeViolationErr{}

func (e *UniqueAttributeViolationErr) Error() string {
	return fmt.Sprintf("Invalid value {%v} for unique attribute {%s}. Value already being used.", e.Value, e.AttributeName)
}
