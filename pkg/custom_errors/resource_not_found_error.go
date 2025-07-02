package custom_errors

type ResourceNotFoundError struct{}

var ErrNotFound = &ResourceNotFoundError{}

func (e *ResourceNotFoundError) Error() string {
	return "resource not found."
}
