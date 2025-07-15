package utils

import (
	"app/pkg/custom_errors"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetURLParamAs retrieves the value of a URL parameter from the request by key (`urlParamKey`),
// and parses it using the given `parser` function, returning the desired type.
//
// Example usage (parsing an integer parameter):
//
//	id, err := GetURLParamAs[int](r, "userId", strconv.Atoi)
func GetURLParamAs[T any](request *http.Request, urlParamKey string, parser func(string) (T, error)) (T, error) {
	valueStr := chi.URLParam(request, urlParamKey)

	result, err := parser(valueStr)
	if err != nil {
		return result, &custom_errors.UrlParamDecodeError{UrlParam: urlParamKey, BaseErr: err}
	}

	return result, nil
}

func GetQueryParamAs[T any](request *http.Request, queryParamKey string, parser func(string) (T, error)) (*T, error) {
	valueStr := request.URL.Query().Get(queryParamKey)

	if valueStr == "" {
		return nil, nil
	}

	result, err := parser(valueStr)
	if err != nil {
		return nil, &custom_errors.QueryParamDecodeError{QueryParam: queryParamKey, BaseErr: err}
	}

	return &result, nil
}

// BodyInstantiableStruct should be implemented by any struct type intended to be
// instantiated from a request body. The VerifyMandatoryFieldsPresence method is used
// to enforce presence/validation of required fields after decoding from JSON.
type BodyInstantiableStruct interface {
	// VerifyMandatoryFieldsPresence checks that all required fields in the struct
	// are present and returns an error if not, or nil if validation passes.
	VerifyMandatoryFieldsPresence() error
}

// InstantiateVarFromBody attempts to decode the request body as JSON into the given variable of type T,
// where T must implement BodyInstantiableStruct (for mandatory field validation).
//
// If there is an error decoding the JSON, and it is due to a type mismatch,
// a DecodeError is returned indicating the specific field and expected type.
//
// After successful decoding, VerifyMandatoryFieldsPresence is called and its error (if any) is returned.
func InstantiateVarFromBody[T BodyInstantiableStruct](body *io.ReadCloser, variable T) (T, error) {
	err := json.NewDecoder(*body).Decode(&variable)
	if err != nil {
		var cmpErr *json.UnmarshalTypeError

		if errors.As(err, &cmpErr) {
			return variable, &custom_errors.DecodeError{
				FieldName:    cmpErr.Field,
				ExpectedType: cmpErr.Type.String(),
				FoundType:    cmpErr.Value,
				BaseErr:      err,
			}
		}

		return variable, &custom_errors.InvalidBodyError{}
	}

	return variable, variable.VerifyMandatoryFieldsPresence()
}

// GetHeader returns the value of the specified HTTP header key.
// The Header.Get method retrieves the first value associated with the given key.
// The lookup is case-insensitive and will return an empty string if the key is not present.
func GetHeader(headers *http.Header, key string) string {
	return headers.Get(key)
}
