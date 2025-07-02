package custom_errors

import "fmt"

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
