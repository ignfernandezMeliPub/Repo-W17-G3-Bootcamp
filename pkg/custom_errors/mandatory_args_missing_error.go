package custom_errors

import "fmt"

type MandatoryArgMissingErr struct {
	Argument string
}

var ErrMandatoryArgMissing = &MandatoryArgMissingErr{}

func (e *MandatoryArgMissingErr) Error() string {
	return fmt.Sprintf("Argument {%v} is mandatory", e.Argument)
}

type UniqueAttributeViolationErr struct {
	AttributeName string
	Value         any
}

func (e *UniqueAttributeViolationErr) Error() string {
	return fmt.Sprintf("Invalid value {%v} for unique attribute {%s}. Value already being used.", e.Value, e.AttributeName)
}
