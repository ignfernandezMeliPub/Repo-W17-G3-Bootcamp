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
	return fmt.Sprintf("Failed to decode url param {%s}. Please verify format.", e.UrlParam)
}
