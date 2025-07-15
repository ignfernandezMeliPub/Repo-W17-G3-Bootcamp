package custom_errors

import "fmt"

type QueryParamDecodeError struct {
	// The param that failed to be decoded.
	QueryParam string

	// The underlying decoding error.
	BaseErr error
}

var QueryParamDecodeErrorI = &QueryParamDecodeError{}

// Error @inheritdoc
func (e *QueryParamDecodeError) Error() string {
	return fmt.Sprintf("Failed to decode query param {%s}. Please verify format.", e.QueryParam)
}
