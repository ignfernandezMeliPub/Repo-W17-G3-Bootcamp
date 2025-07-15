package utils

import (
	"app/pkg/custom_errors"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/bootcamp-go/web/response"
)

const sqlUniqueAttributeViolationErrString = "Error 1062 (23000): Duplicate entry"

func ResponseHttpError(w http.ResponseWriter, err error) {
	var status int
	var message string

	errorMsg := err.Error()
	println(errorMsg)

	switch {
	case errors.As(err, &custom_errors.ErrNotFound) || errors.Is(err, sql.ErrNoRows):
		status = http.StatusNotFound
		message = "Not found"
	case errors.As(err, &custom_errors.ErrInvalidBodyError) || errors.As(err, &custom_errors.ErrDecodeError) || errors.As(err, &custom_errors.UrlParamDecodeErrorI) || errors.As(err, &custom_errors.QueryParamDecodeErrorI):
		status = http.StatusBadRequest
		message = "Bad request"
	case errors.As(err, &custom_errors.ErrUniqueAttributeViolationError) || strings.Contains(err.Error(), sqlUniqueAttributeViolationErrString) || errors.As(err, &custom_errors.ErrForeignKeyViolation):
		status = http.StatusConflict
		message = "Conflict"
	case errors.As(err, &custom_errors.ErrInvalidArgs) || errors.As(err, &custom_errors.ErrMandatoryArgMissing):
		status = http.StatusUnprocessableEntity
		message = "Unprocessable Entity"
	default:
		status = http.StatusInternalServerError
		message = "Internal server error"
		errorMsg = "Internal server error"
	}

	response.JSON(w, status, map[string]any{
		"message": message,
		"error":   errorMsg,
	})
}
