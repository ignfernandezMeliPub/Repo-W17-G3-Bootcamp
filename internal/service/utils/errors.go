package utils

import "errors"

func ExpectError(actual error, expectedError error, onNoError error) error {
	if actual == nil {
		return onNoError
	} else if !errors.As(actual, &expectedError) {
		return actual
	}
	return nil
}

func ExpectErrorOrNilCondition(actual error, condition bool, expectedError error, onNoError error) error {
	if actual == nil && condition {
		return onNoError
	} else if !errors.As(actual, &expectedError) {
		return actual
	}
	return nil
}
