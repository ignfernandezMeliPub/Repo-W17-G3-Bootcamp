package custom_errors

import "fmt"

type ResourceConflictError struct {
	Argument  string
	Value     string
	ExtraInfo string
}

var ErrConflictError = &ResourceConflictError{}

func (e *ResourceConflictError) Error() string {
	if e.ExtraInfo != "" {
		return fmt.Sprintf("Conflict: {%v} for arg {%s}. %s.", e.Value, e.Argument, e.ExtraInfo)
	} else {
		return fmt.Sprintf("Conflict:{%v} for arg {%s}.", e.Value, e.Argument)
	}
}
