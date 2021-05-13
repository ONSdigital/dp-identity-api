package apierrors

import (
	"errors"
	"strings"
)

var InvalidTokenError = errors.New("Invalid token")
var MissingTokenMessage = "No Authorization token was provided"
var MalformedTokenMessage = "The provided token does not meet the required format"

type ErrorStructure struct {
	Errors []IndividualError `json:"errors"`
}

type IndividualError struct {
	SpecificError string `json:"error"`
	Message       string `json:"message"`
	Source        Source `json:"source"`
}

type Source struct {
	Field string `json:"field"`
	Param string `json:"param"`
}

func IndividualErrorBuilder(err error, message, sourceField, sourceParam string) (individualError IndividualError) {

	individualError = IndividualError{
		SpecificError: error.Error(err),
		Message:       message,
		Source: Source{
			Field: sourceField,
			Param: sourceParam},
	}

	return individualError
}

func ErrorResponseBodyBuilder(listOfErrors []IndividualError) (errorResponseBody ErrorStructure) {

	errorResponseBody = ErrorStructure{
		Errors: listOfErrors,
	}

	return errorResponseBody
}

func IdentifyInternalError(authErr error) (isInternalError bool) {
	possibleInternalErrors := []string{
		"InternalErrorException",
		"SerializationError",
		"ReadError",
		"ResponseTimeout",
		"InvalidPresignExpireError",
		"RequestCanceled",
		"RequestError",
	}

	for _, internalError := range possibleInternalErrors {
		if strings.Contains(authErr.Error(), internalError) {
			return true
		}
	}

	//strings.Contains(authErr.Error(), "InternalErrorException") internalError == authErr.Error()
	return false

}
