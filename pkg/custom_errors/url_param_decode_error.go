package custom_errors

import "fmt"

type UrlParamDecodeError struct {
	// The param that failed to be decoded.
	UrlParam string

	// The underlying decoding error.
	BaseErr error
}

var UrlParamDecodeErrorI = &UrlParamDecodeError{}

// Error @inheritdoc
func (e *UrlParamDecodeError) Error() string {
	return fmt.Sprintf("field '%s' failed to be decoded. Please verify format.", e.UrlParam)
}
