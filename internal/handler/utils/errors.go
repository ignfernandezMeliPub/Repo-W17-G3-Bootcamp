package utils

import (
	"app/pkg/custom_errors"
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

func MapErrorToHttpStatus(err error) int {
	switch {
	case errors.As(err, &custom_errors.ErrNotFound):
		return http.StatusNotFound
	case errors.As(err, &custom_errors.ErrNotFound) || errors.As(err, &custom_errors.ErrInvalidArgs) || errors.As(err, &custom_errors.ErrDecodeError):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func ResponseHttpError(w http.ResponseWriter, err error) {
	var status int
	var message string

	error_msg := err.Error()

	switch {
	case errors.As(err, &custom_errors.ErrNotFound):
		status = http.StatusNotFound
		message = "Not found"
	case errors.As(err, &custom_errors.ErrNotFound) || errors.As(err, &custom_errors.ErrInvalidArgs) || errors.As(err, &custom_errors.ErrDecodeError):
		status = http.StatusBadRequest
		message = "Bad request"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
		error_msg = "Internal server error" // Don't expose unhandled error
	}

	response.JSON(w, status, map[string]any{
		"message": message,
		"error":   error_msg,
	})
}
