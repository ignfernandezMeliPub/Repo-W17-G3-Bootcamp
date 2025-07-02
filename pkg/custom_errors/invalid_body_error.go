package custom_errors

type InvalidBodyError struct{}

var ErrInvalidBodyError = &InvalidBodyError{}

// Error @inheritdoc
func (e *InvalidBodyError) Error() string {
	return "Invalid body"
}
