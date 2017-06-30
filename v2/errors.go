package v2

import (
	"fmt"
	"net/http"
)

// HTTPStatusCodeError is an error type that provides additional information
// based on the Open Service Broker API conventions for returning information
// about errors.  If the response body provided by the broker to any client
// operation is malformed, an error of this type will be returned with the
// ResponseError field set to the unmarshalling error.
//
// These errors may optionally provide a machine-readable error message and
// human-readable description.
//
// The IsHTTPError method checks whether an error is of this type.
//
// Checks for important errors in the API specification are implemented as
// utility methods:
//
// - IsGoneError
// - IsConflictError
// - IsAsyncRequiredError
// - IsAppGUIDRequiredError
type HTTPStatusCodeError struct {
	// StatusCode is the HTTP status code returned by the broker.
	StatusCode int
	// ErrorMessage is a machine-readable error string that may be returned by
	// the broker.
	ErrorMessage *string
	// Description is a human-readable description of the error that may be
	// returned by the broker.
	Description *string
	// ResponseError is set to the error that occured when unmarshalling a
	// response body from the broker.
	ResponseError error
}

func (e HTTPStatusCodeError) Error() string {
	return fmt.Sprintf("Status: %v; ErrorMessage: %v; Description: %v; ResponseError: %v", e.StatusCode, e.ErrorMessage, e.Description, e.ResponseError)
}

// IsHTTPError returns whether the error represents an HTTPStatusCodeError.  A
// client method returning an HTTP error indicates that the broker returned an
// error code and a correctly formed response body.
func IsHTTPError(err error) bool {
	_, ok := err.(HTTPStatusCodeError)
	return ok
}

// IsGoneError returns whether the error represents an HTTP GONE status.
func IsGoneError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	return statusCodeError.StatusCode == http.StatusGone
}

// IsConflictError returns whether the error represents a conflict.
func IsConflictError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	return statusCodeError.StatusCode == http.StatusConflict
}

const (
	AsyncErrorMessage     = "AsyncRequired"
	AsyncErrorDescription = "This service plan requires client support for asynchronous service operations."

	AppGUIDRequiredErrorMessage     = "RequiresApp"
	AppGUIDRequiredErrorDescription = "This service supports generation of credentials through binding an application only."
)

// IsAsyncRequiredError returns whether the error corresponds to the
// conventional way of indicating that a service requires asynchronous
// operations to perform an action.
func IsAsyncRequiredError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	if statusCodeError.StatusCode != http.StatusUnprocessableEntity {
		return false
	}

	if *statusCodeError.ErrorMessage != AsyncErrorMessage {
		return false
	}

	return *statusCodeError.Description == AsyncErrorDescription
}

// IsAppGUIDRequiredError returns whether the error corresponds to the
// conventional way of indicating that a service only supports credential-type
// bindings.
func IsAppGUIDRequiredError(err error) bool {
	statusCodeError, ok := err.(HTTPStatusCodeError)
	if !ok {
		return false
	}

	if statusCodeError.StatusCode != http.StatusUnprocessableEntity {
		return false
	}

	if *statusCodeError.ErrorMessage != AppGUIDRequiredErrorMessage {
		return false
	}

	return *statusCodeError.Description == AppGUIDRequiredErrorDescription
}
