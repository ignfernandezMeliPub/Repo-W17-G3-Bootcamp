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
	// Verify checks that all fields in the struct pass validations.
	Verify() error
}

// InstantiateVarFromBody attempts to decode the request body as JSON into the given variable of type T,
// where T must implement BodyInstantiableStruct (for mandatory field validation).
//
// If there is an error decoding the JSON, and it is due to a type mismatch,
// a DecodeError is returned indicating the specific field and expected type.
//
// After successful decoding, Verify is called and its error (if any) is returned.
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

	return variable, variable.Verify()
}
