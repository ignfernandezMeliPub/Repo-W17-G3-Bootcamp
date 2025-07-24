package utils

import (
	"app/pkg/custom_errors"
	"bytes"
	"context"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}

func (t testStruct) Verify() error {
	if t.Name == nil {
		return &custom_errors.MandatoryArgMissingErr{Argument: "name"}
	}

	if strings.TrimSpace(*t.Name) == "" {
		return &custom_errors.InvalidArgValueErr{
			Argument:  "name",
			Value:     *t.Name,
			ExtraInfo: "Value must be non-empty",
		}
	}
	return nil
}

func TestInstantiateVarFromBody(t *testing.T) {

	t.Run("should instantiate the struct successfully", func(t *testing.T) {
		body := `{"name":"Alice","age":30}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))

		var result testStruct
		result, err := InstantiateVarFromBody(&req.Body, result)

		if err != nil {
			t.Fatalf("InstantiateVarFromBody failed: %v", err)
		}

		require.NoError(t, err)
		require.Equal(t, "Alice", *result.Name)
		require.Equal(t, 30, *result.Age)
	})

	t.Run("should return an error if the body is invalid", func(t *testing.T) {
		body := `{"name":"","age":30`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
		var result testStruct
		_, err := InstantiateVarFromBody(&req.Body, result)

		require.Error(t, err)
		require.IsType(t, &custom_errors.InvalidBodyError{}, err)
	})

	t.Run("should return an error if the body is valid but the datatype is invalid", func(t *testing.T) {
		body := `{"name":"Alice","age":"30"}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))
		var result testStruct
		_, err := InstantiateVarFromBody(&req.Body, result)

		require.Error(t, err)
		require.IsType(t, &custom_errors.DecodeError{}, err)
	})
}

func TestGetURLParamAs(t *testing.T) {
	t.Run("should return the param as the expected type", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/test/123", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		param, err := GetURLParamAs(req, "id", strconv.Atoi)
		require.NoError(t, err)
		require.Equal(t, 123, param)
	})

	t.Run("should return an error if the param is not parseable to the expected type", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/test/abc", nil)
		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("id", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

		_, err := GetURLParamAs(req, "id", strconv.Atoi)
		require.Error(t, err)
		require.IsType(t, &custom_errors.UrlParamDecodeError{}, err)
	})
}

func TestGetQueryParamAs(t *testing.T) {
	t.Run("should return the param as the expected type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?id=123", nil)
		param, err := GetQueryParamAs(req, "id", strconv.Atoi)
		require.NoError(t, err)
		require.Equal(t, 123, *param)
	})

	t.Run("should return an error if the param is not parseable to the expected type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?id=abc", nil)
		_, err := GetQueryParamAs(req, "id", strconv.Atoi)
		require.Error(t, err)
		require.IsType(t, &custom_errors.QueryParamDecodeError{}, err)
	})

	t.Run("should return nil if the param is not present", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		param, err := GetQueryParamAs(req, "id", strconv.Atoi)
		require.NoError(t, err)
		require.Nil(t, param)
	})
}
