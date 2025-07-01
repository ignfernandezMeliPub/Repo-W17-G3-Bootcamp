package custom_errors

import "fmt"

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
