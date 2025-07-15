package custom_errors

import "fmt"

type ForeignKeyViolationError struct {
	ConstraintName string
	IsParentRow    bool
	Details        string
}

var ErrForeignKeyViolation = &ForeignKeyViolationError{}

func (e *ForeignKeyViolationError) Error() string {
	if e.IsParentRow {
		if e.ConstraintName != "" {
			return fmt.Sprintf("The entity cannot be deleted because there are other records that depend on it: %s ", e.ConstraintName)
		}
		return fmt.Sprintf("The entity cannot be deleted because there are other records that depend on it.")
	} else {
		if e.ConstraintName != "" {
			return fmt.Sprintf("Unknown entity identifier value: %s ", e.ConstraintName)
		}
		return fmt.Sprintf("Unknown entity identifier value")
	}
}
