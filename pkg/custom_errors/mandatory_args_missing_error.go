package custom_errors

import "fmt"

type MandatoryArgMissingErr struct {
	Argument string
}

var ErrMandatoryArgMissing = &MandatoryArgMissingErr{}

func (e *MandatoryArgMissingErr) Error() string {
	return fmt.Sprintf("Argument/s {%v} is/are mandatory", e.Argument)
}
