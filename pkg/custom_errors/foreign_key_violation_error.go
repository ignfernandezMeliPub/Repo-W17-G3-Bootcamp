package custom_errors

import "fmt"

type ForeignKeyViolationError struct {
	ConstraintName string
	Details        string
}

var ErrForeignKeyViolation = &ForeignKeyViolationError{}

func (e *ForeignKeyViolationError) Error() string {
	if e.ConstraintName != "" {
		return fmt.Sprintf("Unknown entity identifier value: %s ", e.ConstraintName)
	}
	return fmt.Sprintf("Unknown entity identifier value")
}
