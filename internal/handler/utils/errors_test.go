package utils

import (
	"app/pkg/custom_errors"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseHttpError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{name: "ErrNotFound", err: custom_errors.ErrNotFound, want: http.StatusNotFound},
		{name: "ErrInvalidBodyError", err: custom_errors.ErrInvalidBodyError, want: http.StatusBadRequest},
		{name: "ErrDecodeError", err: custom_errors.ErrDecodeError, want: http.StatusBadRequest},
		{name: "ErrUrlParamDecodeErrorI", err: custom_errors.UrlParamDecodeErrorI, want: http.StatusBadRequest},
		{name: "ErrQueryParamDecodeErrorI", err: custom_errors.QueryParamDecodeErrorI, want: http.StatusBadRequest},
		{name: "ErrUniqueAttributeViolationError", err: custom_errors.ErrUniqueAttributeViolationError, want: http.StatusConflict},
		{name: "ErrForeignKeyViolation", err: custom_errors.ErrForeignKeyViolation, want: http.StatusConflict},
		{name: "ErrInvalidArgs", err: custom_errors.ErrInvalidArgs, want: http.StatusUnprocessableEntity},
		{name: "ErrMandatoryArgMissing", err: custom_errors.ErrMandatoryArgMissing, want: http.StatusUnprocessableEntity},
		{name: "Unhandled error", err: errors.New("unhandled error"), want: http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ResponseHttpError(w, tc.err)

			require.Equal(t, tc.want, w.Result().StatusCode)
		})
	}

	t.Run("Should hide internal server error message", func(t *testing.T) {
		message := "This message must be hidden"
		w := httptest.NewRecorder()
		ResponseHttpError(w, errors.New(message))

		require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		response := map[string]any{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.NotContains(t, response["error"], message)
	})
}
