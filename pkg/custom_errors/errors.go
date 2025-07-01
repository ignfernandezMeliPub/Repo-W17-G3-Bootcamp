package custom_errors

import "fmt"

type ResourceNotFoundError struct{}

var ErrNotFound = &ResourceNotFoundError{}

func (e *ResourceNotFoundError) Error() string {
	return "resource not found."
}

type InvalidArgValueErr struct {
	Argument  string
	Value     any
	ExtraInfo string
}

var ErrInvalidArgs = &InvalidArgValueErr{}

func (e *InvalidArgValueErr) Error() string {
	if e.ExtraInfo != "" {
		return fmt.Sprintf("Invalid Value {%v} for arg {%s}. %s.", e.Value, e.Argument, e.ExtraInfo)
	} else {
		return fmt.Sprintf("Invalid Value {%v} for arg {%s}.", e.Value, e.Argument)
	}
}

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

type DecodeError struct {
	// The struct field that failed to decode.
	FieldName string

	// Expected Go type for the field.
	ExpectedType string

	// The actual value/type found in the JSON body.
	FoundType string

	// The underlying decoding error.
	BaseErr error
}

var ErrDecodeError = &DecodeError{}

// Error @inheritdoc
func (e *DecodeError) Error() string {
	return fmt.Sprintf("field '%s' is expected to be of type '%s', but found '%s'", e.FieldName, e.ExpectedType, e.FoundType)
}

// Unwrap returns the underlying error for error unwrapping compatibility.
func (e *DecodeError) Unwrap() error {
	return e.BaseErr
}
